package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yahao333/gort/internal/core"
)

var planCmd = &cobra.Command{
	Use:   "plan [environment]",
	Short: "Plan infrastructure changes for an environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env := args[0]
		deployer := core.NewDeployer()
		// TODO: Implement plan logic
		return nil
	},
}

func init() {
	planCmd.Flags().StringP("version", "v", "", "Version to plan")
}
