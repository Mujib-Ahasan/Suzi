package cli

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var attackGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a yaml file from attack configuration",
	Long: `Generate a yaml file from attack configuration.

This is not like ordinary attack command but generate a attack yaml
file from the attack configuration`,
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

		currentAttack.AttackContentType, err = ValidateContentType(currentAttack.AttackContentType)
		if err != nil {
			return err
		}
		var payload []byte
		if currentAttack.AttackBodyFile != "" || currentAttack.AttackBody != "" {
			payload, err = currentAttack.ValidateBody()
			currentAttack.AttackBody = string(payload)
			currentAttack.AttackBodyFile = ""
		}

		if currentAttack.EmailBool {
			currentAttack.Email = &EmailConfig{}
			currentAttack.Email.To = currentAttack.EmailTo
			currentAttack.Email.From = currentAttack.EmailFrom
			currentAttack.Email.SMTP.Host = currentAttack.SmtpHost
			currentAttack.Email.SMTP.Port = currentAttack.SmtpPort
			currentAttack.Email.SMTP.Retries = currentAttack.SmtpRetries
			currentAttack.Email.SMTP.TimeoutSeconds = currentAttack.SmtpTimeoutS
			currentAttack.Email.SMTP.TLS = currentAttack.SmtpTLS

			err = currentAttack.Email.Validate()
			if err != nil {
				return err
			}
		}

		if output == "" {
			if err := WriteYAMLToStdout(currentAttack); err != nil {
				return err
			}
		} else {
			if err := WriteToYAML(currentAttack, output); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	attackCmd.AddCommand(attackGenerateCmd)

	attackGenerateCmd.Flags().StringVar(&currentAttack.AttackURL, "url", "", "Target URL for the attack (required)")
	attackGenerateCmd.Flags().IntVar(&currentAttack.AttackRate, "rate", 5, "Requests per second")
	attackGenerateCmd.Flags().IntVar(&currentAttack.AttackReq, "req", 25, "Total number of requests to send")
	attackGenerateCmd.Flags().StringVar(&currentAttack.AttackMethod, "method", "GET", "HTTP method (GET, POST, PUT, DELETE, etc.)")
	attackGenerateCmd.Flags().StringVar(&currentAttack.AttackBody, "body", "", "Inline POST body")
	attackGenerateCmd.Flags().StringVar(&currentAttack.AttackBodyFile, "body-file", "", "Read POST body from file")
	attackGenerateCmd.Flags().StringVar(&currentAttack.AttackContentType, "content-type", "", "Attack content type")
	attackGenerateCmd.Flags().StringVar(&currentAttack.Header, "header", "", "Headers for the request")
	attackGenerateCmd.Flags().StringVar(&currentAttack.AttackType, "attack-type", "basic", "Type of attack (basic/burst/rampup/random)")
	attackGenerateCmd.Flags().BoolVar(&currentAttack.EmailBool, "email", false, "Email Enable or not")

	attackGenerateCmd.Flags().StringVar(&currentAttack.EmailTo, "emailto", "you@local.test", "Comma-separated list of recipients")
	attackGenerateCmd.Flags().StringVar(&currentAttack.SmtpHost, "smtpHost", "localhost", "SMTP host (e.g. smtp.gmail.com)")
	attackGenerateCmd.Flags().IntVar(&currentAttack.SmtpPort, "smtpPort", 1025, "SMTP port (e.g. 587 or 465)")
	attackGenerateCmd.Flags().StringVar(&currentAttack.SmtpUser, "smtp-user", os.Getenv("SMTP_USER"), "SMTP username (default from env SMTP_USER)")
	attackGenerateCmd.Flags().StringVar(&currentAttack.SmtpPass, "smtp-pass", os.Getenv("SMTP_PASS"), "SMTP password/app password (default from env SMTP_PASS)")
	attackGenerateCmd.Flags().StringVar(&currentAttack.EmailFrom, "emailFrom", "Suzi <noreply@gmail.com>", "From header")
	attackGenerateCmd.Flags().BoolVar(&currentAttack.SmtpTLS, "smtpTLS", false, "Use TLS (SMTPS/STARTTLS)")
	attackGenerateCmd.Flags().IntVar(&currentAttack.SmtpRetries, "smtp-retries", 3, "Email send retries")
	attackGenerateCmd.Flags().IntVar(&currentAttack.SmtpTimeoutS, "smtp-timeout", 10, "Email send timeout in seconds")
}
