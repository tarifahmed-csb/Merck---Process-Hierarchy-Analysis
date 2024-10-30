package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Merck-Fall2024-Final")
	fmt.Println("Choose an option:")
	fmt.Println("1. Populate")
	fmt.Println("2. Delete all items")
	fmt.Println("3. Run Query")

	// User input
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	// Defining table
	tableName := "Merck-Fall2024-Final"

	// Running the selected option
	switch choice {
	case "1":
		PopulateDatabase(tableName)
	case "2":
		DeleteAllItemsFromTable(tableName)
	case "3":
		query(tableName)
	default:
		fmt.Println("Invalid option. Please choose a valid number.")
	}
}
