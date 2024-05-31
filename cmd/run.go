/*
Copyright © 2024 Felipe Correa <dev.felipecls@gmail.com>
*/
package cmd

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	url         string
	totalReqs   int
	concurrency int
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run load test",
	Long:  `Run a load test on the specified URL with given concurrency and request count.`,
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" || totalReqs <= 0 || concurrency <= 0 {
			fmt.Println("Invalid parameters. Use --url, --requests, and --concurrency.")
			return
		}

		fmt.Printf("Starting load tests: %s\n", url)
		fmt.Printf("Total requests: %d, Concurrency: %d\n", totalReqs, concurrency)

		runLoadTest(url, totalReqs, concurrency)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVar(&url, "url", "", "URL of the service to be tested")
	runCmd.Flags().IntVar(&totalReqs, "requests", 0, "Total number of requests")
	runCmd.Flags().IntVar(&concurrency, "concurrency", 1, "Number of concurrent calls")
}

func runLoadTest(url string, totalReqs, concurrency int) {
	var wg sync.WaitGroup
	results := make(chan int, totalReqs)

	startTime := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < totalReqs/concurrency; j++ {
				status := makeRequest(url)
				results <- status
			}
		}()
	}

	wg.Wait()
	close(results)

	totalTime := time.Since(startTime)
	generateReport(results, totalTime)
}

func makeRequest(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func generateReport(results chan int, totalTime time.Duration) {
	total := 0
	status200 := 0
	statusCount := make(map[int]int)

	for status := range results {
		total++
		if status == 200 {
			status200++
		}
		statusCount[status]++
	}

	fmt.Printf("Load Test Report:\n")
	fmt.Printf("Total Time: %v\n", totalTime)
	fmt.Printf("Total Requests: %d\n", total)
	fmt.Printf("Requests with Status 200: %d\n", status200)

	fmt.Println("HTTP Status Distribution:")
	for status, count := range statusCount {
		fmt.Printf("Status %d: %d\n", status, count)
	}
}
