package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Janina2021kisu."
	dbname   = "merck"
	fileName = "baseProcess.txt"
)

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

// getProcessName prompts the user to enter a process name.
func getProcessName() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the process name: ")
	processName, _ := reader.ReadString('\n')
	return strings.TrimSpace(processName)
}

// runSaveBaseData executes saveBaseData.go as a separate process.
func runSaveBaseData() {
	cmd := exec.Command("go", "run", "saveBaseData.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running saveBaseData.go...")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error running saveBaseData.go: %v\n", err)
	}
	fmt.Println("saveBaseData.go executed successfully.")
}

// getMaxProcessFromFile reads baseProcess.txt and determines the maximum process value.
func getMaxProcessFromFile(fileName string) (int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, fmt.Errorf("failed to open file %s: %v", fileName, err)
	}
	defer file.Close()

	maxProcess := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Assuming the file has a line like: "HierarchyID: 1, Process: 1, ..."
		parts := strings.Split(line, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "Process:") {
				// Extract the number by removing "Process:" prefix
				processStr := strings.TrimSpace(strings.Split(part, ":")[1])
				process, err := strconv.Atoi(processStr)
				if err != nil {
					return 0, fmt.Errorf("failed to convert process to integer: %v (line: %s)", err, line)
				}
				if process > maxProcess {
					maxProcess = process
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file %s: %v", fileName, err)
	}

	return maxProcess, nil
}

// insertProcessToDB inserts a new process into the database.
func insertProcessToDB(db *sql.DB, process int, label string) error {
	query := `
		INSERT INTO table_hierarchy (process, stage, operation, action, measure_id, label)
		VALUES ($1, NULL, NULL, NULL, NULL, $2)
	`
	_, err := db.Exec(query, process, label)
	if err != nil {
		return fmt.Errorf("failed to insert new process into database: %v", err)
	}
	return nil
}

// Function to read lines from a file and return them as a slice of strings
func readLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// Function to extract a specific field's value from a line
func extractFieldValue(line, field string) string {
	parts := strings.Split(line, ", ")
	for _, part := range parts {
		if strings.HasPrefix(part, field) {
			return strings.TrimPrefix(part, field+": ")
		}
	}
	return ""
}

// Function to randomly select `n` elements from a slice
func randomSelect(slice []string, n int) []string {
	rand.Seed(time.Now().UnixNano())
	selected := make([]string, n)
	indices := rand.Perm(len(slice))[:n]
	for i, idx := range indices {
		selected[i] = slice[idx]
	}
	return selected
}

// insertHierarchyValues inserts the selected values (process, stage, operation, action, measure) into the table_hierarchy
func insertHierarchyValues(db *sql.DB, process, stage, operation, action, measureID, label string) error {
	// Prepare the SQL query to insert values into the table_hierarchy
	query := `
		INSERT INTO table_hierarchy (process, stage, operation, action, measure_id, label)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := db.Exec(query, process, stage, operation, action, measureID, label)
	if err != nil {
		return fmt.Errorf("failed to insert values into database: %v", err)
	}
	return nil
}

// Function to assign 5 random stages, 4 random operations, 5 random actions, and 5 random measure_ids to a specific process
func assignRandomStagesOperationsActionsMeasures(db *sql.DB, processes, stages, operations, actions, measures []string, targetProcess string) {
	for _, processLine := range processes {
		processLabel := extractFieldValue(processLine, "Label") // Extract process label

		// Only process the specific process entered by the user
		if strings.EqualFold(processLabel, targetProcess) {
			processValue := extractFieldValue(processLine, "Process") // Extract process value

			// Randomly select 5 stages for the current process
			selectedStages := randomSelect(stages, 5)
			for _, stageLine := range selectedStages {
				stageValue := extractFieldValue(stageLine, "Stage") // Extract stage value
				stageLabel := extractFieldValue(stageLine, "Label") // Extract stage label

				// For each stage, randomly select 4 operations
				selectedOperations := randomSelect(operations, 4)
				for _, operationLine := range selectedOperations {
					operationValue := extractFieldValue(operationLine, "Operation") // Extract operation value
					operationLabel := extractFieldValue(operationLine, "Label")     // Extract operation label

					// For each operation, randomly select 5 actions
					selectedActions := randomSelect(actions, 5)
					for _, actionLine := range selectedActions {
						actionValue := extractFieldValue(actionLine, "Action") // Extract action value
						actionLabel := extractFieldValue(actionLine, "Label")  // Extract action label

						// For each action, randomly select 5 measures
						selectedMeasures := randomSelect(measures, 5)
						for _, measureLine := range selectedMeasures {
							measureValue := extractFieldValue(measureLine, "MeasureID") // Extract measure_id value
							measureLabel := extractFieldValue(measureLine, "Label")     // Extract measure label

							// Concatenate labels and create a combined label
							combinedLabel := fmt.Sprintf("%s-%s-%s-%s-%s", processLabel, stageLabel, operationLabel, actionLabel, measureLabel)

							// Insert into the database
							err := insertHierarchyValues(db, processValue, stageValue, operationValue, actionValue, measureValue, combinedLabel)
							if err != nil {
								log.Printf("Error inserting values into database: %v", err)
								continue // Skip the current iteration if there is an error
							}

							// Print confirmation
							fmt.Printf("Inserted Process: %s, Stage: %s, Operation: %s, Action: %s, MeasureID: %s, Combined Label: %s\n",
								processValue, stageValue, operationValue, actionValue, measureValue, combinedLabel)
						}
					}
				}
			}
		}
	}
}

// processBaseDataFlow performs the complete process of regenerating base data,
// finding the max process, inserting a new process, and assigning random stages,
// operations, actions, and measures.
func populateDatabase(processName string) error {
	// Connect to the database
	db, err := connectToDB()
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Run saveBaseData.go to regenerate baseProcess.txt
	runSaveBaseData()

	// Find the highest process value in baseProcess.txt
	maxProcess, err := getMaxProcessFromFile(fileName)
	if err != nil {
		return fmt.Errorf("error finding max process: %v", err)
	}
	fmt.Printf("The highest process value in %s is: %d\n", fileName, maxProcess)

	// Increment maxProcess and insert into the database
	newProcess := maxProcess + 1
	err = insertProcessToDB(db, newProcess, processName)
	if err != nil {
		return fmt.Errorf("error inserting new process into the database: %v", err)
	}
	fmt.Printf("Successfully inserted new process: %d with label: %s\n", newProcess, processName)

	// Run saveBaseData.go again to ensure the file is updated
	runSaveBaseData()

	// Read lines from files
	processes, err := readLines("baseProcess.txt")
	if err != nil {
		return fmt.Errorf("error reading baseProcess.txt: %v", err)
	}

	stages, err := readLines("baseStage.txt")
	if err != nil {
		return fmt.Errorf("error reading baseStage.txt: %v", err)
	}

	operations, err := readLines("baseOperation.txt")
	if err != nil {
		return fmt.Errorf("error reading baseOperation.txt: %v", err)
	}

	actions, err := readLines("baseAction.txt")
	if err != nil {
		return fmt.Errorf("error reading baseAction.txt: %v", err)
	}

	measures, err := readLines("baseMeasure_id.txt")
	if err != nil {
		return fmt.Errorf("error reading baseMeasure_id.txt: %v", err)
	}

	// Assign random stages, operations, actions, and measures to the specified process
	assignRandomStagesOperationsActionsMeasures(db, processes, stages, operations, actions, measures, processName)

	return nil
}

func getMeasureIDsForProcess(db *sql.DB, process string) (status string, duration string, results string, errMsg string) {
	query := `
	SELECT 
		measure_id, 
		label
	FROM 
		table_hierarchy
	WHERE 
		process = (
			SELECT process
			FROM table_hierarchy 
			WHERE label = $1
			LIMIT 1
		)
		AND measure_id IS NOT NULL;
	`

	// Record the start time
	startTime := time.Now()

	// Execute the query
	rows, err := db.Query(query, process)
	if err != nil {
		errMsg = fmt.Sprintf("Error executing query: %v", err)
		return "failed", "", "", errMsg
	}
	defer rows.Close()
	// Measure execution time
	durationTime := time.Since(startTime)

	// Process query results
	var resultBuilder strings.Builder // To construct the results string
	count := 0

	for rows.Next() {
		var measureID, label string
		if err := rows.Scan(&measureID, &label); err != nil {
			errMsg = fmt.Sprintf("Error scanning row: %v", err)
			return "failed", "", "", errMsg
		}
		// Append the current row to the results
		resultBuilder.WriteString(fmt.Sprintf("Measure ID: %s, Label: %s\n", measureID, label))
		count++
	}

	if err := rows.Err(); err != nil {
		errMsg = fmt.Sprintf("Error iterating rows: %v", err)
		return "failed", "", "", errMsg
	}

	// Prepare output
	status = "OK"
	duration = fmt.Sprintf("%v", durationTime)
	results = resultBuilder.String()

	return status, duration, results, ""
}

func main() {

	// Connect to the database
	db, err := connectToDB()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer db.Close()
	fmt.Println("Choose an action:")
	fmt.Println("1: Enter process name and populate database")
	fmt.Println("2: Query Measure IDs and Labels for a process")
	fmt.Print("Enter your choice: ")

	var choice int
	_, err = fmt.Scan(&choice)
	if err != nil {
		log.Fatalf("Invalid input: %v", err)
	}

	// Switch statement for handling choices
	switch choice {
	case 1:
		processName := getProcessName()
		fmt.Printf("You entered process name: %s\n", processName)

		// Handle database population
		if err := populateDatabase(processName); err != nil {
			log.Fatalf("Error during process base data flow: %v\n", err)
		}
	case 2:
		// Prompt the user for a process ID
		var process string
		fmt.Print("Enter process ID: ")
		_, err := fmt.Scan(&process)
		if err != nil {
			log.Fatalf("Invalid input: %v", err)
		}

		// Execute the query and display results
		status, duration, results, errMsg := getMeasureIDsForProcess(db, process)
		if status == "failed" {
			fmt.Printf("Error: %s\n", errMsg)
		} else {
			fmt.Printf("Status: %s\nDuration: %s\nResults:\n%s", status, duration, results)
		}
	}
}
