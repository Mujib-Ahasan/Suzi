package attacks

import (
	"fmt"
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
func makeRequest(url string, method string, wg *sync.WaitGroup, resp_results chan<- rs.Result, timeout int) {
	defer wg.Done()
	start := time.Now()

	// makign a request....
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Printf("Error_1:  %v\n", err)
		resp_results <- rs.Result{Error: err}
		return
	}
	//setting timeout.
	client_body.Timeout = time.Duration(timeout) * time.Second
	// sending that request....
	resp, err := client_body.Do(req)
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
