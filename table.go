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
func plotResults(results []Result) {
	line := charts.NewLine()
	//xData holds the request numbering
	xData := make([]int, len(results))
	for i := range xData {
		xData[i] = i + 1
	}
	//yData holds the number of miliseconds it took!
	yData := make([]opts.LineData, len(results))
	for i, result := range results {
		// a single data point for a line chart.
		yData[i] = opts.LineData{Value: result.Elapsed.Milliseconds()}
	}

	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{ //will add more options later.
		Title:    "LTRs", // load test results
		Subtitle: "RTPS", //Response time per second
	}))
	line.SetXAxis(xData).AddSeries("Time for Response", yData)

	//save the chart as an html file.
	f, err := os.Create("results.html")
	if err != nil {
		fmt.Println("Could not create results.html:", err)
		return
	}
	defer f.Close()

	if err := line.Render(f); err != nil {
		fmt.Println("Could not render chart:", err)
	}

	fmt.Println("Results plotted to results.html")
}
