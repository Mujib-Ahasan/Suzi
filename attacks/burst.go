package attacks

import (
	"sync"

	// pr "github.com/Mujib-Ahasan/Suzi/cmd"
	rs "github.com/Mujib-Ahasan/Suzi/common"
	sr "github.com/Mujib-Ahasan/Suzi/core"
)

func BurstAttack(opts Options) rs.PResultIn {
	makeHandshake()

	var wg sync.WaitGroup
	resultsChan := make(chan rs.Result, opts.Requests)

	for i := 0; i < opts.Requests; i++ {
		wg.Add(1)
		go makeRequest(opts, &wg, resultsChan)
	}

	wg.Wait()
	close(resultsChan)

	results := make([]rs.Result, 0, opts.Requests)
	for r := range resultsChan {
		results = append(results, r)
	}
	sc := sr.ShowResults(results, opts.Requests, "burst", opts.Method)
	return rs.PResultIn{PRes: sc, NRes: results}
}
