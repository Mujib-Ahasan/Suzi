package cli

import "time"

type Attack struct {
	AttackURL         string        `yaml:"url"`
	AttackRate        int           `yaml:"rate"`
	AttackReq         int           `yaml:"requests"`
	AttackTimeout     time.Duration `yaml:"timeout,omitempty"`
	AttackType        string        `yaml:"attack-type"`
	AttackMethod      string        `yaml:"method"`
	AttackBody        string        `yaml:"body,omitempty"`
	AttackBodyFile    string        `yaml:"body-file,omitempty"`
	AttackContentType string        `yaml:"content-type,omitempty"`
	Header            string        `yaml:"header,omitempty"`
	EmailBool         bool          `yaml:"emailEnable"`
	Email             *EmailConfig  `yaml:"email,omitempty"`
	EmailTo           string        `yaml:"-"`
	SmtpHost          string        `yaml:"-"`
	SmtpPort          int           `yaml:"-"`
	SmtpUser          string        `yaml:"-"`
	SmtpPass          string        `yaml:"-"`
	EmailFrom         string        `yaml:"-"`
	SmtpTLS           bool          `yaml:"-"`
	SmtpRetries       int           `yaml:"-"`
	SmtpTimeoutS      int           `yaml:"-"`
}

type EmailConfig struct {
	To   string     `yaml:"to"`
	From string     `yaml:"from"`
	SMTP SMTPConfig `yaml:"smtp,omitempty"`
}

type SMTPConfig struct {
	Host           string `yaml:"host,omitempty"`
	Port           int    `yaml:"port,omitempty"`
	User           string `yaml:"user,omitempty"`
	Pass           string `yaml:"pass,omitempty"`
	TLS            bool   `yaml:"tls,omitempty"`
	Retries        int    `yaml:"retries,omitempty"`
	TimeoutSeconds int    `yaml:"timeout_seconds,omitempty"`
}
