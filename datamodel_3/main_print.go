package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Janina2021kisu."
	dbname   = "merck"
)

func main() {
	// Declare connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Open the database connection
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Ping the database to ensure connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	// Updated query to select from table_hierarchy where process=1, stage=3, action=8, and measure_id=6
	query := `SELECT * FROM table_hierarchy WHERE process=1;`

	// Start measuring the time
	startTime := time.Now()

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	// Measure the elapsed time
	elapsedTime := time.Since(startTime)

	// Print the elapsed time
	fmt.Printf("Query executed in %v\n", elapsedTime)

	// Retrieve column names
	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("Failed to retrieve columns: %v", err)
	}

	// Print the column names
	fmt.Println(columns)

	// Create a slice of interface{}'s to hold each column value from a row
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for rows.Next() {
		// For each column, create a pointer to the corresponding value
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// Scan the row and fill the values slice
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		// Print each row's values
		for i, col := range columns {
			val := values[i]

			// If the column contains a null value
			if val == nil {
				fmt.Printf("%s: NULL\n", col)
			} else {
				fmt.Printf("%s: %v\n", col, val)
			}
		}
		fmt.Println("-----------------------")
	}

	// Check for errors after iterating over rows
	err = rows.Err()
	if err != nil {
		log.Fatalf("Error iterating over rows: %v", err)
	}
}
