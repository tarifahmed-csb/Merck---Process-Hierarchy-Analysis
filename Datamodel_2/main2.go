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
	var printHierarchy = `
	SELECT ps.process, ps.stage, so.operation, oa.action
	FROM process_stage ps
	JOIN stage_operation so ON ps.stage = so.stage
	LEFT JOIN operation_action oa ON so.operation = oa.operation
	WHERE ps.process = 'process_1';
	`

	var printProcessMeasurement = `
	SELECT measurement, measurement_id, process
	FROM measurements 
	WHERE process=$1;
	`
	var printStageMeasurement = `
	SELECT measurement, measurement_id, stage
	FROM measurements 
	WHERE stage=$1;
	`
	var printOperationMeasurement = `
	SELECT measurement, measurement_id, operation
	FROM measurements 
	WHERE operation=$1;
	`
	var printActionMeasurement = `
	SELECT measurement, measurement_id, action
	FROM measurements 
	WHERE action=$1;
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
	PrintDistinctValues("Processes:", processMap)
	PrintDistinctValues("Stages:", stageMap)
	PrintDistinctValues("Operations:", operationMap)
	PrintDistinctValues("Actions:", actionMap)

	// Print measurement data for each distinct process
	PrintMeasurementData(db, printProcessMeasurement, "process", processMap)

	// Print measurement data for each distinct stage
	PrintMeasurementData(db, printStageMeasurement, "stage", stageMap)

	// Print measurement data for each distinct operation
	PrintMeasurementData(db, printOperationMeasurement, "operation", operationMap)

	// Print measurement data for each distinct action
	PrintMeasurementData(db, printActionMeasurement, "action", actionMap)

	elapsed := time.Since(start)
	fmt.Printf("Query took %s\n", elapsed)
}

// PrintDistinctValues prints distinct values stored in a map
func PrintDistinctValues(label string, m map[string]bool) {
	fmt.Println(label, GetKeys(m))
}

// PrintMeasurementData prints measurement data for each distinct value
func PrintMeasurementData(db *sql.DB, query, valueType string, valueMap map[string]bool) {
	for value := range valueMap {
		rows, err := db.Query(query, value)
		CheckError(err)
		defer rows.Close()

		fmt.Printf("Measurement data for %s %s:\n", valueType, value)
		for rows.Next() {
			var measurement string
			var measurementID string
			var v string
			err := rows.Scan(&measurement, &measurementID, &v)
			CheckError(err)
			fmt.Printf("Measurement: %s, Measurement ID: %s, %s: %s\n", measurement, measurementID, valueType, v)
		}
	}
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
