package cmd

import (
	"fmt"
	"time"

	"github.com/Mujib-Ahasan/Suzi/attacks"

	"github.com/spf13/cobra"
)

var (
	attackURL     string
	attackRate    int
	attackReq     int
	attackTimeout time.Duration
	attackType    string
	attackMethod  string
)

var attackRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an HTTP attack with configurable parameters",
	Long: `Execute HTTP load attacks using different strategies.

Supported attack types:
  - constant
  - ramp
  - spike
  - custom patterns

Example:
  yourproject attack --url https://example.com --rate 100 --req 1000 --attack-type constant`,
	RunE: func(cmd *cobra.Command, args []string) error {

		opts := attacks.Options{
			URL:      attackURL,
			Rate:     attackRate,
			Requests: attackReq,
			Timeout:  attackTimeout,
			Type:     attackType,
			Method:   attackMethod,
		}

		fmt.Println("Starting attack:")
		fmt.Printf("  URL: %s\n", opts.URL)
		fmt.Printf("  Rate: %d RPS\n", opts.Rate)
		fmt.Printf("  Requests: %d\n", opts.Requests)
		fmt.Printf("  Timeout: %s\n", opts.Timeout)
		fmt.Printf("  Attack Type: %s\n", opts.Type)
		fmt.Printf("  Method: %s\n", opts.Method)
		attackList := attacks.Run(opts, false)

		if err := attackList[0].Results.Err; err != nil {
			return fmt.Errorf("error: %v ", err)
		}

		return nil
	},
}

func init() {
	attackCmd.AddCommand(attackRunCmd)

	attackRunCmd.Flags().StringVar(&attackURL, "url", "", "Target URL for the attack (required)")
	attackRunCmd.Flags().IntVar(&attackRate, "rate", 10, "Requests per second")
	attackRunCmd.Flags().IntVar(&attackReq, "req", 100, "Total number of requests to send")
	attackRunCmd.Flags().DurationVar(&attackTimeout, "timeout", 5*time.Second, "Request timeout duration")
	attackRunCmd.Flags().StringVar(&attackType, "attack-type", "constant", "Type of attack (constant/ramp/spike)")
	attackRunCmd.Flags().StringVar(&attackMethod, "method", "GET", "HTTP method (GET, POST, PUT, DELETE, etc.)")

	attackRunCmd.MarkFlagRequired("url")

	// Optional: Add flag auto-completions
	attackRunCmd.RegisterFlagCompletionFunc("attack-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"basic", "burst", "rampup", "random"}, cobra.ShellCompDirectiveNoFileComp
	})

	attackRunCmd.RegisterFlagCompletionFunc("method", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"}, cobra.ShellCompDirectiveNoFileComp
	})
}
