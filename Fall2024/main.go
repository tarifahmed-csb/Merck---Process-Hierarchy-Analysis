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
	fmt.Println("1. Populate the database")
	fmt.Println("2. Delete all items from the database")
	fmt.Println("3. Get all stages from the database")

	//User input
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := scanner.Text()

	//Definiting table
	tableName := "Merck-Fall2024"

	//Running the option
	switch strings.TrimSpace(choice) {
	case "1":
		PopulateDatabase(tableName)
	case "2":
		DeleteAllItemsFromTable(tableName)
	case "3":
		GetAllStages(tableName)
	}
}
