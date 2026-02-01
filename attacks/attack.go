package attacks

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	rs "github.com/Mujib-Ahasan/Suzi/common"
)

var (
	client_body *http.Client
	once        sync.Once
)

func makeHandshake() {
	once.Do(func() {
		slog.Debug("HTTP client initialized")
		client_body = &http.Client{}
	})
}

func makeRequest(opts Options, wg *sync.WaitGroup, resp_results chan<- rs.Result) {
	defer wg.Done()
	start := time.Now()

	body := opts.ConvertBody()

	req, err := http.NewRequest(opts.Method, opts.URL, body)
	if err != nil {
		resp_results <- rs.Result{Error: err}
		return
	}

	if (opts.Method == "POST" || opts.Method == "PUT" || opts.Method == "PATCH") && len(opts.Body) > 0 {
		req.Header.Set("Content-Type", opts.ContentType)
	}

	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	// per-request timeout (SAFE)
	ctx, cancel := context.WithTimeout(req.Context(), opts.Timeout*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := client_body.Do(req)
	elapsed := time.Since(start)
	if err != nil {
		resp_results <- rs.Result{Error: err}
		return
	}
	defer resp.Body.Close()

	cr := &countingReadCloser{rc: resp.Body}
	_, _ = io.Copy(io.Discard, cr)

	resp_results <- rs.Result{
		Status:        resp.Status,
		Elapsed:       elapsed,
		ResponseBytes: cr.bytes,
	}
}

func (opts Options) ConvertBody() io.Reader {
	var body io.Reader
	if (opts.Method == "POST" || opts.Method == "PUT" || opts.Method == "PATCH") && len(opts.Body) > 0 {
		body = bytes.NewReader(opts.Body)
	}

	return body
}
