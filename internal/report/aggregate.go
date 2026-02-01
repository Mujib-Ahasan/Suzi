package report

import (
	"fmt"
	"time"

	"github.com/Mujib-Ahasan/Suzi/common"
	"github.com/Mujib-Ahasan/Suzi/mail"
)

func FromAttackList(list []common.PlotC) LoadTestResultAll {
	var (
		totalRequests   int
		successRequests int
		failedRequests  int
		errorMap        = make(map[string]*ErrorStat)
	)
	result := []LoadTestResult{}
	for _, r := range list {

		start := time.Now()

		totalRequests = r.Results.PRes.Number_Of_Requests
		successRequests = r.Results.PRes.Success_Count
		failedRequests = r.Results.PRes.Error_Count

		for _, e := range r.Results.NRes {
			if e.Error != nil {
				key := e.Error.Error()

				if _, ok := errorMap[key]; !ok {
					errorMap[key] = &ErrorStat{
						Type:  fmt.Sprintf(key),
						Count: 0,
					}
				}
				errorMap[key].Count++
			}
		}

		latency := LatencyStats{
			MinMs: time.Duration(r.Results.PRes.Min.Milliseconds()),
			MaxMs: time.Duration(r.Results.PRes.Max.Milliseconds()),
			AvgMs: time.Duration(r.Results.PRes.ART.Milliseconds()),
			P50Ms: time.Duration(r.Results.PRes.P50.Milliseconds()),
			P90Ms: time.Duration(r.Results.PRes.P90.Milliseconds()),
			P95Ms: time.Duration(r.Results.PRes.P95.Milliseconds()),
			P99Ms: time.Duration(r.Results.PRes.P99.Milliseconds()),
		}
		errors := make([]ErrorStat, 0, len(errorMap))
		for _, e := range errorMap {
			errors = append(errors, *e)
		}

		finished := time.Now()

		testResult := LoadTestResult{
			Test: TestInfo{
				Name:          "Suzi Load testing report",
				AttackType:    r.Attack,
				TargetURL:     r.Results.PRes.URl,
				Method:        r.Results.PRes.Method,
				TotalRequests: r.Results.PRes.Number_Of_Requests,
			},
			Version: "v1",
			Summary: Summary{
				RequestsTotal:   totalRequests,
				RequestsSuccess: successRequests,
				RequestsFailed:  failedRequests,
				SuccessRate:     percent(successRequests, totalRequests),
			},
			Latency: latency,
			Throughput: Throughput{
				RequestsPerSecond: r.Results.PRes.RequestsPerSecond,
				ResBytesPerSecond: float64(r.Results.PRes.TotalRespBytes) / r.Results.PRes.TotalElapsedTime.Seconds(),
			},
			Errors: errors,
			Timestamp: Timestamp{
				StartedAt:  start.Format(time.RFC3339),
				FinishedAt: finished.Format(time.RFC3339),
			},
		}

		result = append(result, testResult)

	}

	resultAll := LoadTestResultAll{
		Result: result,
	}

	return resultAll
}

func percent(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return (float64(a) / float64(b)) * 100
}

func FromAttackListEmail(list []common.PlotC, email mail.Config) LoadTestResultAll {
	var (
		totalRequests   int
		successRequests int
		failedRequests  int
		errorMap        = make(map[string]*ErrorStat)
	)
	result := []LoadTestResult{}
	mail := EmailConfig{}
	for _, r := range list {

		start := time.Now()

		totalRequests = r.Results.PRes.Number_Of_Requests
		successRequests = r.Results.PRes.Success_Count
		failedRequests = r.Results.PRes.Error_Count

		for _, e := range r.Results.NRes {
			if e.Error != nil {
				key := e.Error.Error()

				if _, ok := errorMap[key]; !ok {
					errorMap[key] = &ErrorStat{
						Type:  fmt.Sprintf(key),
						Count: 0,
					}
				}
				errorMap[key].Count++
			}
		}

		latency := LatencyStats{
			MinMs: time.Duration(r.Results.PRes.Min.Milliseconds()),
			MaxMs: time.Duration(r.Results.PRes.Max.Milliseconds()),
			AvgMs: time.Duration(r.Results.PRes.ART.Milliseconds()),
			P50Ms: time.Duration(r.Results.PRes.P50.Milliseconds()),
			P90Ms: time.Duration(r.Results.PRes.P90.Milliseconds()),
			P95Ms: time.Duration(r.Results.PRes.P95.Milliseconds()),
			P99Ms: time.Duration(r.Results.PRes.P99.Milliseconds()),
		}

		mail = EmailConfig{
			From:        email.FromEmail,
			To:          email.ToEmail,
			DialTimeout: email.DialTimeout,
			SendTimeout: email.SendTimeout,
			Retries:     email.Retries,
		}
		errors := make([]ErrorStat, 0, len(errorMap))
		for _, e := range errorMap {
			errors = append(errors, *e)
		}

		finished := time.Now()

		testResult := LoadTestResult{
			Test: TestInfo{
				Name:          "Suzi Load testing report",
				AttackType:    r.Attack,
				TargetURL:     r.Results.PRes.URl,
				Method:        r.Results.PRes.Method,
				TotalRequests: r.Results.PRes.Number_Of_Requests,
			},
			Version: "v1",
			Summary: Summary{
				RequestsTotal:   totalRequests,
				RequestsSuccess: successRequests,
				RequestsFailed:  failedRequests,
				SuccessRate:     percent(successRequests, totalRequests),
			},
			Latency: latency,
			Throughput: Throughput{
				RequestsPerSecond: r.Results.PRes.RequestsPerSecond,
				ResBytesPerSecond: float64(r.Results.PRes.TotalRespBytes) / r.Results.PRes.TotalElapsedTime.Seconds(),
			},
			Errors: errors,
			Timestamp: Timestamp{
				StartedAt:  start.Format(time.RFC3339),
				FinishedAt: finished.Format(time.RFC3339),
			},
		}

		result = append(result, testResult)

	}
	resultAll := LoadTestResultAll{
		Result: result,
		Email:  &mail,
	}

	return resultAll

}
