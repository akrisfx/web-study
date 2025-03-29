package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type LoadTester struct {
	URL            string
	NumWorkers     int
	TestDuration   time.Duration
	RequestCounter int64
	Client         *http.Client
}

func NewLoadTester(url string, workers int, duration time.Duration) *LoadTester {
	return &LoadTester{
		URL:          url,
		NumWorkers:   workers,
		TestDuration: duration,
		Client: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 100,
				MaxConnsPerHost:     100,
			},
		},
	}
}

func (lt *LoadTester) Worker(wg *sync.WaitGroup, stop <-chan struct{}) {
	defer wg.Done()

	for {
		select {
		case <-stop:
			return
		default:
			resp, err := lt.Client.Get(lt.URL)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					atomic.AddInt64(&lt.RequestCounter, 1)
				}
			}
		}
	}
}

func (lt *LoadTester) Run() {
	var wg sync.WaitGroup
	stopChan := make(chan struct{})

	// Start workers
	for i := 0; i < lt.NumWorkers; i++ {
		wg.Add(1)
		go lt.Worker(&wg, stopChan)
	}

	// Run test for specified duration
	time.Sleep(lt.TestDuration)

	// Signal workers to stop
	close(stopChan)

	// Wait for all workers to finish
	wg.Wait()

	// Calculate RPS
	rps := float64(lt.RequestCounter) / lt.TestDuration.Seconds()
	fmt.Printf("Total Requests: %d\n", lt.RequestCounter)
	fmt.Printf("Test Duration: %.2f seconds\n", lt.TestDuration.Seconds())
	fmt.Printf("Requests Per Second (RPS): %.2f\n", rps)
}

// RunLoadTest executes a load test with the provided parameters
// This replaces the main function to avoid conflicts
func RunLoadTest(url string, workers int, durationSeconds int) {
	duration := time.Duration(durationSeconds) * time.Second
	
	tester := NewLoadTester(
		url,
		workers,
		duration,
	)

	fmt.Printf("Starting load test with %d workers for %d seconds\n", workers, durationSeconds)
	fmt.Printf("Target URL: %s\n", url)
	
	tester.Run()
}

// Command line usage example:
// To use this from command line, you can create a CLI tool by adding this to your main.go
/*
	if len(os.Args) > 1 && os.Args[1] == "loadtest" {
		url := "http://localhost:1440/api/users"
		workers := 10
		duration := 10
		
		// Parse custom flags if provided
		if len(os.Args) > 2 {
			url = os.Args[2]
		}
		if len(os.Args) > 3 {
			workers, _ = strconv.Atoi(os.Args[3])
		}
		if len(os.Args) > 4 {
			duration, _ = strconv.Atoi(os.Args[4])
		}
		
		RunLoadTest(url, workers, duration)
		return
	}
*/
