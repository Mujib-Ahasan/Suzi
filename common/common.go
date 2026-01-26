package common

import (
	"time"
)

type Result struct {
	Status  string
	Elapsed time.Duration
	Error   error
}

type PResultIn struct {
	PRes InResult
	NRes []Result
	Err  error
}

type PlotC struct {
	Results PResultIn
	Attack  string
}

type InResult struct {
	URl                string
	Method             string
	Attack             string
	RequestsPerSecond  int
	Number_Of_Requests int
	Success_Count      int
	Error_Count        int
	ART                time.Duration
	P50                time.Duration
	P90                time.Duration
	P95                time.Duration
	P99                time.Duration
	Max                time.Duration
	Min                time.Duration
}
