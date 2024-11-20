package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/yahao333/gort/internal/config"
	"github.com/yahao333/gort/internal/core"
	"github.com/yahao333/gort/internal/logging"
	"github.com/yahao333/gort/internal/plugin"
	"github.com/yahao333/gort/internal/state"
)

type deployOptions struct {
	version     string
	force       bool
	parallel    int
	timeout     time.Duration
	configFile  string
	stateDir    string
	pluginDir   string
	backupState bool
}

var deployOpts = &deployOptions{}

var deployCmd = &cobra.Command{
	Use:   "deploy [environment]",
	Short: "Deploy infrastructure to an environment",
	Long: `Deploy infrastructure to the specified environment.
    
This command will:
1. Load and validate environment configuration
2. Plan the deployment changes
3. Execute the deployment
4. Update the state`,
	Args: cobra.ExactArgs(1),
	RunE: runDeploy,
}

func init() {
	// Add flags
	deployCmd.Flags().StringVarP(&deployOpts.version, "version", "v", "", "Version to deploy")
	deployCmd.Flags().BoolVarP(&deployOpts.force, "force", "f", false, "Force deployment without confirmation")
	deployCmd.Flags().IntVarP(&deployOpts.parallel, "parallel", "p", 1, "Number of parallel deployments")
	deployCmd.Flags().DurationVar(&deployOpts.timeout, "timeout", 30*time.Minute, "Deployment timeout")
	deployCmd.Flags().StringVar(&deployOpts.configFile, "config", "gort.yaml", "Path to config file")
	deployCmd.Flags().StringVar(&deployOpts.stateDir, "state-dir", ".gort/state", "Directory for state files")
	deployCmd.Flags().StringVar(&deployOpts.pluginDir, "plugin-dir", ".gort/plugins", "Directory for plugins")
	deployCmd.Flags().BoolVar(&deployOpts.backupState, "backup-state", true, "Backup state before deployment")
}

func runDeploy(cmd *cobra.Command, args []string) error {
	// Setup context with timeout
	ctx, cancel := context.WithTimeout(cmd.Context(), deployOpts.timeout)
	defer cancel()

	// Initialize logger
	logger := logging.NewLogger(os.Getenv("DEBUG") == "true")
	logger.Info("Starting deployment process")

	// Get environment name
	envName := args[0]

	// Load configuration
	cfg, err := loadConfig(deployOpts.configFile, envName)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize state manager
	stateManager := state.NewStateManager(deployOpts.stateDir)

	// Backup state if enabled
	if deployOpts.backupState {
		if err := backupState(stateManager, envName); err != nil {
			return fmt.Errorf("failed to backup state: %w", err)
		}
	}

	// Initialize plugin manager
	pluginManager := plugin.NewPluginManager(deployOpts.pluginDir, logger)
	if err := pluginManager.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize plugin manager: %w", err)
	}

	// Create deployer
	deployer := core.NewDeployer(
		stateManager,
		pluginManager,
		logger,
		core.DeployerOptions{
			Parallel: deployOpts.parallel,
			Force:    deployOpts.force,
		},
	)

	// Create deployment plan
	plan, err := deployer.Plan(ctx, envName, cfg)
	if err != nil {
		return fmt.Errorf("failed to create deployment plan: %w", err)
	}

	// Show plan and confirm if not forced
	if !deployOpts.force {
		if err := confirmDeployment(plan); err != nil {
			return err
		}
	}

	// Execute deployment
	result, err := deployer.Deploy(ctx, plan)
	if err != nil {
		logger.Errorf("Deployment failed: %v", err)
		if err := handleDeploymentFailure(ctx, deployer, plan); err != nil {
			logger.Errorf("Failed to handle deployment failure: %v", err)
		}
		return fmt.Errorf("deployment failed: %w", err)
	}

	// Show deployment results
	showDeploymentResults(result)

	logger.Info("Deployment completed successfully")
	return nil
}

func loadConfig(configFile, envName string) (*config.Config, error) {
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return nil, err
	}

	// Validate environment exists
	if _, exists := cfg.Environments[envName]; !exists {
		return nil, fmt.Errorf("environment '%s' not found in configuration", envName)
	}

	return cfg, nil
}

func backupState(sm *state.StateManager, envName string) error {
	backupPath := fmt.Sprintf("%s.backup-%s", envName, time.Now().Format("20060102-150405"))
	return sm.BackupState(envName, backupPath)
}

func confirmDeployment(plan *core.DeploymentPlan) error {
	fmt.Println("\nDeployment Plan:")
	fmt.Println("================")
	fmt.Printf("Environment: %s\n", plan.Environment)
	fmt.Printf("Resources to Add: %d\n", len(plan.AddResources))
	fmt.Printf("Resources to Update: %d\n", len(plan.UpdateResources))
	fmt.Printf("Resources to Delete: %d\n", len(plan.DeleteResources))
	fmt.Println("\nDo you want to proceed? (yes/no)")

	var response string
	fmt.Scanln(&response)
	if response != "yes" {
		return fmt.Errorf("deployment cancelled by user")
	}

	return nil
}

func handleDeploymentFailure(ctx context.Context, deployer *core.Deployer, plan *core.DeploymentPlan) error {
	fmt.Println("\nDeployment failed. Attempting rollback...")

	if err := deployer.Rollback(ctx, plan); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	fmt.Println("Rollback completed successfully")
	return nil
}

func showDeploymentResults(result *core.DeploymentResult) {
	fmt.Println("\nDeployment Results:")
	fmt.Println("===================")
	fmt.Printf("Environment: %s\n", result.Environment)
	fmt.Printf("Duration: %s\n", result.Duration)
	fmt.Printf("Resources Created: %d\n", len(result.CreatedResources))
	fmt.Printf("Resources Updated: %d\n", len(result.UpdatedResources))
	fmt.Printf("Resources Deleted: %d\n", len(result.DeletedResources))

	if len(result.Outputs) > 0 {
		fmt.Println("\nOutputs:")
		for k, v := range result.Outputs {
			fmt.Printf("%s: %v\n", k, v)
		}
	}
}
