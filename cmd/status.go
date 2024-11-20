package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yahao333/gort/internal/state"
)

var statusCmd = &cobra.Command{
	Use:   "status [environment]",
	Short: "Show the current status of an environment",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var env string
		if len(args) > 0 {
			env = args[0]
		}

		sm := state.NewStateManager(".")
		state, err := sm.LoadState(env)
		if err != nil {
			return fmt.Errorf("failed to load state: %w", err)
		}

		// Print status
		fmt.Printf("Environment: %s\n", state.Environment)
		fmt.Printf("Last Update: %s\n", state.LastUpdate)
		fmt.Printf("Resources: %d\n", len(state.Resources))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
