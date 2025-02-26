package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"sort"
// 	"strings"
// 	"sync"
// 	"time"

// 	_ "github.com/go-sql-driver/mysql"
// )

// const (
// 	tps      = 300               // 50 queries per second
// 	duration = 30 * time.Second // benchmark duration
// )

// var letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// // randString generates a random string of n characters.
// func randString(n int) string {
// 	b := make([]byte, n)
// 	for i := range b {
// 		b[i] = letterBytes[rand.Intn(len(letterBytes))]
// 	}
// 	return string(b)
// }

// // randomDate generates a random date within 2025 (between Jan 1 and Dec 31).
// func randomDate() string {
// 	start := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
// 	daysToAdd := rand.Intn(365) // 0 to 364 days
// 	randomDate := start.AddDate(0, 0, daysToAdd)
// 	return randomDate.Format("2006-01-02")
// }

// // benchmarkQuery fires queries at the specified TPS for the given duration.
// // Each query is built with random keys and values. It returns a slice of latencies.
// func benchmarkQuery(db *sql.DB, tps int, duration time.Duration) ([]time.Duration, error) {
// 	latencyCh := make(chan time.Duration, 10000)
// 	var wg sync.WaitGroup

// 	ticker := time.NewTicker(time.Second / time.Duration(tps))
// 	done := time.After(duration)

// 	// Launch queries until duration expires.
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		for {
// 			select {
// 			case <-done:
// 				ticker.Stop()
// 				return
// 			case <-ticker.C:
// 				wg.Add(1)
// 				go func() {
// 					defer wg.Done()
// 					startTime := time.Now()

// 					// Build the query dynamically with 4 rows.
// 					var builder strings.Builder
// 					builder.WriteString("INSERT INTO TxVolRiskCheckData (dataKey, value, date) VALUES ")
// 					values := make([]interface{}, 0, 12) // 4 rows * 3 columns
// 					for i := 0; i < 4; i++ {
// 						if i > 0 {
// 							builder.WriteString(",")
// 						}
// 						builder.WriteString("(?, ?, ?)")
// 						// Generate random data for each row.
// 						dataKey := randString(100)
// 						val := float64(rand.Intn(10000)) / 100.0
// 						dateStr := randomDate()
// 						values = append(values, dataKey, val, dateStr)
// 					}
// 					// Append duplicate-key update clause.
// 					builder.WriteString(" ON DUPLICATE KEY UPDATE value = value+1")
// 					query := builder.String()

// 					// Execute the query.
// 					_, err := db.Exec(query, values...)
// 					elapsed := time.Since(startTime)
// 					if err != nil {
// 						log.Printf("Query error: %v", err)
// 					}
// 					latencyCh <- elapsed
// 				}()
// 			}
// 		}
// 	}()

// 	// Wait until all fired queries complete.
// 	wg.Wait()
// 	close(latencyCh)

// 	var latencies []time.Duration
// 	for l := range latencyCh {
// 		latencies = append(latencies, l)
// 	}
// 	return latencies, nil
// }

// func main() {
// 	rand.Seed(time.Now().UnixNano())

// 	// Update DSN with your credentials and database details.
// 	dsn := "root:password@tcp(localhost:3306)/localDB"
// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatalf("Error connecting to DB: %v", err)
// 	}
// 	defer db.Close()

// 	if err := db.Ping(); err != nil {
// 		log.Fatalf("Error pinging DB: %v", err)
// 	}

// 	latencies, err := benchmarkQuery(db, tps, duration)
// 	if err != nil {
// 		log.Fatalf("Benchmark error: %v", err)
// 	}

// 	if len(latencies) == 0 {
// 		fmt.Println("No queries executed.")
// 		return
// 	}

// 	// Compute statistics.
// 	minLatency := latencies[0]
// 	maxLatency := latencies[0]
// 	var total time.Duration
// 	for _, l := range latencies {
// 		if l < minLatency {
// 			minLatency = l
// 		}
// 		if l > maxLatency {
// 			maxLatency = l
// 		}
// 		total += l
// 	}
// 	avgLatency := total / time.Duration(len(latencies))

// 	// Compute median.
// 	sort.Slice(latencies, func(i, j int) bool {
// 		return latencies[i] < latencies[j]
// 	})
// 	medianLatency := latencies[len(latencies)/2]

// 	fmt.Printf("Executed %d queries\n", len(latencies))
// 	fmt.Printf("Min latency: %v\n", minLatency)
// 	fmt.Printf("Max latency: %v\n", maxLatency)
// 	fmt.Printf("Average latency: %v\n", avgLatency)
// 	fmt.Printf("Median latency: %v\n", medianLatency)
// }
