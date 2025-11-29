package cmd

import "github.com/spf13/cobra"

var attackCmd = &cobra.Command{
	Use:   "attack",
	Short: "HTTP load testing related subcommands",
	Long:  "Run HTTP load tests and send reports.",
}

func init() {
	rootCmd.AddCommand(attackCmd)

}
