package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

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

// Function to assign 2 random stages, 2 random operations, 2 random actions, and 2 random measure_ids to a specific process
func assignRandomStagesOperationsActionsMeasures(processes, stages, operations, actions, measures []string, targetProcess string) {
	for _, processLine := range processes {
		processLabel := extractFieldValue(processLine, "Label") // Extract process label

		// Only process the specific process entered by the user
		if strings.EqualFold(processLabel, targetProcess) {
			processValue := extractFieldValue(processLine, "Process") // Extract process value

			// Randomly select 2 stages for the current process
			selectedStages := randomSelect(stages, 2)
			for _, stageLine := range selectedStages {
				stageValue := extractFieldValue(stageLine, "Stage") // Extract stage value
				stageLabel := extractFieldValue(stageLine, "Label") // Extract stage label

				// For each stage, randomly select 2 operations
				selectedOperations := randomSelect(operations, 2)
				for _, operationLine := range selectedOperations {
					operationValue := extractFieldValue(operationLine, "Operation") // Extract operation value
					operationLabel := extractFieldValue(operationLine, "Label")     // Extract operation label

					// For each operation, randomly select 2 actions
					selectedActions := randomSelect(actions, 2)
					for _, actionLine := range selectedActions {
						actionValue := extractFieldValue(actionLine, "Action") // Extract action value
						actionLabel := extractFieldValue(actionLine, "Label")  // Extract action label

						// For each action, randomly select 2 measures
						selectedMeasures := randomSelect(measures, 2)
						for _, measureLine := range selectedMeasures {
							measureValue := extractFieldValue(measureLine, "MeasureID") // Extract measure_id value
							measureLabel := extractFieldValue(measureLine, "Label")     // Extract measure label

							// Concatenate labels and print the result
							combinedLabel := fmt.Sprintf("%s-%s-%s-%s-%s", processLabel, stageLabel, operationLabel, actionLabel, measureLabel)
							fmt.Printf("Process: %s, Stage: %s, Operation: %s, Action: %s, MeasureID: %s, Combined Label: %s\n",
								processValue, stageValue, operationValue, actionValue, measureValue, combinedLabel)
						}
					}
				}
			}
		}
	}
}

func main() {
	// Read lines from files
	processes, err := readLines("baseProcess.txt")
	if err != nil {
		fmt.Println("Error reading baseProcess.txt:", err)
		return
	}

	stages, err := readLines("baseStage.txt")
	if err != nil {
		fmt.Println("Error reading baseStage.txt:", err)
		return
	}

	operations, err := readLines("baseOperation.txt")
	if err != nil {
		fmt.Println("Error reading baseOperation.txt:", err)
		return
	}

	actions, err := readLines("baseAction.txt")
	if err != nil {
		fmt.Println("Error reading baseAction.txt:", err)
		return
	}

	measures, err := readLines("baseMeasure_id.txt")
	if err != nil {
		fmt.Println("Error reading baseMeasure_id.txt:", err)
		return
	}

	// Prompt user for specific process name
	fmt.Print("Enter the process name to process: ")
	var targetProcess string
	fmt.Scanln(&targetProcess)

	// Call the function to assign and print the values with combined labels for the specified process
	assignRandomStagesOperationsActionsMeasures(processes, stages, operations, actions, measures, targetProcess)
}
