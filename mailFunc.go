package main

import (
	"context"
	"fmt"
	"strings"
)

func (cfg *Config) sendMail(emailTo string, reportHTML string) {
	m := New(*cfg)
	// fmt.Println("debug-mailFunc_1")

	html := fmt.Sprintf(reportHTML)

	// Plain text fallback (optional)
	text := fmt.Sprintf(
		"Suzi Load Test Report:\n",
	)

	// fmt.Println("debug-mailFunc_2")

	// Attaching the generated results.html created by me.
	// atts := []Attachment{}
	// files, err := filepath.Glob("*.png")
	// if err != nil {
	// 	fmt.Println("glob error:", err)
	// 	return
	// }
	// fmt.Println("debug-mailFunc_3")

	// for _, f := range files {

	// 	atts = append(atts, Attachment{
	// 		Path: f,
	// 		Name: filepath.Base(f),
	// 	})
	// }
	// fmt.Println("debug-mailFunc_4")

	var recipients []string
	if emailTo != "" {
		for _, s := range strings.Split(emailTo, ",") {
			recipients = append(recipients, strings.TrimSpace(s))
		}
	}

	// fmt.Println("debug-mailFunc_5")

	ctx := context.Background()
	if err := m.Send(ctx, recipients, "Suzi Load Test Report", html, text); err != nil {
		fmt.Println("email send failed:", err)
	} else {
		fmt.Println("email sent successfully!!!")
	}
}
