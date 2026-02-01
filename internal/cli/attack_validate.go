package cli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var validateFile string
var applyValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate an attack YAML file",
	Long: `Validate an attack YAML file without running the load test.

This checks:
- YAML syntax
- Required fields
- Value constraints and relationships`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if validateFile == "" {
			return fmt.Errorf("config file is required (use -f or --file)")
		}

		cfg, err := loadAttackFromYAML(validateFile)
		if err != nil {
			return err
		}

		_, _, _, _, err = cfg.LoadAndValidateAttack(validateFile)
		if err != nil {
			return err
		}

		if set != "" {
			fmt.Printf("attack file is ok for --set %s \n", set)
		}
		fmt.Println("validation succesfully; file is ready to be applied")
		return nil
	},
}

func init() {
	attackCmd.AddCommand(applyValidateCmd)

	applyValidateCmd.Flags().StringVarP(&validateFile, "file", "f", "", "Path to the attack YAML file")
	applyValidateCmd.Flags().StringVar(&set, "set", "", "override values then apply")
}

func (cfg *Attack) isValidURL() error {
	u, err := url.ParseRequestURI(cfg.AttackURL)
	if err != nil {
		return fmt.Errorf("invalid url syntax: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported url scheme: %s", u.Scheme)
	}

	if u.Host == "" {
		return fmt.Errorf("url must include host")
	}

	return nil
}
