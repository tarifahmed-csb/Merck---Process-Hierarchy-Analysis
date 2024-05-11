package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Janina2021kisu."
	dbname   = "database2"
)

type HierarchyRow struct {
	Process   string
	Stage     string
	Operation string
	Action    string
}

func main() {
	// Declare connection
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()

	// Query for hierarchy
	var printHierarchy string = `
	SELECT ps.process, ps.stage, so.operation, oa.action
	FROM process_stage ps
	JOIN stage_operation so ON ps.stage = so.stage
	LEFT JOIN operation_action oa ON so.operation = oa.operation
	WHERE ps.process = 'process_2';
	`

	//measurement data
	var printProcessMeasurement string = `
	select measurement measurement_id process
	from measurements 
	where process=$1;
	`
	var printStageMeasurement string = `
	select measurement measurement_id stage
	from measurements 
	where stage=$1;
	`

	// Measure query performance
	start := time.Now()

	// Execute query
	rows, err := db.Query(printHierarchy)
	CheckError(err)
	defer rows.Close()

	// Create maps to store distinct values
	processMap := make(map[string]bool)
	stageMap := make(map[string]bool)
	operationMap := make(map[string]bool)
	actionMap := make(map[string]bool)

	// Iterate over the results and store distinct values
	for rows.Next() {
		var row HierarchyRow
		err := rows.Scan(&row.Process, &row.Stage, &row.Operation, &row.Action)
		CheckError(err)

		// Store distinct values in maps
		processMap[row.Process] = true
		stageMap[row.Stage] = true
		operationMap[row.Operation] = true
		actionMap[row.Action] = true
	}

	// Print distinct values
	fmt.Println("Distinct Processes:", GetKeys(processMap))
	fmt.Println("Distinct Stages:", GetKeys(stageMap))
	fmt.Println("Distinct Operations:", GetKeys(operationMap))
	fmt.Println("Distinct Actions:", GetKeys(actionMap))

	// Print measurement data for each distinct process
	for process := range processMap {
		rows, err := db.Query(printProcessMeasurement, process)
		CheckError(err)
		defer rows.Close()

		fmt.Printf("Measurement data for process %s:\n", process)
		for rows.Next() {
			var measurement string
			var measurementID string
			var process string
			err := rows.Scan(&measurement, &measurementID, &process)
			CheckError(err)
			fmt.Printf("Measurement: %s, Measurement ID: %s, Process: %s\n", measurement, measurementID, process)
		}
	}

	// Print measurement data for each distinct stage
	for stage := range stageMap {
		rows, err := db.Query(printStageMeasurement, stage)
		CheckError(err)
		defer rows.Close()

		fmt.Printf("Measurement data for stage %s:\n", stage)
		for rows.Next() {
			var measurement string
			var measurementID string
			var stage string
			err := rows.Scan(&measurement, &measurementID, &stage)
			CheckError(err)
			fmt.Printf("Measurement: %s, Measurement ID: %s, Stage: %s\n", measurement, measurementID, stage)
		}
	}

	// Similar logic for operation and action...

	elapsed := time.Since(start)
	fmt.Printf("Query took %s\n", elapsed)
}

// GetKeys returns the keys of a map as a slice
func GetKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// CheckError function to handle errors
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
