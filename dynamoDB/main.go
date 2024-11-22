package main

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Req struct {
	Type string `json:"type"`
	Name string `json:"name,omitempty"`
	Proc string `json:"proc,omitempty"`
}

func main() {
	// Initialize Echo instance
	e := echo.New()

	// Middleware setup
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	// Route setup
	e.POST("/build", buildHandler)
	e.POST("/query", queryHandler)

	// Start server
	port := ":1010"
	fmt.Printf("Server running on http://localhost%s\n", port)
	e.Logger.Fatal(e.Start(port))
}

// buildHandler handles the /build route
func buildHandler(c echo.Context) error {
	req := new(Req)

	// Bind JSON request to struct
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate request type
	if req.Type != "build" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid type for build route")
	}

	// Call your build function
	status, timeTaken, errMsg := PopulateDatabase(req.Name)

	// Construct the response
	resp := map[string]string{
		"status": status,
		"time":   timeTaken,
		"err":    errMsg,
	}

	// Send response as JSON
	return c.JSON(http.StatusOK, resp)
}

// queryHandler handles the /query route
func queryHandler(c echo.Context) error {
	req := new(Req)

	// Bind JSON request to struct
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Declare variables for response
	var status, timeTaken, results, errMsg string

	// Determine the type of query
	switch req.Type {
	case "measures":
		status, timeTaken, results, errMsg = query(req.Proc, "Measure")
	case "results":
		status, timeTaken, results, errMsg = query(req.Proc, "Result")
	case "materials":
		status, timeTaken, results, errMsg = query(req.Proc, "RawMaterial")
	default:
		errMsg = "Invalid query type"
		status = "failed"
		timeTaken = "0s"
		results = ""
	}

	// Construct the response
	resp := map[string]string{
		"status": 	status,
		"time":   	timeTaken,
		"results": 	results, 
		"err":    	errMsg,
	}

	// Send response as JSON
	return c.JSON(http.StatusOK, resp)
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
