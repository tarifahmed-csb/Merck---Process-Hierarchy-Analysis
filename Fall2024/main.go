package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Merck-Fall2024")
	fmt.Println("Choose an option:")
	fmt.Println("1. Populate")
	fmt.Println("2. Delete all items")
	fmt.Println("3. Get all stages")
	fmt.Println("4. Get all operations")
	fmt.Println("5. Run Parent Children Query")

	// User input
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := scanner.Text()

	// Defining table
	tableName := "Merck-Fall2024"

	// Running the option
	switch strings.TrimSpace(choice) {
	case "1":
		PopulateDatabase(tableName)
	case "2":
		DeleteAllItemsFromTable(tableName)
	case "3":
		GetAllStages(tableName)
	case "4":
		GetAllOperations(tableName)
	case "5":
		err := accessPattern(tableName)
		if err != nil {
			fmt.Println("Error:", err)
		}
	default:
		fmt.Println("Invalid option. Please choose a valid number.")
	}
}
