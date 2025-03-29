package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Request types for mixed benchmarking
type RequestType struct {
	Method   string
	Endpoint string
	Body     interface{}
}

// BenchmarkEndpoint runs a benchmark against a specified endpoint
func BenchmarkEndpoint(b *testing.B) {
	endpoints := []string{
		"http://localhost:1440/api/users",
		"http://localhost:1440/api/collection-points",
		"http://localhost:1440/api/waste-types",
		"http://localhost:1440/api/reports",
	}

	concurrencyLevels := []int{1, 5, 10, 25, 50, 100}

	for _, endpoint := range endpoints {
		for _, concurrency := range concurrencyLevels {
			b.Run(fmt.Sprintf("%s-Concurrency_%d", getEndpointName(endpoint), concurrency), func(b *testing.B) {
				benchmarkWithConcurrency(b, endpoint, concurrency)
			})
		}
	}
}

func getEndpointName(url string) string {
	// Extract the last part of the URL path
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == '/' {
			return url[i+1:]
		}
	}
	return url
}

func benchmarkWithConcurrency(b *testing.B, url string, concurrency int) {
	// Skip the timer during setup
	b.ResetTimer()
	b.SetParallelism(concurrency)

	var wg sync.WaitGroup
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: concurrency * 2,
			MaxConnsPerHost:     concurrency * 2,
		},
	}

	requestsCompleted := int64(0)
	requestsSuccessful := int64(0)

	// Start test with a fixed number of iterations
	iterations := b.N
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for atomic.AddInt64(&requestsCompleted, 1) <= int64(iterations) {
				resp, err := client.Get(url)
				if err == nil {
					resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						atomic.AddInt64(&requestsSuccessful, 1)
					}
				}
			}
		}()
	}

	wg.Wait()

	b.ReportMetric(float64(requestsSuccessful)/float64(iterations), "success_rate")
}

// BenchmarkRPSMeasurement runs a fixed-duration test to measure RPS
func BenchmarkRPSMeasurement(b *testing.B) {
	endpoints := []string{
		"http://localhost:1440/api/users",
		"http://localhost:1440/api/collection-points",
		"http://localhost:1440/api/waste-types",
		"http://localhost:1440/api/reports",
	}

	for _, endpoint := range endpoints {
		b.Run(fmt.Sprintf("RPS_%s", getEndpointName(endpoint)), func(b *testing.B) {
			// Configuration
			concurrency := 50
			testDuration := 5 * time.Second
			
			// Reset timer to exclude setup time
			b.ResetTimer()
			
			// Setup HTTP client
			client := &http.Client{
				Timeout: 2 * time.Second,
				Transport: &http.Transport{
					MaxIdleConnsPerHost: concurrency,
					MaxConnsPerHost:     concurrency,
				},
			}
			
			var wg sync.WaitGroup
			var requestCount int64
			
			// Signal channel to stop workers
			stopCh := make(chan struct{})
			
			// Start workers
			wg.Add(concurrency)
			for i := 0; i < concurrency; i++ {
				go func() {
					defer wg.Done()
					for {
						select {
						case <-stopCh:
							return
						default:
							resp, err := client.Get(endpoint)
							if err == nil {
								resp.Body.Close()
								if resp.StatusCode == http.StatusOK {
									atomic.AddInt64(&requestCount, 1)
								}
							}
						}
					}
				}()
			}
			
			// Run test for specified duration
			time.Sleep(testDuration)
			close(stopCh)
			wg.Wait()
			
			// Calculate RPS
			rps := float64(requestCount) / testDuration.Seconds()
			b.ReportMetric(rps, "requests/sec")
			fmt.Printf("Endpoint: %s - Total Requests: %d, RPS: %.2f\n", 
				endpoint, requestCount, rps)
		})
	}
}

// BenchmarkMixedRequests tests a mix of GET, POST, PUT, and DELETE operations
func BenchmarkMixedRequests(b *testing.B) {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())
	
	// Sample data for POST/PUT requests
	sampleUser := map[string]interface{}{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
	}
	
	sampleReport := map[string]interface{}{
		"title":       "Test Report",
		"description": "This is a test report",
		"location":    "Test Location",
	}
	
	sampleWasteType := map[string]interface{}{
		"name":        "Test Waste",
		"description": "This is a test waste type",
	}
	
	// Define a mix of different request types
	requestTypes := []RequestType{
		// GET requests
		{Method: "GET", Endpoint: "http://localhost:1440/api/users"},
		{Method: "GET", Endpoint: "http://localhost:1440/api/reports"},
		{Method: "GET", Endpoint: "http://localhost:1440/api/waste-types"},
		{Method: "GET", Endpoint: "http://localhost:1440/api/collection-points"},
		
		// POST requests
		{Method: "POST", Endpoint: "http://localhost:1440/api/users", Body: sampleUser},
		{Method: "POST", Endpoint: "http://localhost:1440/api/reports", Body: sampleReport},
		{Method: "POST", Endpoint: "http://localhost:1440/api/waste-types", Body: sampleWasteType},
		
		// PUT requests (assuming IDs 1-5 exist)
		{Method: "PUT", Endpoint: "http://localhost:1440/api/users/1", Body: sampleUser},
		{Method: "PUT", Endpoint: "http://localhost:1440/api/reports/1", Body: sampleReport},
		
		// DELETE requests (assuming IDs exist)
		{Method: "DELETE", Endpoint: "http://localhost:1440/api/reports/2"},
		{Method: "DELETE", Endpoint: "http://localhost:1440/api/users/3"},
	}
	
	concurrencyLevels := []int{10, 50, 100}
	
	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("MixedRequests_Concurrency_%d", concurrency), func(b *testing.B) {
			// Reset timer to exclude setup time
			b.ResetTimer()
			
			// Setup HTTP client
			client := &http.Client{
				Timeout: 3 * time.Second,
				Transport: &http.Transport{
					MaxIdleConnsPerHost: concurrency * 2,
					MaxConnsPerHost:     concurrency * 2,
				},
			}
			
			var wg sync.WaitGroup
			var totalRequests int64
			var successfulRequests int64
			
			// Start test
			iterations := b.N
			wg.Add(concurrency)
			
			for i := 0; i < concurrency; i++ {
				go func() {
					defer wg.Done()
					for atomic.AddInt64(&totalRequests, 1) <= int64(iterations) {
						// Choose a random request type
						requestType := requestTypes[rand.Intn(len(requestTypes))]
						
						var req *http.Request
						var err error
						
						switch requestType.Method {
						case "GET", "DELETE":
							req, err = http.NewRequest(requestType.Method, requestType.Endpoint, nil)
						case "POST", "PUT":
							jsonData, _ := json.Marshal(requestType.Body)
							req, err = http.NewRequest(requestType.Method, requestType.Endpoint, bytes.NewBuffer(jsonData))
							req.Header.Set("Content-Type", "application/json")
						}
						
						if err != nil {
							continue
						}
						
						resp, err := client.Do(req)
						if err == nil {
							resp.Body.Close()
							if resp.StatusCode >= 200 && resp.StatusCode < 300 {
								atomic.AddInt64(&successfulRequests, 1)
							}
						}
					}
				}()
			}
			
			wg.Wait()
			
			// Report success rate
			b.ReportMetric(float64(successfulRequests)/float64(iterations), "success_rate")
		})
	}
}

// BenchmarkRealWorldScenario simulates a more realistic usage pattern
func BenchmarkRealWorldScenario(b *testing.B) {
	testDuration := 10 * time.Second
	
	b.Run("RealisticUsage", func(b *testing.B) {
		// Reset timer to exclude setup time
		b.ResetTimer()
		
		// Number of simulated users with different behavior patterns
		numBrowserUsers := 30  // Users browsing/reading
		numReporterUsers := 15 // Users creating reports/content
		numAdminUsers := 5     // Users performing admin operations
		
		client := &http.Client{
			Timeout: 3 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 200,
				MaxConnsPerHost:     200,
			},
		}
		
		var wg sync.WaitGroup
		var browserCount, reporterCount, adminCount int64
		
		// Signal channel to stop workers
		stopCh := make(chan struct{})
		
		// Browser users - mostly GET requests
		wg.Add(numBrowserUsers)
		for i := 0; i < numBrowserUsers; i++ {
			go func() {
				defer wg.Done()
				browsePaths := []string{
					"http://localhost:1440/api/reports",
					"http://localhost:1440/api/collection-points",
					"http://localhost:1440/api/waste-types",
				}
				
				for {
					select {
					case <-stopCh:
						return
					default:
						endpoint := browsePaths[rand.Intn(len(browsePaths))]
						resp, err := client.Get(endpoint)
						if err == nil {
							resp.Body.Close()
							if resp.StatusCode == http.StatusOK {
								atomic.AddInt64(&browserCount, 1)
							}
						}
						// Simulate user think time
						// time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
					}
				}
			}()
		}
		
		// Reporter users - mix of GET and POST
		wg.Add(numReporterUsers)
		for i := 0; i < numReporterUsers; i++ {
			go func() {
				defer wg.Done()
				
				reportData := map[string]interface{}{
					"title":       fmt.Sprintf("Report %d", rand.Intn(1000)),
					"description": "Test description",
					"location":    "Test location",
				}
				
				for {
					select {
					case <-stopCh:
						return
					default:
						// 20% chance to create a report, 80% to browse
						if rand.Float32() < 0.2 {
							// Create report
							jsonData, _ := json.Marshal(reportData)
							req, _ := http.NewRequest("POST", "http://localhost:1440/api/reports", bytes.NewBuffer(jsonData))
							req.Header.Set("Content-Type", "application/json")
							resp, err := client.Do(req)
							if err == nil {
								resp.Body.Close()
								if resp.StatusCode < 300 {
									atomic.AddInt64(&reporterCount, 1)
								}
							}
							// Longer delay after posting
							time.Sleep(time.Duration(1000+rand.Intn(1000)) * time.Millisecond)
						} else {
							// Browse reports
							resp, err := client.Get("http://localhost:1440/api/reports")
							if err == nil {
								resp.Body.Close()
								if resp.StatusCode == http.StatusOK {
									atomic.AddInt64(&reporterCount, 1)
								}
							}
							// time.Sleep(time.Duration(300+rand.Intn(700)) * time.Millisecond)
						}
					}
				}
			}()
		}
		
		// Admin users - all kinds of operations
		wg.Add(numAdminUsers)
		for i := 0; i < numAdminUsers; i++ {
			go func() {
				defer wg.Done()
				
				userData := map[string]interface{}{
					"name":     fmt.Sprintf("User %d", rand.Intn(1000)),
					"email":    fmt.Sprintf("user%d@example.com", rand.Intn(10000)),
					"password": "password123",
				}
				
				adminOps := []string{"GET", "POST", "PUT", "DELETE"}
				adminEndpoints := []string{
					"http://localhost:1440/api/users",
					"http://localhost:1440/api/waste-types",
					"http://localhost:1440/api/collection-points",
				}
				
				for {
					select {
					case <-stopCh:
						return
					default:
						op := adminOps[rand.Intn(len(adminOps))]
						baseEndpoint := adminEndpoints[rand.Intn(len(adminEndpoints))]
						endpoint := baseEndpoint
						
						// Add ID for PUT/DELETE
						if op == "PUT" || op == "DELETE" {
							endpoint = fmt.Sprintf("%s/%d", baseEndpoint, 1+rand.Intn(5))
						}
						
						var req *http.Request
						if op == "POST" || op == "PUT" {
							jsonData, _ := json.Marshal(userData)
							req, _ = http.NewRequest(op, endpoint, bytes.NewBuffer(jsonData))
							req.Header.Set("Content-Type", "application/json")
						} else {
							req, _ = http.NewRequest(op, endpoint, nil)
						}
						
						resp, err := client.Do(req)
						if err == nil {
							resp.Body.Close()
							if resp.StatusCode < 300 {
								atomic.AddInt64(&adminCount, 1)
							}
						}
						
						// Admin operations have variable delays
						// time.Sleep(time.Duration(500+rand.Intn(1500)) * time.Millisecond)
					}
				}
			}()
		}
		
		// Run test for specified duration
		time.Sleep(testDuration)
		close(stopCh)
		wg.Wait()
		
		// Calculate and report metrics
		totalRequests := browserCount + reporterCount + adminCount
		browserRPS := float64(browserCount) / testDuration.Seconds()
		reporterRPS := float64(reporterCount) / testDuration.Seconds()
		adminRPS := float64(adminCount) / testDuration.Seconds()
		totalRPS := float64(totalRequests) / testDuration.Seconds()
		
		b.ReportMetric(browserRPS, "browser_rps")
		b.ReportMetric(reporterRPS, "reporter_rps")
		b.ReportMetric(adminRPS, "admin_rps")
		b.ReportMetric(totalRPS, "total_rps")
		
		fmt.Printf("Realistic Test Results:\n")
		fmt.Printf("Browser Users: %d requests (%.2f RPS)\n", browserCount, browserRPS)
		fmt.Printf("Reporter Users: %d requests (%.2f RPS)\n", reporterCount, reporterRPS)
		fmt.Printf("Admin Users: %d requests (%.2f RPS)\n", adminCount, adminRPS)
		fmt.Printf("Total: %d requests (%.2f RPS)\n", totalRequests, totalRPS)
	})
}
