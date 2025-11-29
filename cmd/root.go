package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose bool
)

// rootCmd is the base command for the CLI.
var rootCmd = &cobra.Command{
	Use:   "suzi",
	Short: "A high-performance HTTP load-testing tool",
	Long: `Suzi is a powerful and extensible HTTP load-testing tool.

It supports multiple attack strategies (constant, ramp, burst, custom patterns),
result summarization, and integrations (mail, dashboards).`,

	// Root command should NOT run any attack directly.
	// It should only show help if called without subcommands.
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute is called by main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func init() {
	// Global flags available to all subcommands
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}
