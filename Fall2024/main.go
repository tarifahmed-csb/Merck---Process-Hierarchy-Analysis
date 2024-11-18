package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Req struct {
	Type string `json:"type"`
	Name string `json:"name,omitempty"`
	Proc string `json:"proc,omitempty"` 
}

func main() {
	//Setting up HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/build", buildHandler)
	mux.HandleFunc("/query", queryHandler)

	// Start server
	port := ":1011" 
		// Neptune : port 1010
		// Dynamo : port 1011
		// PostgreSQL : port 1012
	fmt.Printf("Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}

//Building the table 
func buildHandler(response http.ResponseWriter, request *http.Request) {
	var req Req

	//Decode the request 
	err := json.NewDecoder(request.Body).Decode(&req)

	//Error
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	//Error
	if req.Type != "build" {
		http.Error(response, "Invalid type for build route", http.StatusBadRequest)
		return
	}

	//Call your build function
	status, timeTaken, errMsg := PopulateDatabase(req.Name) 

	// Construct the response
	resp := map[string]string{
		"status": status,
		"time":   timeTaken,
		"err":    errMsg,
	}

	// Encode and send the response
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(resp)
}

func queryHandler(response http.ResponseWriter, request *http.Request) {
	var req Req

	// Decode the request body
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	//Determine the type of query 
	//Call query function 
	var entityType string
	switch req.Type {
	case "measures":
		status, timeTaken, errMsg := query(req.Proc, "Measure")
	case "results":
		status, timeTaken, errMsg := query(req.Pro, "Result")
	case "materials":
		status, timeTaken, errMsg := query(req.Pro, "RawMaterial")
	default:
		entityType = "Process" // Default to "Process" if none of the above
	}

	// Construct the response
	resp := map[string]string{
		"status": status,
		"time":   timeTaken,
		"err":    errMsg,
	}

	// Encode and send the response
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(resp)
}

// func main() {
// 	fmt.Println("Merck-Fall2024-Final")
// 	fmt.Println("Choose an option:")
// 	fmt.Println("1. Populate")
// 	fmt.Println("2. Delete all items")
// 	fmt.Println("3. Run Query")

// 	// User input
// 	scanner := bufio.NewScanner(os.Stdin)
// 	scanner.Scan()
// 	choice := strings.TrimSpace(scanner.Text())

// 	// Defining table
// 	tableName := "Merck-Fall2024-Final"

// 	// Running the selected option
// 	switch choice {
// 	case "1":
// 		PopulateDatabase(tableName)
// 	case "2":
// 		DeleteAllItemsFromTable(tableName)
// 	case "3":
// 		fmt.Println("Select query type:")
// 		fmt.Println("1. Get Measures")
// 		fmt.Println("2. Get Results")
// 		fmt.Println("3. Get Raw Materials")

// 		scanner.Scan()
// 		queryChoice := strings.TrimSpace(scanner.Text())

// 		switch queryChoice {
// 		case "1":
// 			getMeasures(tableName)
// 		case "2":
// 			getResults(tableName)
// 		case "3":
// 			getRawMaterials(tableName)
// 		default:
// 			fmt.Println("Invalid query option.")
// 		}
// 	default:
// 		fmt.Println("Invalid option. Please choose a valid number.")
// 	}
// }
