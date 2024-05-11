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
	dbname   = "postgres"
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
	SELECT 
		process_stages.process, 
		process_stages.stage, 
		stages_operations.operation AS operation,
		operations_actions.action
	FROM 
		process_stages
	INNER JOIN 
		stages_operations ON process_stages.stage = stages_operations.stage
	INNER JOIN 
		operations_actions ON stages_operations.operation = operations_actions.operation
	WHERE 
		process_stages.process = 'napa';
	`

	// Query for stage measurements
	var printStageMeasurements string = `SELECT * FROM stages_measurements WHERE stage = $1`

	//query for operation measurments
	var printOperationMeasurements string = `Select * from operations_measurements Where operation = $1`

	//query for action measurement
	var printactionMeasurements string = `Select * from actions_measurements where action = $1`

	//query for process measurement
	//var printProcessMeasurements string = `select * from process_measurements where process = $1`

	// Measure query execution time
	startTime := time.Now()

	// Execute the query for hierarchy
	hierarchyRows, err := queryHierarchy(db, printHierarchy)
	CheckError(err)

	// Print hierarchy table
	printHierarchyTable(hierarchyRows)

	//print process measurements

	// Execute and print stage measurements
	printStageMeasurementsTable(db, printStageMeasurements, hierarchyRows)

	// Execute and print operations measurements
	printOperationsMeasurementsTable(db, printOperationMeasurements, hierarchyRows)

	// Execute and print operations measurements
	printActionMeasurementsTable(db, printactionMeasurements, hierarchyRows)

	// Calculate and print query execution time
	elapsed := time.Since(startTime)
	fmt.Printf("Query execution time: %v\n", elapsed)
}

func queryHierarchy(db *sql.DB, query string) ([]HierarchyRow, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hierarchyRows []HierarchyRow
	for rows.Next() {
		var row HierarchyRow
		if err := rows.Scan(&row.Process, &row.Stage, &row.Operation, &row.Action); err != nil {
			return nil, err
		}
		hierarchyRows = append(hierarchyRows, row)
	}
	return hierarchyRows, nil

}

func printHierarchyTable(rows []HierarchyRow) {
	fmt.Println("Hierarchy Table:")
	fmt.Printf("%-15s %-15s %-15s %-15s\n", "Process", "Stage", "Operation", "Action")
	for _, row := range rows {
		fmt.Printf("%-15s %-15s %-15s %-15s\n", row.Process, row.Stage, row.Operation, row.Action)
	}
	fmt.Printf("\n")
}

func printStageMeasurementsTable(db *sql.DB, query string, hierarchyRows []HierarchyRow) {
	stageLabelPrinted := false // Flag to track if stage label has been printed

	// Keep track of printed measurements for each stage
	printedMeasurements := make(map[string]map[string]bool)

	for _, row := range hierarchyRows {
		stageRows, err := db.Query(query, row.Stage)
		CheckError(err)

		if !stageLabelPrinted {
			fmt.Printf("%-15s %-15s\n", "Stage", "Measurement")
			stageLabelPrinted = true // Set the flag to true after printing the label
		}

		// Initialize the map for the current stage if it doesn't exist
		if printedMeasurements[row.Stage] == nil {
			printedMeasurements[row.Stage] = make(map[string]bool)
		}

		var stage, measurement string
		for stageRows.Next() {
			if err := stageRows.Scan(&stage, &measurement); err != nil {
				CheckError(err)
				continue
			}

			// Check if the measurement has already been printed for this stage
			if !printedMeasurements[row.Stage][measurement] {
				fmt.Printf("%-15s %-15s\n", stage, measurement)
				printedMeasurements[row.Stage][measurement] = true
			}
		}
		stageRows.Close()
	}
	fmt.Printf("\n")
}

func printOperationsMeasurementsTable(db *sql.DB, query string, hierarchyRows []HierarchyRow) {
	operationLabelPrinted := false // Flag to track if operation label has been printed

	// Keep track of printed measurements for each operation
	printedMeasurements := make(map[string]map[string]bool)

	for _, row := range hierarchyRows {
		operationRows, err := db.Query(query, row.Operation)
		CheckError(err)

		if !operationLabelPrinted {
			fmt.Printf("%-15s %-15s\n", "Operation", "Measurement")
			operationLabelPrinted = true // Set the flag to true after printing the label
		}

		// Initialize the map for the current operation if it doesn't exist
		if printedMeasurements[row.Operation] == nil {
			printedMeasurements[row.Operation] = make(map[string]bool)
		}

		var operation, measurement sql.NullString
		for operationRows.Next() {
			if err := operationRows.Scan(&operation, &measurement); err != nil {
				CheckError(err)
				continue
			}

			// Check if the measurement has already been printed for this operation
			if !printedMeasurements[row.Operation][measurement.String] {
				if measurement.Valid {
					fmt.Printf("%-15s %-15s\n", row.Operation, measurement.String)
				} else {
					fmt.Printf("%-15s %-15s\n", row.Operation, "NULL")
				}
				printedMeasurements[row.Operation][measurement.String] = true
			}
		}
		operationRows.Close()
	}
	fmt.Printf("\n")
}

func printActionMeasurementsTable(db *sql.DB, query string, hierarchyRows []HierarchyRow) {
	actionLabelPrinted := false // Flag to track if action label has been printed

	// Keep track of printed measurements for each action
	printedMeasurements := make(map[string]map[string]bool)

	for _, row := range hierarchyRows {
		actionRows, err := db.Query(query, row.Action)
		CheckError(err)

		if !actionLabelPrinted {
			fmt.Printf("%-15s %-15s\n", "Action", "Measurement")
			actionLabelPrinted = true // Set the flag to true after printing the label
		}

		// Initialize the map for the current action if it doesn't exist
		if printedMeasurements[row.Action] == nil {
			printedMeasurements[row.Action] = make(map[string]bool)
		}

		var action, measurement sql.NullString
		for actionRows.Next() {
			if err := actionRows.Scan(&action, &measurement); err != nil {
				CheckError(err)
				continue
			}

			// Check if the measurement has already been printed for this action
			if !printedMeasurements[row.Action][measurement.String] {
				fmt.Printf("%-15s %-15s\n", action.String, measurement.String)
				printedMeasurements[row.Action][measurement.String] = true
			}
		}
		actionRows.Close()
	}
	fmt.Printf("\n")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
