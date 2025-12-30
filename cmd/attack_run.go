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

var currentAttack Attack

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

		currentAttack.AttackMethod = strings.ToUpper(currentAttack.AttackMethod)

		err := currentAttack.ValidateMethod()
		if err != nil {
			return err
		}

		err = currentAttack.ValidateBeforeAttack()
		if err != nil {
			return err
		}

		payload, err := currentAttack.ValidateBody()
		if err != nil {
			return err
		}

		attackContentType, err := ValidateContentType(currentAttack.AttackContentType)
		if err != nil {
			return err
		}

		header := currentAttack.ValidateHeader()

		opts := attacks.Options{
			URL:         currentAttack.AttackURL,
			Rate:        currentAttack.AttackRate,
			Requests:    currentAttack.AttackReq,
			Timeout:     currentAttack.AttackTimeout,
			Type:        currentAttack.AttackType,
			Method:      currentAttack.AttackMethod,
			Body:        payload,
			ContentType: attackContentType,
			Headers:     header,
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

	attackRunCmd.Flags().StringVar(&currentAttack.AttackURL, "url", "", "Target URL for the attack (required)")
	attackRunCmd.Flags().IntVar(&currentAttack.AttackRate, "rate", 10, "Requests per second")
	attackRunCmd.Flags().IntVar(&currentAttack.AttackReq, "req", 100, "Total number of requests to send")
	attackRunCmd.Flags().DurationVar(&currentAttack.AttackTimeout, "timeout", 5*time.Second, "Request timeout duration")
	attackRunCmd.Flags().StringVar(&currentAttack.AttackType, "attack-type", "constant", "Type of attack (basic/burst/rampup/random)")
	attackRunCmd.Flags().StringVar(&currentAttack.AttackMethod, "method", "GET", "HTTP method (GET, POST, PUT, DELETE, etc.)")
	attackRunCmd.Flags().StringVar(&currentAttack.AttackBody, "body", "", "Inline POST body")
	attackRunCmd.Flags().StringVar(&currentAttack.AttackBodyFile, "body-file", "", "Read POST body from file")
	attackRunCmd.Flags().StringVar(&currentAttack.AttackContentType, "content-type", "", "Attack content type")
	attackRunCmd.Flags().StringVar(&currentAttack.AttackContentType, "content-type", "", "Attack content type")
	attackRunCmd.Flags().StringVar(&currentAttack.Header, "header", "", "Headers for the request")

	attackRunCmd.MarkFlagRequired("url")

	// Optional: Add flag auto-completions
	attackRunCmd.RegisterFlagCompletionFunc("attack-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"basic", "burst", "rampup", "random"}, cobra.ShellCompDirectiveNoFileComp
	})

	attackRunCmd.RegisterFlagCompletionFunc("method", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"}, cobra.ShellCompDirectiveNoFileComp
	})
}

func (cAttack Attack) ValidateJSON(data []byte) error {
	var js json.RawMessage
	return json.Unmarshal(data, &js)
}

func (cAttack Attack) ValidateMethod() error {
	methods := [...]string{"POST", "GET", "PUT", "DELETE", "PATCH"}
	for _, e := range methods {
		if e == cAttack.AttackMethod {
			return nil
		}
	}
	return fmt.Errorf("Error!! Invalid HTTP Method! Valid: POST, GET, PUT, DELETE, PATCH")
}

func (cAttack Attack) ValidateBody() ([]byte, error) {
	var payload []byte
	if currentAttack.AttackMethod != "GET" && currentAttack.AttackBody != "" {
		payload = []byte(currentAttack.AttackBody)
	}

	if currentAttack.AttackMethod != "GET" && currentAttack.AttackBodyFile != "" {
		b, err := os.ReadFile(currentAttack.AttackBodyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read Body file: %w", err)
		}
		payload = b
	}
	if len(payload) > 0 && currentAttack.AttackMethod != "GET" {
		err := currentAttack.ValidateJSON(payload)
		fmt.Printf("body: %s ", string(payload))
		if err != nil {
			return nil, fmt.Errorf("JSON provided is not valid")
		}
	}

	return payload, nil
}

func (cAttack Attack) ValidateBeforeAttack() error {

	if currentAttack.AttackBody != "" && currentAttack.AttackBodyFile != "" {
		return fmt.Errorf("both attackBody and attackBodyFile cannot be set")
	}

	if (currentAttack.AttackMethod == "PUT" || currentAttack.AttackMethod == "POST" || currentAttack.AttackMethod == "PATCH") && (currentAttack.AttackBody == "" && currentAttack.AttackBodyFile == "") {
		return fmt.Errorf("%s requires a Body", currentAttack.AttackMethod)
	}

	return nil

}
