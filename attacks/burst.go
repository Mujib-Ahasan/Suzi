package attacks

import (
	"fmt"
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
		go makeRequest(opts.URL, opts.Method, &wg, resultsChan, opts.Timeout)
	}

	wg.Wait()
	close(resultsChan)

	results := make([]rs.Result, 0, opts.Requests)
	for r := range resultsChan {
		results = append(results, r)
	}
	sc := sr.ShowResults(results, opts.Requests, "burst")
	fmt.Printf("%+v\n", sc)
	return rs.PResultIn{PRes: sc, NRes: results}
}
