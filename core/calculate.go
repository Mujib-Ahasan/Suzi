package core

import "time"

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
