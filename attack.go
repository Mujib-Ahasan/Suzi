package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	Status  string
	Elapsed time.Duration
	Error   error
}

// setting it global so to avoid handshake each time(no TCP/TLS each time)
var client_body = &http.Client{}

// this function sends the HTTP request and send response woth some data through chanel.
func makeRequest(url string, method string, wg *sync.WaitGroup, resp_results chan<- Result, timeout int) {
	defer wg.Done()
	start := time.Now()

	// makign a request....
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Printf("Error_1:  %v\n", err)
		resp_results <- Result{Error: err}
		return
	}
	//setting timeout.
	client_body.Timeout = time.Duration(timeout) * time.Second
	// sending that request....
	resp, err := client_body.Do(req)
	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("Error_2:  %v\n", err)
		resp_results <- Result{Error: err}
		return
	}
	defer resp.Body.Close()
	// sending the respose result through channel....
	resp_results <- Result{Status: resp.Status, Elapsed: elapsed}
}

func basicAttack(url string, numRequests int, rate int, method string, timeout int) []Result {
	// numRequests=total number of request to be fired.
	// ratee=requests per second (RPS).
	var wg sync.WaitGroup
	//multiple channels
	results_chan := make(chan Result, numRequests)
	//emits signals in  every 1/rate second.
	// time.NewTicker(d) returns a Ticker that repeatedly sends the current time on its channel ticker.C every d duration.
	ticker := time.NewTicker(time.Second / time.Duration(rate))
	defer ticker.Stop()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		// blocks until the next tick, spacing requests evenly.
		<-ticker.C
		go makeRequest(url, method, &wg, results_chan, timeout)
	}
	wg.Wait()
	close(results_chan)

	results := make([]Result, 0, numRequests)
	for result := range results_chan {
		results = append(results, result)
	}
	sc := showResults(results, numRequests, "basic")
	fmt.Printf("%+v\n", sc)

	return results
}

func burstAttack(url string, numRequests int, method string, timeout int) []Result {
	var wg sync.WaitGroup
	resultsChan := make(chan Result, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go makeRequest(url, method, &wg, resultsChan, timeout)
	}

	wg.Wait()
	close(resultsChan)

	results := make([]Result, 0, numRequests)
	for r := range resultsChan {
		results = append(results, r)
	}
	sc := showResults(results, numRequests, "burst")
	fmt.Printf("%+v\n", sc)
	return results
}

func randomLoadAttack(url string, numRequests int, method string, rate int, timeout int) []Result {
	var wg sync.WaitGroup
	resultsChan := make(chan Result, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go makeRequest(url, method, &wg, resultsChan, timeout)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond / time.Duration(rate))
	}
	wg.Wait()
	close(resultsChan)
	results := make([]Result, 0, numRequests)
	for result := range resultsChan {
		results = append(results, result)
	}
	sc := showResults(results, numRequests, "random")
	fmt.Printf("%+v\n", sc)
	return results
}
