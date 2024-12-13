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
	// resp := map[string]string{
	// 	"status":  status,
	// 	"time":    timeTaken,
	// 	"results": results,
	// 	"err":     errMsg,
	// }
	var Resp struct {
		Status  string `json:"status"`
		Time    string `json:"time"`
		Results string `json:"results"`
		Error   string `json:"error"`
	}
	Resp.Status = status
	Resp.Time = timeTaken
	Resp.Results = results
	Resp.Error = errMsg

	// Send response as JSON
	fmt.Print(Resp)
	return c.JSON(http.StatusOK, Resp)
	//return c.JSON(http.StatusOK, resp)
}

// package main

// import (
// 	//"bufio"
// 	"fmt"
// 	//"os"
// 	//"strconv"
// 	// "strings"
// 	// "strings"
// 	"net/http"
// 	"github.com/labstack/echo/v4"
// 	"github.com/labstack/echo/v4/middleware"
// )

// type Req struct {
// 	Type string `json:"type"`
// 	Name string `json:"name,omitempty"`
// 	Proc string `json:"proc,omitempty"`
// }

// func main() {
// 	// Initialize Echo instance
// 	e := echo.New()

// 	// Middleware setup
// 	e.Use(middleware.Logger())
// 	e.Use(middleware.Recover())
// 	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
// 		AllowOrigins: []string{"*"},
// 		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
// 		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
// 	}))

// 	// Route setup
// 	e.POST("/build", buildHandler)
// 	e.POST("/query", queryHandler)

// 	// Start server
// 	port := ":1010"
// 	fmt.Printf("Server running on http://localhost%s\n", port)
// 	e.Logger.Fatal(e.Start(port))
// }

// // buildHandler handles the /build route
// func buildHandler(c echo.Context) error {
// 	req := new(Req)

// 	// Bind JSON request to struct
// 	if err := c.Bind(req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 	}

// 	// Validate request type
// 	if req.Type != "build" {
// 		return echo.NewHTTPError(http.StatusBadRequest, "Invalid type for build route")
// 	}

// 	// Call your build function
// 	status, timeTaken, errMsg := PopulateDatabase(req.Name)

// 	// Construct the response
// 	resp := map[string]string{
// 		"status": status,
// 		"time":   timeTaken,
// 		"err":    errMsg,
// 	}

// 	// Send response as JSON
// 	return c.JSON(http.StatusOK, resp)
// }

// // queryHandler handles the /query route
// func queryHandler(c echo.Context) error {
// 	req := new(Req)

// 	// Bind JSON request to struct
// 	if err := c.Bind(req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 	}

// 	// Declare variables for response
// 	var status, timeTaken, results, errMsg string

// 	// Determine the type of query
// 	switch req.Type {
// 	case "measures":
// 		status, timeTaken, results, errMsg = query(req.Proc, "Measure")
// 	case "process":
// 		status, timeTaken, results, errMsg = query(req.Proc, "Process")
// 	case "materials":
// 		status, timeTaken, results, errMsg = query(req.Proc, "RawMaterial")
// 	default:
// 		errMsg = "Invalid query type"
// 		status = "failed"
// 		timeTaken = "0s"
// 		results = ""
// 	}

// 	// Construct the response
// 	resp := map[string]string{
// 		"status": 	status,
// 		"time":   	timeTaken,
// 		"results": 	results,
// 		"err":    	errMsg,
// 	}

// 	// Send response as JSON
// 	return c.JSON(http.StatusOK, resp)
// }

// func main() {
// 	reader := bufio.NewReader(os.Stdin)

// 	for {
// 		fmt.Println("Merck-Fall2024-Final")
// 		fmt.Println("Choose an option:")
// 		fmt.Println("1. Populate")
// 		fmt.Println("2. Delete all items")
// 		fmt.Println("3. Run Query")
// 		fmt.Println("4. Populate Mass Data")

// 		// User input

// 		input, _ := reader.ReadString('\n') //reads string until \n
// 		input = input[:len(input)-1]        // Remove newline character & \r char

// 		choice, err := strconv.Atoi(input)

// 		if err != nil || choice < 1 || choice > 6 {
// 			fmt.Println("Invalid input. Please enter a number between 1 and 5.")

// 		}
// 		// Defining table
// 		// tableName := "Merck-Fall2024-Final-2"

// 		// Running the selected option
// 		switch choice {
// 		case 1:
// 			fmt.Println("You chose Option 1")
// 			// Ask the user to input a name for a process
// 			fmt.Print("Enter a name for the process: ")
// 			processName, _ := reader.ReadString('\n')
// 			processName = processName[:len(processName)-1] // Remove newline character & /r
// 			PopulateDatabase(processName)

// 		case 4:
// 			fmt.Println("You chose Option 4")
// 			// Ask the user to input a name for a process
// 			fmt.Print("Enter a name for the starting process name: ")
// 			processName, _ := reader.ReadString('\n')
// 			processName = processName[:len(processName)-1] // Remove newline character & /r
// 			for i := 0; i < 10; i++ {
// 				str := strconv.Itoa(i)
// 				PopulateDatabase(processName + str)
// 			}

// 		case 3:
// 			fmt.Println("You chose Option 3")

// 			// Prompt for process name
// 			fmt.Print("Enter a prefix for the process name: ")
// 			prefix, _ := reader.ReadString('\n')
// 			prefix = strings.TrimSpace(prefix)

// 			// Prompt for entity type with numbered options
// 			fmt.Println("Select entity type:")
// 			fmt.Println("1. Measures")
// 			fmt.Println("2. Results")
// 			fmt.Println("3. Raw Materials")
// 			fmt.Print("Enter the number corresponding to the entity type: ")
// 			entityTypeChoice, _ := reader.ReadString('\n')
// 			entityTypeChoice = strings.TrimSpace(entityTypeChoice)

// 			var entityType string
// 			switch entityTypeChoice {
// 			case "1":
// 				entityType = "measures"
// 			case "2":
// 				entityType = "results"
// 			case "3":
// 				entityType = "raw materials"
// 			default:
// 				fmt.Println("Invalid entity type choice. Please try again.")
// 				continue
// 			}

// 			// Call the query function
// 			status, duration, results, errMsg := query(prefix, entityType)

// 			// Display the results
// 			if status == "Failed" {
// 				fmt.Printf("Query failed. Error: %s\n", errMsg)
// 			} else {
// 				fmt.Printf("Query completed in %s. Results:\n%s\n", duration, results)
// 			}
// 		default:
// 			fmt.Println("Invalid option. Please choose a valid number.")
// 		}
// 		// Ask the user if they want to continue
// 		fmt.Print("Do you want to enter another menu option? (y/n): ")
// 		again, _ := reader.ReadString('\n')
// 		again = again[:len(again)-1] // Remove newline character

// 		if again != "y" {
// 			break
// 		}

// 	}
// }
