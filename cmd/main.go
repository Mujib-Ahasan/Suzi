package cmd

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	st "github.com/Mujib-Ahasan/Suzi/attacks"
	cm "github.com/Mujib-Ahasan/Suzi/common"
	pr "github.com/Mujib-Ahasan/Suzi/core"
	ml "github.com/Mujib-Ahasan/Suzi/mail"
)

func main() {
	hfs := flag.NewFlagSet("suzi", flag.ExitOnError)

	url := hfs.String("url", "", "Site where you want to attack")
	// numOfReq := hfs.Int("req", 10, "Number of requests to send")
	// timeout := hfs.Int("timeout", 5, "Request timeout in seconds")
	attacktype := hfs.String("atk", " ", "type of attack")
	plot := hfs.Bool("plot", false, "Do ya wanna plot da test as a timeseries ?")
	emailEnable := hfs.Bool("email", false, "Send results via email")
	numCPUS := hfs.Int("cpus", runtime.NumCPU(), "Number of CPUs to use")
	// method := hfs.String("method", "GET", "HTTP method to use (GET, POST, etc.)")
	// rate := hfs.Int("rate", 1, "Number of requests per second")
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

	var pc cm.PlotC
	var attackAll []cm.PlotC

	switch strings.ToLower(*attacktype) {
	case "mailall":
		attackAll = append(attackAll, cm.PlotC{Results: st.BasicAttack(st.Options{})})
		attackAll = append(attackAll, cm.PlotC{Results: st.BurstAttack(st.Options{})})
		attackAll = append(attackAll, cm.PlotC{Results: st.RandomLoadAttack(st.Options{})})
		attackAll = append(attackAll, cm.PlotC{Results: st.RampUpAttack(st.Options{}, 1, 15)})
	case "basic":
		pc = cm.PlotC{Results: st.BasicAttack(st.Options{})}
	case "burst":
		pc = cm.PlotC{Results: st.BurstAttack(st.Options{})}
	case "random":
		pc = cm.PlotC{Results: st.RandomLoadAttack(st.Options{})}
	case "rampup":
		pc = cm.PlotC{Results: st.RampUpAttack(st.Options{}, 1, 15)}
	default:
		fmt.Println("Unknown attack type:", *attacktype)
		return
	}

	if *plot {
		pr.PlotResults(pc)
	}

	cfg := ml.Config{
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
		reportHTML := ml.BuildEmailReportHTML(attackAll, *url)
		cfg.SendMail(*emailTo, reportHTML)
	}
}
