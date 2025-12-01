package attacks

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	rs "github.com/Mujib-Ahasan/Suzi/common"
	sr "github.com/Mujib-Ahasan/Suzi/core"
)

func RandomLoadAttack(opts Options) rs.PResultIn {
	makeHandshake()

	var wg sync.WaitGroup
	resultsChan := make(chan rs.Result, opts.Requests)

	for i := 0; i < opts.Requests; i++ {
		wg.Add(1)
		go makeRequest(opts, &wg, resultsChan)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond / time.Duration(opts.Rate))
	}
	wg.Wait()
	close(resultsChan)
	results := make([]rs.Result, 0, opts.Requests)
	for result := range resultsChan {
		results = append(results, result)
	}
	sc := sr.ShowResults(results, opts.Requests, "random", opts.Method)
	fmt.Printf("%+v\n", sc)
	return rs.PResultIn{PRes: sc, NRes: results}
}
