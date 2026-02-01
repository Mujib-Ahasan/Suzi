package core

import (
	"log/slog"
	"sort"
	"time"

	rs "github.com/Mujib-Ahasan/Suzi/common"
)

func ShowResults(results []rs.Result, numRequests int, attack string, method string) rs.InResult {
	var totalTime time.Duration
	var avg, p50, p90, p95, p99, max, min time.Duration

	var successCount, errorCount int
	var totalResBytes int64
	latencies := make([]time.Duration, 0, successCount)

	for _, r := range results {
		if r.Error != nil {
			errorCount++
		} else {
			slog.Debug("Response Status:", r.Status)
			slog.Debug("Response Time:", r.Elapsed)
			totalTime += r.Elapsed
			successCount++
			latencies = append(latencies, r.Elapsed)
		}

		totalResBytes += r.ResponseBytes
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
			min = latencies[0]
		}
	}

	return rs.InResult{
		Method:             method,
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
		Min:                min,
		TotalRespBytes:     totalResBytes,
		TotalElapsedTime:   totalTime,
	}
}
