package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var debug bool

// rootCmd is the base command for the CLI.
var rootCmd = &cobra.Command{
	Use:   "suzi",
	Short: "A high-performance HTTP load-testing tool",
	Long: `Suzi is a powerful and extensible HTTP load-testing tool.

It supports multiple attack strategies (constant, ramp, burst, custom patterns),
result summarization, and integrations (mail, dashboards).`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		setupLogger()
		return nil
	},

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
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")
}

func DebugEnabled() bool {
	return debug
}

func setupLogger() {
	level := slog.LevelInfo

	if DebugEnabled() {
		level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})

	slog.SetDefault(slog.New(handler))
}
