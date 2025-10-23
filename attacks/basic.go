package attacks

import (
	"fmt"
	"sync"
	"time"

	// pr "github.com/Mujib-Ahasan/Suzi/cmd"
	rs "github.com/Mujib-Ahasan/Suzi/common"
	sr "github.com/Mujib-Ahasan/Suzi/core"
)

func BasicAttack(url string, numRequests int, rate int, method string, timeout int) rs.PResultIn {
	// numRequests=total number of request to be fired.
	// ratee=requests per second (RPS).
	makeHandshake()
	var wg sync.WaitGroup
	//multiple channels
	results_chan := make(chan rs.Result, numRequests)
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

	results := make([]rs.Result, 0, numRequests)
	for result := range results_chan {
		results = append(results, result)
	}
	sc := sr.ShowResults(results, numRequests, "basic")
	fmt.Printf("%+v\n", sc)

	return rs.PResultIn{PRes: sc, NRes: results}
}
