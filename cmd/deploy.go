package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yahao333/gort/internal/core"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [environment]",
	Short: "Deploy infrastructure to an environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env := args[0]
		deployer := core.NewDeployer()
		return deployer.Deploy(env)
	},
}

func init() {
	deployCmd.Flags().StringP("version", "v", "", "Version to deploy")
	deployCmd.Flags().BoolP("force", "f", false, "Force deployment")
}
