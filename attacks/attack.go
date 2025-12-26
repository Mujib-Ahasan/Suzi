package attacks

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	rs "github.com/Mujib-Ahasan/Suzi/common"
)

// setting it global so to avoid handshake each time(no TCP/TLS each time)
var client_body *http.Client

func makeHandshake() {
	fmt.Println("New connection made!!!")
	client_body = &http.Client{}
}

// = &http.Client{}

// this function sends the HTTP request and send response woth some data through chanel.
func makeRequest(opts Options, wg *sync.WaitGroup, resp_results chan<- rs.Result) {
	defer wg.Done()
	start := time.Now()

	var body io.Reader

	if (opts.Method == "POST" || opts.Method == "PUT") && len(opts.Body) > 0 {
		body = bytes.NewReader(opts.Body)
	}

	// making a request....
	req, err := http.NewRequest(opts.Method, opts.URL, body)
	if err != nil {
		resp_results <- rs.Result{Error: err}
		return
	}
	if (opts.Method == "POST" || opts.Method == "PUT") && len(opts.Body) > 0 {
		req.Header.Set("Content-Type", opts.ContentType)
		fmt.Printf("Content typeeeeeeeeeeeeeee attack.go: %s \n", opts.ContentType)
	}

	//setting timeout.
	client_body.Timeout = opts.Timeout * time.Second
	// sending that request....
	resp, err := client_body.Do(req)
	// bodyByte, _ := io.ReadAll(resp.Body)
	// fmt.Printf("responseeeeeeeeee: %v \n", string(bodyByte))
	// defer resp.Body.Close()
	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("Error_2:  %v\n", err)
		resp_results <- rs.Result{Error: err}
		return
	}
	defer resp.Body.Close()
	// sending the respose result through channel....
	resp_results <- rs.Result{Status: resp.Status, Elapsed: elapsed}
}
