package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gort",
	Short: "GoRT - Infrastructure Release Tool",
	Long: `GoRT is a tool for managing infrastructure deployments
           across different environments using terraform.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add sub-commands
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(planCmd)
}
