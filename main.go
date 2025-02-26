package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"strings"
// 	"time"

// 	_ "github.com/go-sql-driver/mysql"
// )

// const (
// 	numRows       = 10000000 // 10 lakh rows
// 	batchSize     = 1000    // number of rows per batch insert
// 	dataKeyLength = 100     // length of the random dataKey string
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

// func main() {
// 	// Seed the random number generator.
// 	rand.Seed(time.Now().UnixNano())

// 	// Data source name: change password and host as needed.
// 	dsn := "root:password@tcp(localhost:3306)/localDB"
// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatal("Error connecting to DB:", err)
// 	}
// 	defer db.Close()

// 	// Test the database connection.
// 	if err := db.Ping(); err != nil {
// 		log.Fatal("Error pinging DB:", err)
// 	}

// 	// Fixed date for all rows.
// 	dateValue := "2025-02-13"

// 	// Calculate the number of batches.
// 	batches := numRows / batchSize

// 	startTime := time.Now()
// 	for i := 0; i < batches; i++ {
// 		var query strings.Builder
// 		query.WriteString("INSERT INTO TxVolRiskCheckData (dataKey, value, date) VALUES ")
// 		// A slice to hold the arguments for the placeholders.
// 		values := make([]interface{}, 0, batchSize*3)
// 		for j := 0; j < batchSize; j++ {
// 			if j > 0 {
// 				query.WriteString(",")
// 			}
// 			// Each row has three placeholders: (?, ?, ?)
// 			query.WriteString("(?, ?, ?)")
// 			// Generate a random dataKey of length 100.
// 			dataKey := randString(dataKeyLength)
// 			// Generate a random float value, here in the range 0.00 to 99.99.
// 			value := float64(rand.Intn(10000)) / 100.0
// 			values = append(values, dataKey, value, dateValue)
// 		}

// 		// Execute the batch insert.
// 		_, err := db.Exec(query.String(), values...)
// 		if err != nil {
// 			log.Fatalf("Error executing batch insert at batch %d: %v", i, err)
// 		}

// 		// Optional: print progress every 10 batches.
// 		if (i+1)%10 == 0 {
// 			fmt.Printf("Inserted %d rows...\n", (i+1)*batchSize)
// 		}
// 	}

// 	elapsed := time.Since(startTime)
// 	fmt.Printf("Finished inserting %d rows in %s\n", numRows, elapsed)
// }

