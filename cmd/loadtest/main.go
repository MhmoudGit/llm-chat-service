package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	url := flag.String("url", "http://localhost:8080/chat", "URL to chat endpoint")
	concurrency := flag.Int("c", 50, "Number of concurrent requests")
	flag.Parse()

	fmt.Printf("Starting load test against %s with %d concurrent requests...\n", *url, *concurrency)

	var wg sync.WaitGroup
	var successCount int64
	var failCount int64

	start := time.Now()

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := sendRequest(*url, id); err != nil {
				fmt.Printf("[Req %d] Failed: %v\n", id, err)
				atomic.AddInt64(&failCount, 1)
			} else {
				// fmt.Printf("[Req %d] Success\n", id)
				atomic.AddInt64(&successCount, 1)
			}
		}(i)
		time.Sleep(500 * time.Millisecond)
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Println("\n--- Load Test Results ---")
	fmt.Printf("Total Requests: %d\n", *concurrency)
	fmt.Printf("Success: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failCount)
	fmt.Printf("Total Duration: %v\n", duration)

	if failCount == 0 {
		fmt.Println("Status: PASSED")
	} else {
		fmt.Println("Status: COMPLETED (with errors)")
	}
}

func sendRequest(url string, id int) error {
	payload := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "user", "content": "Hello load test"},
		},
		"stream": true,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	// Consume stream
	_, err = io.Copy(io.Discard, resp.Body)
	return err
}
