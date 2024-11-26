package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Janina2021kisu."
	dbname   = "merck"
)

///each process will have 5 stage, 4 operation, 5 action, and 5 measure

// connectToDB establishes a connection to the PostgreSQL database.
func connectToDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}
	return db, nil
}

// queryAndWriteToFile executes a query and writes the results to a file.
func queryAndWriteToFile(db *sql.DB, query, fileName string) error {
	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	// Open the file
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", fileName, err)
	}
	defer file.Close()

	// Iterate through the result rows and write to the file
	for rows.Next() {
		var hierarchyID int
		var process, stage, operation, measureID sql.NullInt64
		var action, label sql.NullString

		err := rows.Scan(&hierarchyID, &process, &stage, &operation, &action, &measureID, &label)
		if err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}

		// Write the row to the file, handling NULL values
		_, err = file.WriteString(fmt.Sprintf(
			"HierarchyID: %d, Process: %v, Stage: %v, Operation: %v, Action: %s, MeasureID: %v, Label: %s\n",
			hierarchyID, process.Int64, stage.Int64, operation.Int64,
			action.String, measureID.Int64, label.String))
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %v", fileName, err)
		}
	}

	// Check for any error encountered during iteration
	if err = rows.Err(); err != nil {
		return fmt.Errorf("row iteration error: %v", err)
	}

	return nil
}

// main function
func main() {
	// Connect to the database
	db, err := connectToDB()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer db.Close()

	// Define queries and output files
	queries := map[string]string{
		"process":    "SELECT * FROM table_hierarchy WHERE stage IS NULL AND operation IS NULL AND action IS NULL AND measure_id IS NULL",
		"stage":      "SELECT * FROM table_hierarchy WHERE process IS NULL AND operation IS NULL AND action IS NULL AND measure_id IS NULL",
		"operation":  "SELECT * FROM table_hierarchy WHERE process IS NULL AND stage IS NULL AND action IS NULL AND measure_id IS NULL",
		"action":     "SELECT * FROM table_hierarchy WHERE process IS NULL AND operation IS NULL AND stage IS NULL AND measure_id IS NULL",
		"measure_id": "SELECT * FROM table_hierarchy WHERE process IS NULL AND stage IS NULL AND action IS NULL AND operation IS NULL",
	}
	files := map[string]string{
		"process":    "baseProcess.txt",
		"stage":      "baseStage.txt",
		"operation":  "baseOperation.txt",
		"action":     "baseAction.txt",
		"measure_id": "baseMeasure_id.txt",
	}

	// Loop over the queries and write results to respective files
	for key, query := range queries {
		fileName := files[key]
		err = queryAndWriteToFile(db, query, fileName)
		if err != nil {
			log.Fatalf("Error processing %s: %v\n", key, err)
		}
	}
}
