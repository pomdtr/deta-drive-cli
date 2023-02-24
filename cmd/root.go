/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/deta/deta-go/deta"
	"github.com/spf13/cobra"
)

type contextKey int

const (
	detaKey contextKey = iota
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "deta-drive",
	Short:        "A companion cli for Deta Drive",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		dataKey, err := cmd.Flags().GetString("data-key")
		if err != nil {
			return err
		}

		if dataKey == "" && os.Getenv("DETA_PROJECT_KEY") == "" {
			return fmt.Errorf("data-key is required if DETA_PROJECT_KEY is not set")
		}

		options := []deta.ConfigOption{}
		if dataKey != "" {
			options = append(options, deta.WithProjectKey(dataKey))
		}

		deta, err := deta.New(options...)
		if err != nil {
			return err
		}

		cmd.SetContext(context.WithValue(cmd.Context(), detaKey, deta))
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("data-key", "", "Data key for the drive")
}
