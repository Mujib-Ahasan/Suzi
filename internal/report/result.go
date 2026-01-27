package report

import "time"

type LoadTestResultAll struct {
	Result []LoadTestResult `json:"result"`
	Email  *EmailConfig     `json:"email,omitempty"`
}

type LoadTestResult struct {
	Version    string       `json:"version"`
	Test       TestInfo     `json:"test"`
	Summary    Summary      `json:"summary"`
	Latency    LatencyStats `json:"latency"`
	Throughput Throughput   `json:"throughput"`
	Errors     []ErrorStat  `json:"errors,omitempty"`
	Timestamp  Timestamp    `json:"timestamp"`
}

type EmailConfig struct {
	From        string        `json:"from"`
	DialTimeout time.Duration `json:"dialTimeout,omitempty"`
	SendTimeout time.Duration `json:"sendTimeout"`
	Retries     int           `json:"retries"`
	To          string        `json:"to"`
}

type TestInfo struct {
	Name          string `json:"name"`
	AttackType    string `json:"attackType"`
	TargetURL     string `json:"targetUrl"`
	Method        string `json:"method"`
	TotalRequests int    `json:"totalRequests"`
}

type Summary struct {
	RequestsTotal   int     `json:"requestsTotal"`
	RequestsSuccess int     `json:"requestsSuccess"`
	RequestsFailed  int     `json:"requestsFailed"`
	SuccessRate     float64 `json:"successRate"`
}

type LatencyStats struct {
	MinMs time.Duration `json:"minMs"`
	MaxMs time.Duration `json:"maxMs"`
	AvgMs time.Duration `json:"avgMs"`
	P50Ms time.Duration `json:"p50Ms"`
	P90Ms time.Duration `json:"p90Ms"`
	P95Ms time.Duration `json:"p95Ms"`
	P99Ms time.Duration `json:"p99Ms"`
}

type Throughput struct {
	RequestsPerSecond int     `json:"requestsPerSecond"`
	BytesPerSecond    float64 `json:"bytesPerSecond,omitempty"`
}

type ErrorStat struct {
	Type  string `json:"type"` // timeout, connection, http_error
	Code  int    `json:"code,omitempty"`
	Count int    `json:"count"`
}

type Timestamp struct {
	StartedAt  string `json:"startedAt"`
	FinishedAt string `json:"finishedAt"`
	DurationMs int64  `json:"durationMs"`
}
