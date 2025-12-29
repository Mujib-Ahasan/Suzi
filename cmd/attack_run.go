package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Mujib-Ahasan/Suzi/attacks"

	"github.com/spf13/cobra"
)

var (
	attackURL         string
	attackRate        int
	attackReq         int
	attackTimeout     time.Duration
	attackType        string
	attackMethod      string
	attackBody        string
	attackBodyFile    string
	attackContentType string
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
		var payload []byte
		attackMethod = strings.ToUpper(attackMethod)

		err := ValidateMethod(attackMethod)
		if err != nil {
			return err
		}

		if attackBody != "" && attackBodyFile != "" {
			return fmt.Errorf("both attackBody and attackBodyFile cannot be set")
		}

		if (attackMethod == "PUT" || attackMethod == "POST") && (attackBody == "" && attackBodyFile == "") {
			return fmt.Errorf("%s requires a Body", attackMethod)
		}

		if attackMethod != "GET" && attackBody != "" {
			payload = []byte(attackBody)
		}

		if attackMethod != "GET" && attackBodyFile != "" {
			b, err := os.ReadFile(attackBodyFile)
			if err != nil {
				return fmt.Errorf("failed to read Body file: %w", err)
			}
			payload = b
		}
		if len(payload) > 0 && attackMethod != "GET" {
			err := ValidateJSON(payload)
			fmt.Printf("body: %s ", string(payload))
			if err != nil {
				return fmt.Errorf("JSON provided is not valid")
			}
		}

		attackContentType, err := ValidateContentType(attackContentType)
		if err != nil {
			return err
		}

		opts := attacks.Options{
			URL:         attackURL,
			Rate:        attackRate,
			Requests:    attackReq,
			Timeout:     attackTimeout,
			Type:        attackType,
			Method:      attackMethod,
			Body:        payload,
			ContentType: attackContentType,
		}
		fmt.Printf("Starting attack: \n")
		fmt.Printf("  URL: %s\n", opts.URL)
		fmt.Printf("  Rate: %d RPS\n", opts.Rate)
		fmt.Printf("  Requests: %d\n", opts.Requests)
		fmt.Printf("  Timeout: %s\n", opts.Timeout)
		fmt.Printf("  Attack Type: %s\n", opts.Type)
		fmt.Printf("  Method: %s\n", opts.Method)
		fmt.Printf(" Body: %s \n", opts.Body)
		fmt.Printf("content type: %s \n", opts.ContentType)

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
	attackRunCmd.Flags().StringVar(&attackType, "attack-type", "constant", "Type of attack (basic/burst/rampup/random)")
	attackRunCmd.Flags().StringVar(&attackMethod, "method", "GET", "HTTP method (GET, POST, PUT, DELETE, etc.)")
	attackRunCmd.Flags().StringVar(&attackBody, "body", "", "Inline POST body")
	attackRunCmd.Flags().StringVar(&attackBodyFile, "body-file", "", "Read POST body from file")
	attackRunCmd.Flags().StringVar(&attackContentType, "content-type", "", "Attack content type")

	attackRunCmd.MarkFlagRequired("url")

	// Optional: Add flag auto-completions
	attackRunCmd.RegisterFlagCompletionFunc("attack-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"basic", "burst", "rampup", "random"}, cobra.ShellCompDirectiveNoFileComp
	})

	attackRunCmd.RegisterFlagCompletionFunc("method", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"}, cobra.ShellCompDirectiveNoFileComp
	})
}

func ValidateJSON(data []byte) error {
	var js json.RawMessage
	return json.Unmarshal(data, &js)
}

func ValidateMethod(method string) error {
	methods := [...]string{"POST", "GET", "PUT", "DELETE", "PATCH"}
	for _, e := range methods {
		if e == method {
			return nil
		}
	}
	return fmt.Errorf("Error!! Invalid HTTP Method! Valid: POST, GET, PUT, DELETE, PATCH")
}
