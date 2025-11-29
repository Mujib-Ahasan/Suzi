package attacks

import (
	"fmt"
	"sync"
	"time"

	// pr "github.com/Mujib-Ahasan/Suzi/cmd"
	rs "github.com/Mujib-Ahasan/Suzi/common"
	sr "github.com/Mujib-Ahasan/Suzi/core"
)

func BasicAttack(opts Options) rs.PResultIn {
	// opts.numRequests=total number of request to be fired.
	// ratee=requests per second (RPS).
	makeHandshake()
	var wg sync.WaitGroup
	//multiple channels
	results_chan := make(chan rs.Result, opts.Requests)
	//emits signals in  every 1/rate second.
	// time.NewTicker(d) returns a Ticker that repeatedly sends the current time on its channel ticker.C every d duration.
	ticker := time.NewTicker(time.Second / time.Duration(opts.Rate))
	defer ticker.Stop()

	for i := 0; i < opts.Requests; i++ {
		wg.Add(1)
		// blocks until the next tick, spacing requests evenly.
		<-ticker.C
		go makeRequest(opts.URL, opts.Method, &wg, results_chan, opts.Timeout)
	}
	wg.Wait()
	close(results_chan)

	results := make([]rs.Result, 0, opts.Requests)
	for result := range results_chan {
		results = append(results, result)
	}
	sc := sr.ShowResults(results, opts.Requests, "basic")
	fmt.Printf("%+v\n", sc)

	return rs.PResultIn{PRes: sc, NRes: results}
}
