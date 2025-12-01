package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Mujib-Ahasan/Suzi/attacks"
	ml "github.com/Mujib-Ahasan/Suzi/mail"

	"github.com/spf13/cobra"
)

var (
	emailTo      string
	smtpHost     string
	smtpPort     int
	smtpUser     string
	smtpPass     string
	emailFrom    string
	smtpTLS      bool
	smtpRetries  int
	smtpTimeoutS int
)

func init() {
	// Register this subcommand under rootCmd
	attackCmd.AddCommand(attackEmailCmd)
	// Flags for the email command
	attackEmailCmd.Flags().StringVar(&attackURL, "url", "", "Target URL for the attack (required)")
	attackEmailCmd.Flags().IntVar(&attackRate, "rate", 5, "Requests per second")
	attackEmailCmd.Flags().IntVar(&attackReq, "req", 25, "Total number of requests to send")
	attackEmailCmd.Flags().StringVar(&attackMethod, "method", "GET", "HTTP method (GET, POST, PUT, DELETE, etc.)")
	attackEmailCmd.Flags().StringVar(&attackBody, "body", "", "Inline POST body")
	attackEmailCmd.Flags().StringVar(&attackBodyFile, "body-file", "", "Read POST body from file")

	attackEmailCmd.Flags().StringVar(&emailTo, "emailTo", "you@local.test", "Comma-separated list of recipients")
	attackEmailCmd.Flags().StringVar(&smtpHost, "smtpHost", "localhost", "SMTP host (e.g. smtp.gmail.com)")
	attackEmailCmd.Flags().IntVar(&smtpPort, "smtpPort", 1025, "SMTP port (e.g. 587 or 465)")
	attackEmailCmd.Flags().StringVar(&smtpUser, "smtp-user", os.Getenv("SMTP_USER"), "SMTP username (default from env SMTP_USER)")
	attackEmailCmd.Flags().StringVar(&smtpPass, "smtp-pass", os.Getenv("SMTP_PASS"), "SMTP password/app password (default from env SMTP_PASS)")
	attackEmailCmd.Flags().StringVar(&emailFrom, "emailFrom", "Suzi <noreply@gmail.com>", "From header")
	attackEmailCmd.Flags().BoolVar(&smtpTLS, "smtpTLS", false, "Use TLS (SMTPS/STARTTLS)")
	attackEmailCmd.Flags().IntVar(&smtpRetries, "smtp-retries", 3, "Email send retries")
	attackEmailCmd.Flags().IntVar(&smtpTimeoutS, "smtp-timeout", 10, "Email send timeout in seconds")
}

var attackEmailCmd = &cobra.Command{
	Use:   "email",
	Short: "Send an email report",
	Long:  "Send an email report using SMTP settings.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var payload []byte
		attackMethod = strings.ToUpper(attackMethod)

		if attackBody != "" && attackBodyFile != "" {
			return fmt.Errorf("both attackBody and attackBodyFile cannot be set")
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
			err := validateJSON(payload)
			if err != nil {
				return fmt.Errorf("JSON provided is not valid")
			}
		}

		opts := attacks.Options{
			URL:      attackURL,
			Rate:     attackRate,
			Requests: attackReq,
			Timeout:  attackTimeout,
			Type:     attackType,
			Method:   attackMethod,
			Body:     payload,
		}

		fmt.Println("Email command invoked with:")
		fmt.Println("SMTP Host:", smtpHost)
		fmt.Println("SMTP Port:", smtpPort)
		fmt.Println("SMTP User:", smtpUser)
		fmt.Println("SMTP TLS:", smtpTLS)
		fmt.Println("Retries:", smtpRetries)
		fmt.Println("Body:", attackBody)

		attackList := attacks.Run(opts, true)

		if err := attackList[0].Results.Err; err != nil {
			return fmt.Errorf("error: %v ", err)
		}

		cfg := ml.Config{
			Host:        smtpHost,
			Port:        smtpPort,
			Username:    smtpUser,
			Password:    smtpPass,
			From:        emailFrom,
			UseTLS:      smtpTLS,
			DialTimeout: 5 * time.Second,
			SendTimeout: time.Duration(smtpTimeoutS) * time.Second,
			Retries:     smtpRetries,
		}

		reportHTML := ml.BuildEmailReportHTML(attackList, attackURL)
		if err := cfg.SendMail(emailTo, reportHTML); err != nil {
			return fmt.Errorf("error: %w", err)
		}

		return nil
	},
}
