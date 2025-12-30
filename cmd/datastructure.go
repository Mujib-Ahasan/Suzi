package cmd

import "time"

type Attack struct {
	AttackURL         string        `yaml:"url"`
	AttackRate        int           `yaml:"rate"`
	AttackReq         int           `yaml:"requests"`
	AttackTimeout     time.Duration `yaml:"timeout"`
	AttackType        string        `yaml:"attack-type"`
	AttackMethod      string        `yaml:"method"`
	AttackBody        string        `yaml:"body"`
	AttackBodyFile    string        `yaml:"body-file"`
	AttackContentType string        `yaml:"content-type"`
	Header            string        `yaml:"header"`
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
	SMTP SMTPConfig `yaml:"smtp"`
}

type SMTPConfig struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	User           string `yaml:"user"`
	Pass           string `yaml:"pass"`
	TLS            bool   `yaml:"tls"`
	Retries        int    `yaml:"retries"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
}
