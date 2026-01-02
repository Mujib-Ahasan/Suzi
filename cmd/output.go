package cmd

import (
	"fmt"

	"github.com/Mujib-Ahasan/Suzi/attacks"
	"github.com/Mujib-Ahasan/Suzi/common"
)

func printQuiet(s []common.PlotC) {
	for _, e := range s {
		fmt.Printf("attack=%s success=%d failures=%d avg=%dms p99=%dms\n", e.Attack, e.Results.PRes.Success_Count, e.Results.PRes.Error_Count, e.Results.PRes.ART.Milliseconds(), e.Results.PRes.P99.Milliseconds())
	}
}

func printVerbose(attacks []common.PlotC) {
	fmt.Println("🚀 Suzi Load Test Result")
	fmt.Println("────────────────────────────────")

	for _, e := range attacks {
		fmt.Printf("Attack Type    : %s\n", e.Attack)

		fmt.Printf("Total Requests : %d\n",
			e.Results.PRes.Success_Count+e.Results.PRes.Error_Count)

		fmt.Printf("Successful     : %d\n", e.Results.PRes.Success_Count)
		fmt.Printf("Failed         : %d\n", e.Results.PRes.Error_Count)

		fmt.Printf("Avg Latency    : %v\n", e.Results.PRes.ART)
		fmt.Printf("p50 Latency    : %v\n", e.Results.PRes.P50)
		fmt.Printf("p90 Latency    : %v\n", e.Results.PRes.P90)
		fmt.Printf("p95 Latency    : %v\n", e.Results.PRes.P95)
		fmt.Printf("P99 Latency    : %v\n", e.Results.PRes.P99)
		fmt.Println("────────────────────────────────")

		fmt.Println()
	}
}

func printVerboseHeader(opts attacks.Options) {
	fmt.Println("🚀 Suzi Load Test")
	fmt.Println("────────────────────────────")
	fmt.Printf("Target        %s\n", opts.URL)
	fmt.Printf("Method        %s\n", opts.Method)
	fmt.Printf("Rate          %d req/s\n", opts.Rate)
	fmt.Printf("Requests      %d\n", opts.Requests)
	fmt.Printf("Timeout       %s\n", opts.Timeout)

	if string(opts.Body) != "" {
		fmt.Printf("Body          %s\n", string(opts.Body))
	}

	if opts.EmailEnabled {
		fmt.Println("Email Alerts  ENABLED")
	} else {
		fmt.Println("Email Alerts  DISABLED")
	}

	fmt.Println("────────────────────────────")
	fmt.Println()
}

func printDefault(attacks []common.PlotC) {
	for _, e := range attacks {
		fmt.Printf("Attack: %s\n", e.Attack)

		total := e.Results.PRes.Success_Count + e.Results.PRes.Error_Count

		fmt.Printf("  Requests : %d\n", total)
		fmt.Printf("  Success  : %d\n", e.Results.PRes.Success_Count)
		fmt.Printf("  Failures : %d\n", e.Results.PRes.Error_Count)
		fmt.Printf("  Avg      : %v\n", e.Results.PRes.ART)
		fmt.Printf("  P99      : %v\n", e.Results.PRes.P99)

		fmt.Println()
	}
}
