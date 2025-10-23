package attacks

import (
	"fmt"
	"sync"

	// pr "github.com/Mujib-Ahasan/Suzi/cmd"
	rs "github.com/Mujib-Ahasan/Suzi/common"
	sr "github.com/Mujib-Ahasan/Suzi/core"
)

func BurstAttack(url string, numRequests int, method string, timeout int) rs.PResultIn {
	makeHandshake()

	var wg sync.WaitGroup
	resultsChan := make(chan rs.Result, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go makeRequest(url, method, &wg, resultsChan, timeout)
	}

	wg.Wait()
	close(resultsChan)

	results := make([]rs.Result, 0, numRequests)
	for r := range resultsChan {
		results = append(results, r)
	}
	sc := sr.ShowResults(results, numRequests, "burst")
	fmt.Printf("%+v\n", sc)
	return rs.PResultIn{PRes: sc, NRes: results}
}
