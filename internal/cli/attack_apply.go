package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/Mujib-Ahasan/Suzi/attacks"
	"github.com/Mujib-Ahasan/Suzi/common"
	"github.com/Mujib-Ahasan/Suzi/internal/report"
	ml "github.com/Mujib-Ahasan/Suzi/mail"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var applyFile string

// applyCmd represents the attack apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply an attack configuration from a YAML file",
	Long: `Apply a load test configuration from a YAML file.

This behaves exactly like 'attack run', but uses a YAML file
as the primary configuration source. CLI flags override YAML values.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if applyFile == "" {
			return fmt.Errorf("config file is required (use -f or --file)")
		}
		format = strings.ToLower(format)
		if format != "text" && format != "json" && format != "yaml" {
			return fmt.Errorf(
				"invalid format %q: supported formats are json and yaml",
				format,
			)
		}

		err := DecodeStrictYAML(applyFile, &currentAttack)
		if err != nil {
			return err
		}

		cfg, err := loadAttackFromYAML(applyFile)
		if err != nil {
			return err
		}
		cfg.AttackMethod = strings.ToUpper(cfg.AttackMethod)

		if err := cfg.ValidateYaml(); err != nil {
			return err
		}

		err = cfg.ValidateMethod()
		if err != nil {
			return err
		}
		err = cfg.ValidateBeforeAttack()
		if err != nil {
			return err
		}
		payload, err := cfg.ValidateBody()
		if err != nil {
			return err
		}

		attackContentType, err := ValidateContentType(currentAttack.AttackContentType)
		if err != nil {
			return err
		}
		emailFlag := false

		if cfg.Email != nil {
			emailFlag = true
		}

		opts := attacks.Options{
			URL:          cfg.AttackURL,
			Rate:         cfg.AttackRate,
			Requests:     cfg.AttackReq,
			Timeout:      cfg.AttackTimeout,
			Type:         cfg.AttackType,
			Method:       cfg.AttackMethod,
			Body:         payload,
			ContentType:  attackContentType,
			EmailEnabled: emailFlag,
		}

		if verbose {
			report.PrintVerboseHeader(opts)
		}
		slog.Debug("Preparing for attack ")
		var attackList []common.PlotC
		var mailCfg ml.Config
		var result report.LoadTestResultAll

		if cfg.Email == nil {
			attackList = attacks.Run(opts)
			result = report.FromAttackList(attackList)
		} else {
			attackList = attacks.Run(opts)
			mailCfg = ml.Config{
				ToEmail:     cfg.Email.To,
				Host:        cfg.Email.SMTP.Host,
				Port:        cfg.Email.SMTP.Port,
				Username:    cfg.Email.SMTP.User,
				Password:    cfg.Email.SMTP.Pass,
				FromEmail:   cfg.Email.From,
				UseTLS:      cfg.Email.SMTP.TLS,
				DialTimeout: 5 * time.Second,
				SendTimeout: time.Duration(cfg.Email.SMTP.TimeoutSeconds) * time.Second,
				Retries:     cfg.Email.SMTP.Retries,
			}
			result = report.FromAttackListEmail(attackList, mailCfg)
		}
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

		if err := attackList[0].Results.Err; err != nil {
			return fmt.Errorf("error: %v ", err)
		}

		if cfg.Email != nil {
			reportHTML := ml.BuildEmailReportHTML(attackList, opts.URL)
			if err := mailCfg.SendMail(cfg.Email.To, reportHTML); err != nil {
				return fmt.Errorf("error: %w", err)
			}
		}
		return nil
	},
}

func init() {
	applyCmd.Flags().StringVarP(&applyFile, "file", "f", "", "Path to attack YAML file")

	attackCmd.AddCommand(applyCmd)
}

func loadAttackFromYAML(path string) (*Attack, error) {
	slog.Debug("Parsing YAML file")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Attack
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	return &cfg, nil
}

func (cAttack Attack) ValidateYaml() error {
	slog.Debug("Validating YAML file")

	if strings.TrimSpace(cAttack.AttackURL) == "" {
		return fmt.Errorf("url is required")
	}

	if cAttack.AttackRate <= 0 {
		return fmt.Errorf("rate must be greater than 0")
	}

	if cAttack.AttackReq <= 0 {
		return fmt.Errorf("requests must be greater than 0")
	}

	if cAttack.AttackMethod == "" {
		cAttack.AttackMethod = "GET"
	}

	if cAttack.AttackTimeout == 0 {
		cAttack.AttackTimeout = 5 * time.Second
	}

	method := strings.ToUpper(cAttack.AttackMethod)
	switch method {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS":
		cAttack.AttackMethod = method
	default:
		return fmt.Errorf("unsupported http method: %s", cAttack.AttackMethod)
	}

	if cAttack.AttackBody != "" && cAttack.AttackBodyFile != "" {
		return fmt.Errorf("only one of body or body-file may be set")
	}

	if (cAttack.AttackBody != "" || cAttack.AttackBodyFile != "") && method == "GET" {
		return fmt.Errorf("request body is not allowed for GET requests")
	}

	if cAttack.Email != nil {
		if err := cAttack.Email.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (e *EmailConfig) Validate() error {
	slog.Debug("Validating Email details \n")
	if strings.TrimSpace(e.To) == "" {
		return fmt.Errorf("email.to is required")
	}

	if strings.TrimSpace(e.From) == "" {
		return fmt.Errorf("email.from is required")
	}

	return e.SMTP.Validate()
}

func (s *SMTPConfig) Validate() error {
	if strings.TrimSpace(s.Host) == "" {
		return fmt.Errorf("email.smtp.host is required")
	}
	if strings.TrimSpace(s.User) == "" {
		s.User = os.Getenv("SMTP_USER")
	}

	if strings.TrimSpace(s.Pass) == "" {
		s.Pass = os.Getenv("SMTP_PASS")
	}

	if s.Port <= 0 || s.Port > 65535 {
		return fmt.Errorf("email.smtp.port must be between 1 and 65535")
	}

	if s.TimeoutSeconds <= 0 {
		s.TimeoutSeconds = 10
	}

	if s.Retries < 0 {
		return fmt.Errorf("email.smtp.retries cannot be negative")
	}

	return nil
}

func DecodeStrictYAML(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)

	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("invalid: %w", err)
	}

	// Optional: ensure only ONE document
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("multiple YAML documents are not supported")
	}

	return nil
}
