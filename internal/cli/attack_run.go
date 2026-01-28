package cli

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/Mujib-Ahasan/Suzi/attacks"
	"github.com/Mujib-Ahasan/Suzi/internal/report"

	"github.com/spf13/cobra"
)

var currentAttack Attack

var attackRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an HTTP attack with configurable parameters",
	Long: `Execute HTTP load attacks using different strategies.

Supported attack types:
  - basic
  - burst
  - random
  - rampup

Example:
  yourproject attack --url https://example.com --rate 100 --req 1000 --attack-type constant`,
	RunE: func(cmd *cobra.Command, args []string) error {

		currentAttack.AttackMethod = strings.ToUpper(currentAttack.AttackMethod)

		err := currentAttack.ValidateMethod()
		if err != nil {
			return err
		}
		format = strings.ToLower(format)
		if format != "text" && format != "json" && format != "yaml" {
			return fmt.Errorf(
				"invalid format %q: supported formats are json and yaml",
				format,
			)
		}
		err = currentAttack.ValidateBeforeAttack()
		if err != nil {
			return err
		}

		payload, err := currentAttack.ValidateBody()
		if err != nil {
			return err
		}

		currentAttack.AttackContentType, err = ValidateContentType(currentAttack.AttackContentType)
		if err != nil {
			return err
		}

		header := currentAttack.ValidateHeader()

		opts := attacks.Options{
			URL:          currentAttack.AttackURL,
			Rate:         currentAttack.AttackRate,
			Requests:     currentAttack.AttackReq,
			Timeout:      currentAttack.AttackTimeout,
			Type:         currentAttack.AttackType,
			Method:       currentAttack.AttackMethod,
			Body:         payload,
			ContentType:  currentAttack.AttackContentType,
			Headers:      header,
			EmailEnabled: false,
		}

		if verbose {
			report.PrintVerboseHeader(opts)
		}
		slog.Debug("Preparing for attack")
		attackList := attacks.Run(opts)

		// 1️⃣ Convert to unified result
		result := report.FromAttackList(attackList)
		// 2️⃣ Output decision
		switch format {
		case "json":
			if output != "" {
				if err := report.WriteJSON(result, output); err != nil {
					return err
				}
			} else {
				if err := report.WriteJSONToStdout(result); err != nil {
					return err
				}
			}
		case "yaml":
			if output != "" {
				if err := report.WriteYAML(result, output); err != nil {
					return err
				}
			} else {
				if err := report.WriteYAMLToStdout(result); err != nil {
					return err
				}
			}
			if output != "" {

			}
		case "text":
			switch {
			case quiet:
				report.PrintQuiet(attackList)
			case verbose:
				report.PrintVerbose(attackList)
			default:
				report.PrintDefault(attackList)
			}
		}
		// Error should be checked at last or would not get the report.
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
	attackRunCmd.Flags().StringVar(&currentAttack.AttackType, "attack-type", "basic", "Type of attack (basic/burst/rampup/random)")
	attackRunCmd.Flags().StringVar(&currentAttack.AttackMethod, "method", "GET", "HTTP method (GET, POST, PUT, DELETE, etc.)")
	attackRunCmd.Flags().StringVar(&currentAttack.AttackBody, "body", "", "Inline POST body")
	attackRunCmd.Flags().StringVar(&currentAttack.AttackBodyFile, "body-file", "", "Read POST body from file")
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
	slog.Debug("Validating JSON")
	var js json.RawMessage
	return json.Unmarshal(data, &js)
}

func (cAttack Attack) ValidateMethod() error {
	slog.Debug("Validating attacking method")
	methods := [...]string{"POST", "GET", "PUT", "DELETE", "PATCH"}
	for _, e := range methods {
		if e == cAttack.AttackMethod {
			return nil
		}
	}
	return fmt.Errorf("Error!! Invalid HTTP Method! Valid: POST, GET, PUT, DELETE, PATCH")
}

func (cAttack Attack) ValidateBody() ([]byte, error) {
	slog.Debug("Validating body")
	var payload []byte
	if cAttack.AttackMethod != "GET" && cAttack.AttackBody != "" {
		payload = []byte(cAttack.AttackBody)
	}

	if cAttack.AttackMethod != "GET" && cAttack.AttackBodyFile != "" {
		b, err := os.ReadFile(cAttack.AttackBodyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read Body file: %w", err)
		}
		payload = b
	}
	if len(payload) > 0 && cAttack.AttackMethod != "GET" {
		err := cAttack.ValidateJSON(payload)
		if err != nil {
			return nil, fmt.Errorf("JSON provided is not valid")
		}
	}

	return payload, nil
}

func (cAttack Attack) ValidateBeforeAttack() error {

	if cAttack.AttackBody != "" && cAttack.AttackBodyFile != "" {
		return fmt.Errorf("both attackBody and attackBodyFile cannot be set")
	}

	if (cAttack.AttackMethod == "PUT" || cAttack.AttackMethod == "POST" || cAttack.AttackMethod == "PATCH") && (cAttack.AttackBody == "" && cAttack.AttackBodyFile == "") {
		return fmt.Errorf("%s requires a Body", cAttack.AttackMethod)
	}

	return nil

}
