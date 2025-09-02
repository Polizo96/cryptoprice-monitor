package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// A request client simulator to check how fast server is responding

const (
	apiURL          = "http://localhost:8080/prices?coin=BTC&exchange=coinbase"
	concurrentUsers = 10
	requestsPerUser = 100
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{Timeout: 5 * time.Second}

	for i := 0; i < requestsPerUser; i++ {
		start := time.Now()
		resp, err := client.Get(apiURL)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("User %d Request %d error: %v\n", id, i+1, err)
		}

		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}

		fmt.Printf("User %d Request %d status: %s, took: %v\n", id, i+1, resp.Status, duration)
	}
}

func main() {
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go worker(i+1, &wg)
	}
	wg.Wait()
	fmt.Printf("All requests completed in %v\n", time.Since(start))
}
