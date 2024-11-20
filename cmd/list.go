package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yahao333/gort/internal/config"
)

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List available environments",
    RunE: func(cmd *cobra.Command, args []string) error {
        cfg, err := config.LoadConfig("")
        if err != nil {
            return fmt.Errorf("failed to load config: %w", err)
        }

        fmt.Println("Available Environments:")
        fmt.Println("----------------------")
        for name, env := range cfg.Environments {
            fmt.Printf("- %s:\n", name)
            fmt.Printf("  Provider: %s\n", env.Provider)
            fmt.Printf("  Region: %s\n", env.Region)
            if len(env.Tags) > 0 {
                fmt.Printf("  Tags: %v\n", env.Tags)
            }
            fmt.Println()
        }
        
        return nil
    },
}

func init() {
    rootCmd.AddCommand(listCmd)
}