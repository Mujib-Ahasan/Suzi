package main

import (
	"bytes"
	"fmt"
	"time"
)

// generate a styled HTML report for email
func BuildEmailReportHTML(plots []plotC, url string) string {
	var buf bytes.Buffer

	intro := `
    <div style="padding:20px; text-align:center; font-family:Arial, sans-serif;">
        <h1 style="color:#2c3e50; margin-bottom:10px;">🚀 Suzi Load Testing Report</h1>
        <p style="color:#555; font-size:15px; margin-top:0;">
            Here's a snapshot of how your system performed under stress.  
            Each block below summarizes one attack with latency percentiles, success & error counts.
        </p>
    </div>
    `
	buf.WriteString(intro)

	for _, plot := range plots {
		// Current time in IST
		ist := time.Now().In(time.FixedZone("IST", 5*3600+1800))
		timestamp := ist.Format("02 Jan 2006 15:04:05 MST")

		tableHTML := fmt.Sprintf(`
    <div style="display:flex; justify-content:center; margin-top:25px;">
        <div style="padding:20px; border:1px solid #cce5cc; border-radius:10px; width:650px; font-family:Arial, sans-serif; box-shadow:0 4px 10px rgba(0,0,0,0.1); background:#f6fff6;">
            
            <!-- Header -->
            <h2 style="text-align:center; color:#2c662d; margin-bottom:5px;">⚡ Suzi Load Test Report</h2>
            <p style="text-align:center; color:#444; margin-top:0;">Attack: <b>%s</b> | Generated on: %s</p>
            
            <!-- Summary -->
            <div style="display:flex; justify-content:space-around; margin:15px 0;">
                <div style="padding:10px; background:#eafaf1; border-radius:6px; flex:1; margin:0 5px; text-align:center; border:1px solid #b6e3c1;">
                    ✅ <b>Success</b><br/>%d
                </div>
                <div style="padding:10px; background:#fdecea; border-radius:6px; flex:1; margin:0 5px; text-align:center; border:1px solid #f5c6c6;">
                    ❌ <b>Errors</b><br/>%d
                </div>
            </div>

            <!-- Percentiles Table -->
            <div style="display:flex; justify-content:center; margin-top:20px;">
                <table cellpadding="10" style="border-collapse:collapse; width:80%%; text-align:center; border:1px solid #ccc; background:#ffffff;">
                    <tr style="background-color:#e6f2e6; font-weight:bold; color:#2c662d;">
                        <th>Metric</th>
                        <th>p50</th>
                        <th>p90</th>
                        <th>p95</th>
                        <th>p99</th>
                    </tr>
                    <tr>
                        <td style="font-weight:bold;">Latency (ms)</td>
                        <td>%.2f</td>
                        <td>%.2f</td>
                        <td>%.2f</td>
                        <td>%.2f</td>
                    </tr>
                </table>
            </div>
        </div>
    </div>
`,
			plot.attack, timestamp,
			plot.results.PRes.Success_Count, plot.results.PRes.Error_Count,
			float64(plot.results.PRes.P50)/float64(time.Millisecond),
			float64(plot.results.PRes.P90)/float64(time.Millisecond),
			float64(plot.results.PRes.P95)/float64(time.Millisecond),
			float64(plot.results.PRes.P99)/float64(time.Millisecond),
		)

		buf.WriteString(tableHTML)
	}

	// Add a footer with button
	footer := fmt.Sprintf(`
    <div style="text-align:center; margin-top:30px;">
        <a href="%s" target="_blank" 
           style="background:#28a745; color:white; padding:12px 25px; text-decoration:none; border-radius:5px; font-size:16px; font-weight:bold;">
           🔗 Visit Tested Website
        </a>
    </div>
`, url)
	buf.WriteString(footer)

	outro := `
    <div style="padding:20px; text-align:center; margin-top:40px; font-family:Arial, sans-serif; font-size:14px; color:#555;">
        <p><b>About Suzi</b>: Suzi is an open-source load testing tool designed to help developers
        benchmark their systems in real scenarios. While Suzi is built with care, results may not be
        perfect and occasional errors can occur, please validate before making critical decisions.</p>
        <p>Contribute, raise issues, or star ⭐ the project on GitHub:<br/>
        <a href="https://github.com/Mujib-Ahasan/Suzi" target="_blank" style="color:#28a745; font-weight:bold;">github.com/Mujib-Ahasan/Suzi</a></p>
    </div>
    `
	buf.WriteString(outro)

	return buf.String()
}
