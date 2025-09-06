package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

type PResultIn struct {
	PRes inResult
	NRes []Result
}

type plotC struct {
	results PResultIn
	attack  string
}

func main() {
	hfs := flag.NewFlagSet("suzi", flag.ExitOnError)

	url := hfs.String("url", "", "Site where you want to attack")
	numOfReq := hfs.Int("req", 10, "Number of requests to send")
	timeout := hfs.Int("timeout", 5, "Request timeout in seconds")
	attacktype := hfs.String("atk", " ", "type of attack")
	plot := hfs.Bool("plot", false, "Do ya wanna plot da test as a timeseries ?")
	emailEnable := hfs.Bool("email", false, "Send results via email")
	numCPUS := hfs.Int("cpus", runtime.NumCPU(), "Number of CPUs to use")
	method := hfs.String("method", "GET", "HTTP method to use (GET, POST, etc.)")
	rate := hfs.Int("rate", 1, "Number of requests per second")
	emailTo := hfs.String("emailTo", "you@local.test", "Comma-separated list of recipients")
	smtpHost := hfs.String("smtpHost", "localhost", "SMTP host (e.g. smtp.gmail.com)")
	smtpPort := hfs.Int("smtpPort", 1025, "SMTP port (eg: 587 or 465)")
	smtpUser := hfs.String("smtp-user", os.Getenv("SMTP_USER"), "SMTP username (default from env SMTP_USER)")
	smtpPass := hfs.String("smtp-pass", os.Getenv("SMTP_PASS"), "SMTP password/app password (default from env SMTP_PASS)")
	emailFrom := hfs.String("emailFrom", "Suzi <noreply@gmail.com>", "From header")
	smtpTLS := hfs.Bool("smtpTLS", false, "Use TLS (SMTPS/STARTTLS)")
	smtpRetries := hfs.Int("smtp-retries", 3, "Email send retries")
	smtpTimeoutS := hfs.Int("smtp-timeout", 10, "Email send timeout in seconds")

	hfs.Parse(os.Args[1:])

	runtime.GOMAXPROCS(*numCPUS)

	var pc plotC
	var attackAll []plotC

	switch strings.ToLower(*attacktype) {
	case "mailall":
		attackAll = append(attackAll, plotC{results: basicAttack(*url, *numOfReq, *rate, *method, *timeout), attack: "Basic"})
		attackAll = append(attackAll, plotC{results: burstAttack(*url, *numOfReq, *method, *timeout), attack: "Burst"})
		attackAll = append(attackAll, plotC{results: randomLoadAttack(*url, *numOfReq, *method, *rate, *timeout), attack: "random"})
		attackAll = append(attackAll, plotC{results: rampUpAttack(*url, *numOfReq, 1, 15, *method, *timeout), attack: "Ramp-Up"})
	case "basic":
		pc = plotC{results: basicAttack(*url, *numOfReq, *rate, *method, *timeout), attack: "Basic"}
	case "burst":
		pc = plotC{results: burstAttack(*url, *numOfReq, *method, *timeout), attack: "Burst"}
	case "random":
		pc = plotC{results: randomLoadAttack(*url, *numOfReq, *method, *rate, *timeout), attack: "random"}
	case "rampup":
		pc = plotC{results: rampUpAttack(*url, *numOfReq, 1, 15, *method, *timeout), attack: "Ramp-Up"}
	default:
		fmt.Println("Unknown attack type:", *attacktype)
		return
	}

	if *plot {
		plotResults(pc)
	}

	cfg := Config{
		Host:        *smtpHost,
		Port:        *smtpPort,
		Username:    *smtpUser,
		Password:    *smtpPass,
		From:        *emailFrom,
		UseTLS:      *smtpTLS,
		DialTimeout: 5 * time.Second,
		SendTimeout: time.Duration(*smtpTimeoutS) * time.Second,
		Retries:     *smtpRetries,
	}

	if *emailEnable {
		reportHTML := BuildEmailReportHTML(attackAll, *url)
		cfg.sendMail(*emailTo, reportHTML)
	}
}
