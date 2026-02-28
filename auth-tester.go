package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	url         = "http://localhost:7070/auth/login"
	concurrency = 50   // number of goroutines
	requests    = 2000 // total requests
)

func main() {
	var wg sync.WaitGroup
	start := time.Now()

	jobs := make(chan struct{}, requests)

	for i := 0; i < requests; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	payload := []byte(`{
		"email": "test@example.com",
		"password": "password123"
	}`)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for range jobs {
				req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err == nil {
					resp.Body.Close()
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Load test finished in:", time.Since(start))
}
