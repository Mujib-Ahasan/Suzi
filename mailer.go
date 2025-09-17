package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
)

type Config struct {
	Host        string        // e.g. "smtp.gmail.com"
	Port        int           // e.g. 587 or 465
	Username    string        // SMTP username (email/user)
	Password    string        // SMTP password or app password
	From        string        // "Suzi <no-reply@domain.com>"
	UseTLS      bool          // true for SMTPS(465) or STARTTLS as needed
	DialTimeout time.Duration // connect timeout
	SendTimeout time.Duration // send timeout
	Retries     int           // e.g. 3
}

// // Attachment file path
// type Attachment struct {
// 	Path string
// 	Name string // optional; if empty, the file name is used
// }

type Mailer struct {
	cfg Config
}

func New(cfg Config) *Mailer { return &Mailer{cfg: cfg} }

func (m *Mailer) Send(ctx context.Context, to []string, subject, htmlBody, textBody string) error {
	// fmt.Println("debug-mailer_1")

	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)
	// fmt.Println("debug-mailer_2")

	e := email.NewEmail()
	// fmt.Println("debug-mailer_3")

	e.From = m.cfg.From
	e.To = to
	e.Subject = subject
	if textBody != "" {
		e.Text = []byte(textBody)
	}
	if htmlBody != "" {
		e.HTML = []byte(htmlBody)
	}

	// fmt.Println("debug-mailer_4")

	// for _, a := range atts {
	// 	if a.Name != "" {
	// 		f, err := os.Open(a.Path)
	// 		if err != nil {
	// 			return fmt.Errorf("attach: %w", err)
	// 		}
	// 		defer f.Close()

	// 		if _, err := e.Attach(f, a.Name, a.Name); err != nil {
	// 			return fmt.Errorf("attach: %w", err)
	// 		}
	// 	} else {
	// 		// normal attach (uses the actual filename)
	// 		if _, err := e.AttachFile(a.Path); err != nil {
	// 			return fmt.Errorf("attach: %w", err)
	// 		}
	// 	}
	// }
	// fmt.Println("debug-mailer_5")

	var (
		auth smtp.Auth
		tlsC *tls.Config
	)
	if m.cfg.Username != "" || m.cfg.Password != "" {
		auth = smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	}
	if m.cfg.UseTLS {
		tlsC = &tls.Config{ServerName: m.cfg.Host}
	}

	// fmt.Println("debug-mailer_6")

	backoff := 500 * time.Millisecond
	attempts := m.cfg.Retries
	if attempts < 1 {
		attempts = 1
	}

	// fmt.Println("debug-mailer_7")
	var lastErr error
	for i := 0; i < attempts; i++ {
		sendCtx, cancel := context.WithTimeout(ctx, m.cfg.SendTimeout)
		defer cancel()

		errCh := make(chan error, 1)
		go func() {
			if m.cfg.UseTLS {
				errCh <- e.SendWithTLS(addr, auth, tlsC)
			} else {
				errCh <- e.Send(addr, auth)
			}
		}()

		select {
		case <-sendCtx.Done():
			lastErr = fmt.Errorf("send timeout: %w", sendCtx.Err())
		case err := <-errCh:
			if err == nil {
				return nil
			}
			lastErr = err
		}
		time.Sleep(backoff)
		backoff *= 2
	}
	return lastErr
}
