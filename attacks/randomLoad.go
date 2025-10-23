package attacks

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	rs "github.com/Mujib-Ahasan/Suzi/common"
	sr "github.com/Mujib-Ahasan/Suzi/core"
)

func RandomLoadAttack(url string, numRequests int, method string, rate int, timeout int) rs.PResultIn {
	makeHandshake()

	var wg sync.WaitGroup
	resultsChan := make(chan rs.Result, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go makeRequest(url, method, &wg, resultsChan, timeout)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond / time.Duration(rate))
	}
	wg.Wait()
	close(resultsChan)
	results := make([]rs.Result, 0, numRequests)
	for result := range resultsChan {
		results = append(results, result)
	}
	sc := sr.ShowResults(results, numRequests, "random")
	fmt.Printf("%+v\n", sc)
	return rs.PResultIn{PRes: sc, NRes: results}
}
