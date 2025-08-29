package main

import (
	"fmt"
	"sort"
	"time"
)

type inResult struct {
	Attack             string
	Number_Of_Requests int
	Success_Count      int
	Error_Count        int
	ART                time.Duration
	P50                time.Duration
	P90                time.Duration
	P95                time.Duration
	P99                time.Duration
	Max                time.Duration
}

func percentileCalculation(latencies []time.Duration, p float64) time.Duration {
	if len(latencies) == 0 {
		return 0
	}
	index := int((p / 100.0) * float64(len(latencies)))
	if index >= len(latencies) {
		index = len(latencies) - 1
	}
	return latencies[index]
}

func showResults(results []Result, numRequests int, attack string) inResult {
	var totalTime time.Duration
	var avg, p50, p90, p95, p99, max time.Duration

	var successCount, errorCount int
	latencies := make([]time.Duration, 0, successCount)
	for _, r := range results {
		if r.Error != nil {
			fmt.Println("Error:", r.Error)
			errorCount++
		} else {
			fmt.Println("Response Status:", r.Status)
			fmt.Println("Response Time:", r.Elapsed)
			totalTime += r.Elapsed
			successCount++
			latencies = append(latencies, r.Elapsed)
		}
	}

	if successCount > 0 {
		sort.Slice(latencies, func(i, j int) bool {
			return latencies[i] < latencies[j]
		})

		if successCount > 0 {
			avg = totalTime / time.Duration(successCount)
			p50 = percentileCalculation(latencies, 50)
			p90 = percentileCalculation(latencies, 90)
			p95 = percentileCalculation(latencies, 95)
			p99 = percentileCalculation(latencies, 99)
			max = latencies[len(latencies)-1]
		}
	}

	return inResult{
		Attack:             attack,
		Number_Of_Requests: numRequests,
		Success_Count:      successCount,
		Error_Count:        errorCount,
		ART:                avg,
		P50:                p50,
		P90:                p90,
		P95:                p95,
		P99:                p99,
		Max:                max,
	}
}
