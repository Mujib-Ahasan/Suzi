package main

import (
	"fmt"
	"os"

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
	xData := make([]int, len(plot.results))
	for i := range xData {
		xData[i] = i + 1
	}
	//yData holds the number of miliseconds it took!
	yData := make([]opts.LineData, len(plot.results))
	for i, result := range plot.results {
		// a single data point for a line chart.
		yData[i] = opts.LineData{Value: result.Elapsed.Milliseconds()}
	}

	// line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{ //will add more options later.
	// 	Title:    "LTRs", // load test results
	// 	Subtitle: "RTPS", //Response time per second
	// }))

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "📊 Load Test Results",
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

	line.SetXAxis(xData).AddSeries("Time for Response", yData)

	//save the chart as an html file.
	file_name := "result-" + plot.attack + ".html"
	f, err := os.Create(file_name)
	if err != nil {
		fmt.Println("Could not create results.html:", err)
		return
	}
	defer f.Close()

	if err := line.Render(f); err != nil {
		fmt.Println("Could not render chart:", err)
	}

	fmt.Println("Results plotted to " + file_name)
}
