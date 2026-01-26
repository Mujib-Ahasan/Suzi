package cli

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/Mujib-Ahasan/Suzi/attacks"
	"github.com/Mujib-Ahasan/Suzi/internal/report"
	ml "github.com/Mujib-Ahasan/Suzi/mail"

	"github.com/spf13/cobra"
)

func init() {
	// Register this subcommand under rootCmd
	attackCmd.AddCommand(attackEmailCmd)
	// Flags for the email command
	attackEmailCmd.Flags().StringVar(&currentAttack.AttackURL, "url", "", "Target URL for the attack (required)")
	attackEmailCmd.Flags().IntVar(&currentAttack.AttackRate, "rate", 5, "Requests per second")
	attackEmailCmd.Flags().IntVar(&currentAttack.AttackReq, "req", 25, "Total number of requests to send")
	attackEmailCmd.Flags().StringVar(&currentAttack.AttackMethod, "method", "GET", "HTTP method (GET, POST, PUT, DELETE, etc.)")
	attackEmailCmd.Flags().StringVar(&currentAttack.AttackBody, "body", "", "Inline POST body")
	attackEmailCmd.Flags().StringVar(&currentAttack.AttackBodyFile, "body-file", "", "Read POST body from file")
	attackEmailCmd.Flags().StringVar(&currentAttack.AttackContentType, "content-type", "", "Attack content type")
	attackEmailCmd.Flags().StringVar(&currentAttack.Header, "header", "", "Headers for the request")

	attackEmailCmd.Flags().StringVar(&currentAttack.EmailTo, "emailTo", "you@local.test", "Comma-separated list of recipients")
	attackEmailCmd.Flags().StringVar(&currentAttack.SmtpHost, "smtpHost", "localhost", "SMTP host (e.g. smtp.gmail.com)")
	attackEmailCmd.Flags().IntVar(&currentAttack.SmtpPort, "smtpPort", 1025, "SMTP port (e.g. 587 or 465)")
	attackEmailCmd.Flags().StringVar(&currentAttack.SmtpUser, "smtp-user", os.Getenv("SMTP_USER"), "SMTP username (default from env SMTP_USER)")
	attackEmailCmd.Flags().StringVar(&currentAttack.SmtpPass, "smtp-pass", os.Getenv("SMTP_PASS"), "SMTP password/app password (default from env SMTP_PASS)")
	attackEmailCmd.Flags().StringVar(&currentAttack.EmailFrom, "emailFrom", "Suzi <noreply@gmail.com>", "From header")
	attackEmailCmd.Flags().BoolVar(&currentAttack.SmtpTLS, "smtpTLS", false, "Use TLS (SMTPS/STARTTLS)")
	attackEmailCmd.Flags().IntVar(&currentAttack.SmtpRetries, "smtp-retries", 3, "Email send retries")
	attackEmailCmd.Flags().IntVar(&currentAttack.SmtpTimeoutS, "smtp-timeout", 10, "Email send timeout in seconds")
}

var attackEmailCmd = &cobra.Command{
	Use:   "email",
	Short: "Send an email report",
	Long:  "Send an email report using SMTP settings.",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentAttack.AttackMethod = strings.ToUpper(currentAttack.AttackMethod)

		err := currentAttack.ValidateMethod()
		if err != nil {
			return err
		}
		format = strings.ToLower(format)
		if format != "" && format != "json" && format != "yaml" {
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

		attackContentType, err := ValidateContentType(currentAttack.AttackContentType)
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
			ContentType:  attackContentType,
			Headers:      header,
			EmailEnabled: true,
		}

		if verbose {
			report.PrintVerboseHeader(opts)
		}
		slog.Debug("Preparing for Email attack")

		attackList := attacks.Run(opts, true)
		if err := attackList[0].Results.Err; err != nil {
			return fmt.Errorf("error: %v ", err)
		}

		cfg := ml.Config{
			Host:        currentAttack.SmtpHost,
			Port:        currentAttack.SmtpPort,
			Username:    currentAttack.SmtpUser,
			Password:    currentAttack.SmtpPass,
			FromEmail:   currentAttack.EmailFrom,
			UseTLS:      currentAttack.SmtpTLS,
			DialTimeout: 5 * time.Second,
			SendTimeout: time.Duration(currentAttack.SmtpTimeoutS) * time.Second,
			Retries:     currentAttack.SmtpRetries,
			ToEmail:     currentAttack.EmailTo,
		}

		reportHTML := ml.BuildEmailReportHTML(attackList, opts.URL)
		if err := cfg.SendMail(currentAttack.EmailTo, reportHTML); err != nil {
			return fmt.Errorf("error: %w", err)
		}

		// 1️⃣ Convert to unified result
		result := report.FromAttackListEmail(attackList, cfg)
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

		return nil
	},
}

func ValidateContentType(ct string) (string, error) {
	if ct == "" {
		return "application/json", nil
	}
	contentTypes := [...]string{"application/javascript ", "application/json", "application/xml"}

	for _, e := range contentTypes {
		if e == ct {
			return e, nil
		}
	}

	return "", fmt.Errorf("Error!!!please choose between application/javascript, application/json, application/xml")
}

func (cAttack Attack) ValidateHeader() map[string]string {
	slog.Debug("Validating header")
	headers := make(map[string]string)
	if strings.TrimSpace(cAttack.Header) == "" {
		return headers
	}
	pairs := strings.Split(cAttack.Header, ",")

	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)

		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if key != "" {
				headers[key] = value
			}
		}
	}

	return headers
}
