package mail

import (
	"context"
	"fmt"
	"strings"
)

func (cfg *Config) SendMail(emailTo string, reportHTML string) error {
	m := New(*cfg)
	html := fmt.Sprintf(reportHTML)

	// Plain text fallback (optional)
	text := "Suzi Load Test Report:\n"

	var recipients []string
	if emailTo != "" {
		for _, s := range strings.Split(emailTo, ",") {
			recipients = append(recipients, strings.TrimSpace(s))
		}
	} else {
		return fmt.Errorf("recipants mail must be mentioned")
	}

	ctx := context.Background()
	if err := m.Send(ctx, recipients, "Suzi Load Test Report", html, text); err != nil {
		return fmt.Errorf("email send failed: %w ", err)
	} else {
		fmt.Println("email sent successfully!!!")
	}

	return nil

}
