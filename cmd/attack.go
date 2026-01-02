package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var quiet bool
var verbose bool

var attackCmd = &cobra.Command{
	Use:   "attack",
	Short: "HTTP load testing related subcommands",
	Long:  "Run HTTP load tests and send reports.",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if quiet && verbose {
			return fmt.Errorf("cannot use --quiet and --verbose together")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(attackCmd)

	attackCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "machine-friendly output")
	attackCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "human-friendly output")
}
