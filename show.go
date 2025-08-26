package main

import (
	"fmt"
	"time"
)

type InResult struct {
	Number_Of_Requests int
	SuccessCount       int
	ErrorCount         int
	ART                time.Duration
}

func showResults(results []Result, numRequests int, attack string) InResult {
	var totalTime time.Duration
	var successCount, errorCount int

	for _, r := range results {
		if r.Error != nil {
			fmt.Println("Error:", r.Error)
			errorCount++
		} else {
			fmt.Println("Response Status:", r.Status)
			fmt.Println("Response Time:", r.Elapsed)
			totalTime += r.Elapsed
			successCount++
		}
	}

	avgTime := time.Duration(0)
	if successCount > 0 {
		avgTime = totalTime / time.Duration(successCount)
	}

	return InResult{
		Number_Of_Requests: numRequests,
		SuccessCount:       successCount,
		ErrorCount:         errorCount,
		ART:                avgTime,
	}
}
