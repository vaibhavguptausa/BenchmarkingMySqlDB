package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rcrowley/go-metrics"
)

// Benchmark configuration constants.
const (
	tps            = 300                     // 50 queries per second.
	duration       = 30 * time.Second       // Benchmark duration.
	keyLen         = 100                    // Each key is 100 characters long.
	delayThreshold = 100 * time.Millisecond // Threshold for logging delayed queries.
)

var letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// randString generates a random string of n characters.
func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// randomDate returns a random date string in the year 2025.
func randomDate() string {
	start := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	daysToAdd := rand.Intn(365)
	return start.AddDate(0, 0, daysToAdd).Format("2006-01-02")
}

// benchmarkQueryRandomKeys fires queries at the specified TPS for the given duration.
// For each query, it builds an INSERT that uses 4 rows, each row with a random key, random value, and random date.
// It records per-query latencies in a go-metrics histogram and also returns per-query details.
func benchmarkQueryRandomKeys(db *sql.DB, tps int, duration time.Duration) ([]time.Duration, error) {
	var wg sync.WaitGroup
	latencyCh := make(chan time.Duration, 10000)

	// Create a histogram for latency metrics using a uniform sample.
	histogram := metrics.NewHistogram(metrics.NewUniformSample(1028))
	metrics.Register("latency", histogram)

	ticker := time.NewTicker(time.Second / time.Duration(tps))
	done := time.After(duration)
    wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				wg.Add(1)
				go func() {
					defer wg.Done()
					startTime := time.Now()

					// Build the query with 4 rows. Each row gets its own random key.
					var builder strings.Builder
					builder.WriteString("INSERT INTO TxVolRiskCheckData (dataKey, value, date) VALUES ")
					values := make([]interface{}, 0, 12) // 4 rows * 3 columns
					for i := 0; i < 4; i++ {
						if i > 0 {
							builder.WriteString(",")
						}
						builder.WriteString("(?, ?, ?)")
						// Generate a random key of length keyLen.
						dataKey := randString(keyLen)
						// Generate a random float value.
						val := float64(rand.Intn(10000)) / 100.0
						// Generate a random date in 2025.
						dateStr := randomDate()
						values = append(values, dataKey, val, dateStr)
					}
					// Append duplicate-key update clause.
					builder.WriteString(" ON DUPLICATE KEY UPDATE value = value+1")
					query := builder.String()

					// Execute the query.
					_, err := db.Exec(query, values...)
					elapsed := time.Since(startTime)
					if err != nil {
						log.Printf("Query error: %v", err)
					}
					latencyCh <- elapsed
					histogram.Update(int64(elapsed.Nanoseconds()))

					// Log if delayed.
					if elapsed > delayThreshold {
						log.Printf("Delayed query at %s took %v; error: %v", startTime.Format(time.RFC3339Nano), elapsed, err)
					}
				}()
			}
		}
	}()

	wg.Wait()
	close(latencyCh)

	var latencies []time.Duration
	for l := range latencyCh {
		latencies = append(latencies, l)
	}

	// Report histogram metrics.
	fmt.Println("Histogram Metrics:")
	fmt.Printf("Count: %d\n", histogram.Count())
	fmt.Printf("Min: %v\n", time.Duration(histogram.Min()))
	fmt.Printf("Max: %v\n", time.Duration(histogram.Max()))
	fmt.Printf("Mean: %v\n", time.Duration(int64(histogram.Mean())))
	fmt.Printf("Median: %v\n", time.Duration(int64(histogram.Percentile(0.5))))
	fmt.Printf("95th percentile: %v\n", time.Duration(int64(histogram.Percentile(0.95))))
	fmt.Printf("99th percentile: %v\n", time.Duration(int64(histogram.Percentile(0.99))))

	return latencies, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Connect to MySQL.
	// Update the DSN with your credentials and database details.
	dsn := "root:password@tcp(localhost:3306)/localDB"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging DB: %v", err)
	}

	latencies, err := benchmarkQueryRandomKeys(db, tps, duration)
	if err != nil {
		log.Fatalf("Benchmark error: %v", err)
	}

	if len(latencies) == 0 {
		fmt.Println("No queries executed.")
		return
	}
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)
	// Compute basic statistics from the latency slice.
	var total time.Duration
	minLatency := latencies[0]
	maxLatency := latencies[0]
	for _, l := range latencies {
		total += l
		if l < minLatency {
			minLatency = l
		}
		if l > maxLatency {
			maxLatency = l
		}
	}
	avgLatency := total / time.Duration(len(latencies))
	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	medianLatency := latencies[len(latencies)/2]

	fmt.Println("Basic Latency Stats:")
	fmt.Printf("Executed %d queries\n", len(latencies))
	fmt.Printf("Min latency: %v\n", minLatency)
	fmt.Printf("Max latency: %v\n", maxLatency)
	fmt.Printf("Average latency: %v\n", avgLatency)
	fmt.Printf("Median latency: %v\n", medianLatency)

	// Optional: Print overall connection pool stats.
	stats := db.Stats()
	fmt.Printf("DB Stats: Open Connections: %d, In Use: %d, Idle: %d\n",
		stats.OpenConnections, stats.InUse, stats.Idle)
}
