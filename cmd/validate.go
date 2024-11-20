package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yahao333/gort/internal/config"
	"github.com/yahao333/gort/internal/provider/terraform"
)

var validateCmd = &cobra.Command{
    Use:   "validate [environment]",
    Short: "Validate configuration and terraform files",
    Args:  cobra.MaximumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // Load and validate config
        cfg, err := config.LoadConfig("")
        if err != nil {
            return fmt.Errorf("configuration validation failed: %w", err)
        }

        // If environment specified, validate specific environment
        if len(args) > 0 {
            env := args[0]
            if _, exists := cfg.Environments[env]; !exists {
                return fmt.Errorf("environment '%s' not found in configuration", env)
            }

            // Validate terraform configuration
            provider := terraform.NewTerraformProvider(".")
            if err := provider.Validate(env); err != nil {
                return fmt.Errorf("terraform validation failed: %w", err)
            }
        }

        fmt.Println("Validation successful!")
        return nil
    },
}

func init() {
    rootCmd.AddCommand(validateCmd)
}