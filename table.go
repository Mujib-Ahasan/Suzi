package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// little bit context:
// github.com/go-echarts/go-echarts/v2/charts: contains the chart types canbe created. for now linechart only.
// github.com/go-echarts/go-echarts/v2/opts: contains the configration structs that describe chart data.

// generate a line table/chart of response time.
func plotResults(plot plotC) {
	line := charts.NewLine()
	//xData holds the request numbering
	xData := make([]int, len(plot.results.NRes))
	for i := range xData {
		xData[i] = i + 1
	}
	//yData holds the number of miliseconds it took!
	yData := make([]opts.LineData, len(plot.results.NRes))
	for i, result := range plot.results.NRes {
		// a single data point for a line chart.
		yData[i] = opts.LineData{Value: result.Elapsed.Milliseconds()}
	}

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Load Test Results",
			Subtitle: "Response Time per Second (RTPS)",
			Left:     "center", // center align title
		}),
		// Add tooltip for better readability
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      opts.Bool(true),
			Trigger:   "axis",
			Formatter: "{b} req<br/>⏱ {c} ms",
		}),
		// Add X axis name
		charts.WithXAxisOpts(opts.XAxis{
			Name:    "Request #",
			NameGap: 25,
		}),
		// Add Y axis name
		charts.WithYAxisOpts(opts.YAxis{
			Name:    "Response Time (ms)",
			NameGap: 30,
			AxisLabel: &opts.AxisLabel{
				Show:      opts.Bool(true),
				Formatter: "{value} ms",
			},
		}),
		// Add legend
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(true),
			Left: "right",
		}),
	)

	tableHTML := fmt.Sprintf(`
    <div style="display:flex; justify-content:center; margin-top:20px;">
        <div style="padding:15px; border:1px solid #ccc; border-radius:8px; width:600px; font-family:Arial; box-shadow:0 4px 8px rgba(0,0,0,0.1);">
            <h3 style="text-align:center; margin-bottom:10px;">⚡ Latency Percentiles</h3>
            <table border="1" cellpadding="8" style="border-collapse:collapse; width:100%%; text-align:center;">
                <tr style="background-color:#f2f2f2;">
                    <th>Metric</th>
                    <th>p50</th>
                    <th>p90</th>
                    <th>p95</th>
                    <th>p99</th>
                </tr>
                <tr>
                    <td>Latency (ms)</td>
                    <td>%.2f</td>
                    <td>%.2f</td>
                    <td>%.2f</td>
                    <td>%.2f</td>
                </tr>
            </table>
        </div>
    </div>
`,
		float64(plot.results.PRes.P50)/float64(time.Millisecond),
		float64(plot.results.PRes.P90)/float64(time.Millisecond),
		float64(plot.results.PRes.P95)/float64(time.Millisecond),
		float64(plot.results.PRes.P99)/float64(time.Millisecond),
	)

	line.SetXAxis(xData).AddSeries("Time for Response", yData)

	//saving the chart as html file.
	file_name := "result-" + plot.attack + ".html"
	f, err := os.Create(file_name)
	if err != nil {
		fmt.Println("Could not create results.html:", err)
		return
	}
	defer f.Close()

	if err := line.Render(f); err != nil {
		fmt.Println("Could not render chart:", err)
		return
	}
	io.WriteString(f, tableHTML)
	fmt.Println("Results plotted to " + file_name)
}
