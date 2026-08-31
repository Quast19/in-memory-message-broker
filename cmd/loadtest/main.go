package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	logFile, err := os.Create("loadtest.log")
	if err != nil {
		fmt.Println("failed to create log file:", err)
		return
	}
	defer logFile.Close()

	var wg sync.WaitGroup
	var mu sync.Mutex
	numRequests := 50
	successCount := 0
	failCount := 0
	var slowest time.Duration
	var fastest time.Duration

	fmt.Fprintf(logFile, "load test started: %s\n", time.Now().Format(time.RFC1123))
	fmt.Fprintf(logFile, "target: http://localhost:8080/publish\n")
	fmt.Fprintf(logFile, "total requests: %d\n", numRequests)
	fmt.Fprintln(logFile, "----------------------------")

	start := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()

			reqStart := time.Now()
			url := fmt.Sprintf("http://localhost:8080/publish?topic=orders&message=order-%d", n)
			resp, err := http.Get(url)
			elapsed := time.Since(reqStart)

			mu.Lock()
			defer mu.Unlock()

			if fastest == 0 || elapsed < fastest {
				fastest = elapsed
			}
			if elapsed > slowest {
				slowest = elapsed
			}

			if err != nil {
				failCount++
				fmt.Fprintf(logFile, "request %3d: FAILED (%v) in %v\n", n, err, elapsed)
				return
			}
			defer resp.Body.Close()

			successCount++
			fmt.Fprintf(logFile, "request %3d: status %s in %v\n", n, resp.Status, elapsed)
		}(i)
	}

	wg.Wait()
	totalElapsed := time.Since(start)

	fmt.Fprintln(logFile, "----------------------------")
	fmt.Fprintf(logFile, "finished:      %s\n", time.Now().Format(time.RFC1123))
	fmt.Fprintf(logFile, "succeeded:     %d\n", successCount)
	fmt.Fprintf(logFile, "failed:        %d\n", failCount)
	fmt.Fprintf(logFile, "total time:    %v\n", totalElapsed)
	fmt.Fprintf(logFile, "fastest req:   %v\n", fastest)
	fmt.Fprintf(logFile, "slowest req:   %v\n", slowest)
	fmt.Fprintf(logFile, "requests/sec:  %.2f\n", float64(numRequests)/totalElapsed.Seconds())
}
