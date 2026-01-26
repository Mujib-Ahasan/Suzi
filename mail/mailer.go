package mail

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
	FromEmail   string        // "Suzi <no-reply@domain.com>"
	UseTLS      bool          // true for SMTPS(465) or STARTTLS as needed
	DialTimeout time.Duration // connect timeout
	SendTimeout time.Duration // send timeout
	Retries     int           // e.g. 3
	ToEmail     string
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

	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)

	e := email.NewEmail()

	e.From = m.cfg.FromEmail
	e.To = to
	e.Subject = subject
	if textBody != "" {
		e.Text = []byte(textBody)
	}
	if htmlBody != "" {
		e.HTML = []byte(htmlBody)
	}

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

	backoff := 500 * time.Millisecond
	attempts := m.cfg.Retries
	if attempts < 1 {
		attempts = 1
	}

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
