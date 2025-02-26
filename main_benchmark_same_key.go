package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"sort"
// 	"sync"
// 	"time"

// 	_ "github.com/go-sql-driver/mysql"
// )

// const (
// 	tps      = 300               // 50 queries per second
// 	duration = 30 * time.Second // benchmark duration
// 	keyLen   = 100              // fixed key length (100 characters)
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

// // randomDate returns a random date string in the year 2025.
// func randomDate() string {
// 	start := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
// 	daysToAdd := rand.Intn(365)
// 	return start.AddDate(0, 0, daysToAdd).Format("2006-01-02")
// }

// // benchmarkQueryFixedKey fires queries at the given TPS for the specified duration.
// // Each query uses the same fixed key but inserts random values and random dates.
// // It returns a slice of the query latencies.
// func benchmarkQueryFixedKey(db *sql.DB, fixedKey string, tps int, duration time.Duration) ([]time.Duration, error) {
// 	latencyCh := make(chan time.Duration, 10000)
// 	var wg sync.WaitGroup

// 	// Ticker to trigger queries at the given TPS.
// 	ticker := time.NewTicker(time.Second / time.Duration(tps))
// 	done := time.After(duration)
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

// 					// For each query, generate a random value and a random date.
// 					randomVal := float64(rand.Intn(10000)) / 100.0
// 					randomDateStr := randomDate()

// 					// Build the query.
// 					query := "INSERT INTO TxVolRiskCheckData (dataKey, value, date) VALUES (?, ?, ?) " +
// 						"ON DUPLICATE KEY UPDATE value = value+1"
// 					_, err := db.Exec(query, fixedKey, randomVal, randomDateStr)
// 					elapsed := time.Since(startTime)
// 					if err != nil {
// 						log.Printf("Query error: %v", err)
// 					}
// 					latencyCh <- elapsed
// 				}()
// 			}
// 		}
// 	}()

// 	// Wait until all queries are done.
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

// 	// Generate one fixed key for all queries.
// 	fixedKey := randString(keyLen)
// 	fmt.Printf("Using fixed key: %s\n", fixedKey)

// 	// Connect to MySQL.
// 	// Update the DSN with your credentials and database details.
// 	dsn := "root:password@tcp(localhost:3306)/localDB"
// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatalf("Error connecting to DB: %v", err)
// 	}
// 	defer db.Close()

// 	if err := db.Ping(); err != nil {
// 		log.Fatalf("Error pinging DB: %v", err)
// 	}
// 	db.SetMaxIdleConns(0)
// 	db.SetMaxOpenConns(10)
// 	// Run the benchmark.
// 	latencies, err := benchmarkQueryFixedKey(db, fixedKey, tps, duration)
// 	if err != nil {
// 		log.Fatalf("Benchmark error: %v", err)
// 	}

// 	if len(latencies) == 0 {
// 		fmt.Println("No queries executed.")
// 		return
// 	}

// 	// Compute min, max, average, and median latencies.
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
