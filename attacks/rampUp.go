package attacks

import (
	"fmt"
	"sync"
	"time"

	// pr "github.com/Mujib-Ahasan/Suzi/cmd"
	rs "github.com/Mujib-Ahasan/Suzi/common"
	sr "github.com/Mujib-Ahasan/Suzi/core"
)

func RampUpAttack(opts Options, startRate int, peakRate int) rs.PResultIn {
	makeHandshake()

	var wg sync.WaitGroup
	resultsChan := make(chan rs.Result, opts.Requests)

	// Linearly ramp rate from startRate → peakRate
	rateStep := float64(peakRate-startRate) / float64(opts.Requests)

	for i := 0; i < opts.Requests; i++ {
		wg.Add(1)

		// Calculate current rate (linearly increasing)
		currentRate := float64(startRate) + (rateStep * float64(i))
		// Convert rate → sleep duration (inverse relation)
		sleepDuration := time.Second / time.Duration(currentRate)
		go makeRequest(opts, &wg, resultsChan)
		// Control pacing here
		time.Sleep(sleepDuration)
	}

	wg.Wait()
	close(resultsChan)
	results := make([]rs.Result, 0, opts.Requests)
	for result := range resultsChan {
		results = append(results, result)
	}

	sc := sr.ShowResults(results, opts.Requests, "rampup", opts.Method)
	fmt.Printf("%+v\n", sc)
	return rs.PResultIn{PRes: sc, NRes: results}
}
