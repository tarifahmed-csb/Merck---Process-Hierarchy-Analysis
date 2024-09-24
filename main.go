package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
)

// "ws://localhost:8182/gremlin" to connect to local gremlin server; ws is websockets used for TCP connections
const database_url = "ws://localhost:8182/gremlin"

func main() {

	// Creating the connection to the server.
	driverRemoteConnection, err := gremlingo.NewDriverRemoteConnection(database_url,
		func(settings *gremlingo.DriverRemoteConnectionSettings) {
			settings.TraversalSource = "g"
		})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Deffered Cleanup; Will be called when this function (main) reaches the end
	defer driverRemoteConnection.Close()

	// Creating graph traversal, this traverses the graph and initializes a traversal object
	//which is used to navigate & query the graph data stored remotely
	g := gremlingo.Traversal_().WithRemote(driverRemoteConnection)

	//initialize a reader (takes in user input)
	reader := bufio.NewReader(os.Stdin)

	//start for demo program
	for {
		printMenu()

		input, _ := reader.ReadString('\n') //reads string until \n
		// input = strings.TrimSpace(input)
		input = input[:len(input)-2] // Remove newline character & \r char

		choice, err := strconv.Atoi(input)

		if err != nil || choice < 1 || choice > 6 {
			fmt.Println("Invalid input. Please enter a number between 1 and 5.")
			continue
		}

		switch choice {
		case 1:
			fmt.Println("You chose Option 1")
			// Ask the user to input a name for a process
			fmt.Print("Enter a name for the process: ")
			processName, _ := reader.ReadString('\n')
			processName = processName[:len(processName)-1] // Remove newline character
			insertNewProcess(g, processName)
		case 2:
			fmt.Println("You chose Option 2")
			getAllProcesses(g)

		case 3:
			fmt.Println("You chose Option 3")
			fmt.Println("Enter the name of the Process")
			processName, _ := reader.ReadString('\n')
			processName = processName[:len(processName)-1] // Remove newline character
			getAllStages(g, processName)
		case 4:
			fmt.Println("You chose Option 4")
			fmt.Println("Enter the name of the Process")
			processName, _ := reader.ReadString('\n')
			processName = processName[:len(processName)-1] // Remove newline character
			getAllChildren(g, processName)
		case 5:
			fmt.Println("You chose Option 5")
			fmt.Println("Enter the name of the Process")
			processName, _ := reader.ReadString('\n')
			processName = processName[:len(processName)-1] // Remove newline character
			getAllMeasures(g, processName)
		case 6:
			fmt.Println("You chose Option 6")
			fmt.Println("Enter the name of the Process")

			processName, _ := reader.ReadString('\n')
			processName = processName[:len(processName)-1] // Remove newline character
			getAllResults(g, processName)
		}

		// Ask the user if they want to continue
		fmt.Print("Do you want to enter another menu option? (y/n): ")
		again, _ := reader.ReadString('\n')
		again = again[:len(again)-2] // Remove newline character

		if again != "y" {
			break
		}
	}

	// rawmats := output.RawMaterials

	// for _, rawmat := range rawmats {
	// 	err = insertRawMat(g, rawmat)
	// }

	//from a process give all measures,
	//linking raw materials to processes (jumping from one process to another)
	//

	// result, err := g.AddE("has").From(g.V("test-1")).To(g.V("test-1-1-1-M1")).Next()

	// fmt.Println("here is result:")

	// fmt.Println(result)
	// fmt.Println("here is err")

	// fmt.Println(err)

	// Perform traversal
	// results, err := g.V().Limit(5).ToList()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // Print results
	// fmt.Println("printing here")
	// for _, r := range results {
	// 	fmt.Println(r.GetString())
	// }
}

func printMenu() {
	fmt.Println("Choose an action (1-5):")
	fmt.Println("1. Insert a new process")
	fmt.Println("2. Query all Processes")
	fmt.Println("3. Query all Stages of a given Process")
	fmt.Println("4. Query all children of a given Process")
	fmt.Println("5. Query all measures for a given Process")
	fmt.Println("6. Query all results for a given Process")
	fmt.Print("Enter your choice: ")

}

func insertNewProcess(g *gremlingo.GraphTraversalSource, name string) {
	fmt.Println("Inserting new Process: " + name + "\nLoading data...")

	//creating data
	output, err := ModelData(name)
	if err != nil {
		log.Fatal(err)
	}

	//Store process, the stages, operations, & actions and any measures into variables
	start := time.Now()

	//Here we store process
	processName, err := insertProcess(g, output.Hierarchy.Process)
	if err != nil {
		log.Fatal(err)
	}

	//Here we iterate through processes to store the next levels in the hierarchy
	process_stages := output.Hierarchy.Stages

	//this outer loop traverses stages
	for _, stage := range process_stages {
		//this inserts the stage along with any connected measures
		stageID, err := insertStage(g, stage)
		if err != nil {
			log.Fatal(err)
		}
		//this connects newly inserted stage with the process
		err = edgeProcessStage(g, processName, stageID)
		if err != nil {
			log.Fatal(err)
		}
		operations := stage.Operations

		//this loop traverses operations
		for _, operation := range operations {
			//this inserts the operation along with any connected measures
			operationID, err := insertOperation(g, operation)
			if err != nil {
				log.Fatal(err)
			}
			//this connects newly inserted operation with the stage
			err = edgeStageOperation(g, stageID, operationID)
			if err != nil {
				log.Fatal(err)
			}

			actions := operation.Actions

			//this loop traverses actions
			for _, action := range actions {
				//this inserts the action along with any connected measures
				actionID, err := insertAction(g, action)
				if err != nil {
					log.Fatal(err)
				}
				//this connects newly inserted action with operation
				err = edgeOperationAction(g, operationID, actionID)
				if err != nil {
					log.Fatal(err)
				}
				measures := action.Measures

				//this loop traverses measures
				for _, measure := range measures {

					measureID, err := insertMeasure(g, measure)
					if err != nil {
						log.Fatal(err)
					}
					err = edgeActionMeasure(g, actionID, measureID)
					if err != nil {
						log.Fatal(err)
					}
				}

			}
		}
	}
	//Now iterate through x_paths & adds the vertexes
	//Plus corresponding edges from measure to xpath
	x_paths := output.Xpath

	for _, x_path := range x_paths {
		err := insertXPathnEdge(g, x_path)
		if err != nil {
			log.Fatal(err)
		}

	}

	ress := output.Results

	for _, result := range ress {
		err = insertResult(g, result)
		//fmt.Println("in here")
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(result.Result)
		// fmt.Println(result.ResultName)
		//RESULT NAME HOLDS MEAURE ID

		//insertXPath(g, x_path)
	}
	end := time.Now()

	fmt.Println("Database Updated!")
	elapsedTime := end.Sub(start)
	fmt.Println("Elapsed Time:", elapsedTime)
}

func getAllResults(g *gremlingo.GraphTraversalSource, name string) {
	fmt.Println("\n\nGetting all results for Process " + name + ":")
	start := time.Now()

	result, err := g.V(name).Repeat(gremlingo.T__.Out()).Until(gremlingo.T__.HasLabel("result")).ToList()

	end := time.Now()
	if err != nil {
		log.Fatal(err)
	}
	for _, vertex := range result {
		fmt.Println(vertex)
	}
	elapsedTime := end.Sub(start)
	fmt.Println("Elapsed Time:", elapsedTime)
}

func getAllProcesses(g *gremlingo.GraphTraversalSource) {
	fmt.Println("\n\nHere are all Processes:")
	start := time.Now()
	result, err := g.V().HasLabel("Process").ToList()
	end := time.Now()
	if err != nil {
		log.Fatal(err)
	}

	for _, vertex := range result {
		fmt.Println(vertex)
	}
	elapsedTime := end.Sub(start)
	fmt.Println("Elapsed Time:", elapsedTime)
}

func getAllStages(g *gremlingo.GraphTraversalSource, name string) {
	start := time.Now()
	result, err := g.V(name).Out().ToList()
	end := time.Now()
	if err != nil {
		log.Fatal(err)
	}

	for _, vertex := range result {
		fmt.Println(vertex)
	}
	elapsedTime := end.Sub(start)
	fmt.Println("Elapsed Time:", elapsedTime)

}

func getAllMeasures(g *gremlingo.GraphTraversalSource, name string) {
	fmt.Println("\n\nGetting all measures for Process " + name + ":")
	start := time.Now()

	result, err := g.V(name).Repeat(gremlingo.T__.Out()).Until(gremlingo.T__.HasLabel("Measure")).ToList()

	end := time.Now()
	if err != nil {
		log.Fatal(err)
	}
	for _, vertex := range result {
		fmt.Println(vertex)
	}
	elapsedTime := end.Sub(start)
	fmt.Println("Elapsed Time:", elapsedTime)
}

func getAllChildren(g *gremlingo.GraphTraversalSource, name string) {
	fmt.Println("\n\nGetting all children for Process " + name + ":")
	start := time.Now()

	result, err := g.V(name).Repeat(gremlingo.T__.Out()).Until(gremlingo.T__.HasLabel("Measure")).Path().ToList()

	end := time.Now()
	if err != nil {
		log.Fatal(err)
	}
	for _, vertex := range result {
		fmt.Println(vertex)
	}
	elapsedTime := end.Sub(start)
	fmt.Println("Elapsed Time:", elapsedTime)
}

// fucntion to insert a result vertex within the graph
// returns an error if insertion fails, else returns nil
// func insertRawMat(g *gremlingo.GraphTraversalSource, rawMat RawMaterials) error {
// 	x_path_name := result.Xpath
// 	measureID := result.ResultName
// 	resultID := result.BatchID

// 	res, err := g.AddV("result").Property("name", resultID).Property("materialNum", result.MaterialNum).Property("batchID", resultID).Property("xpath", x_path_name).Property("result", result.Result).Property("DOM", result.DOM).Property("site", result.Site).Property("measureID", resultID).Next()
// 	if err != nil {
// 		return errors.New("failed to execute query for xpath:" + x_path_name + "\nError:" + err.Error())
// 	}

// 	vertex, err := res.GetVertex()
// 	if err != nil {
// 		return errors.New("Not a vertex" + err.Error())
// 	}

// 	//result vertex inserted, now adding edge from respective measureID to Xpath vertex
// 	//newID := measureID + x_path_name
// 	_, err = g.AddE("links").From(g.V(x_path_name)).To(g.V(vertex.Id)).Next()
// 	fmt.Println("here is xpath")
// 	fmt.Println(x_path_name)
// 	fmt.Println("here is vertex id")
// 	fmt.Println(resultID)

// 	if err != nil {
// 		return errors.New("failed to insert edge for result with measure:" + measureID + "\nError:" + err.Error())
// 	}

// 	return nil
// }

// fucntion to insert a result vertex within the graph
// returns an error if insertion fails, else returns nil
func insertResult(g *gremlingo.GraphTraversalSource, result Results) error {
	x_path_name := result.Xpath
	measureID := result.ResultName
	resultID := result.BatchID

	res, err := g.AddV("result").Property("name", resultID).Property("materialNum", result.MaterialNum).Property("batchID", resultID).Property("xpath", x_path_name).Property("result", result.Result).Property("DOM", result.DOM).Property("site", result.Site).Property("measureID", resultID).Next()
	if err != nil {
		return errors.New("failed to execute query for xpath:" + x_path_name + "\nError:" + err.Error())
	}

	vertex, err := res.GetVertex()
	if err != nil {
		return errors.New("Not a vertex" + err.Error())
	}

	//result vertex inserted, now adding edge from respective measureID to Xpath vertex
	//newID := measureID + x_path_name
	_, err = g.AddE("links").From(g.V(x_path_name)).To(g.V(vertex.Id)).Next()
	//fmt.Println("here is xpath")
	//fmt.Println(x_path_name)
	//fmt.Println("here is vertex id")
	//fmt.Println(resultID)

	if err != nil {
		return errors.New("failed to insert edge for result with measure:" + measureID + "\nError:" + err.Error())
	}

	return nil
}

// fucntion to insert a XPATH vertex within the graph
// returns an error if insertion fails, else returns nil
func insertXPathnEdge(g *gremlingo.GraphTraversalSource, x_path Xpath) error {
	x_path_name := x_path.Xpath
	measureID := x_path.MeasureID
	_, err := g.AddV("xpath").Property("name", x_path_name).Property("site", x_path.Site).Property("measureID", x_path.MeasureID).Property(gremlingo.T.Id, x_path.Xpath).As(x_path.Xpath).Next()

	if err != nil {
		return errors.New("failed to execute query for xpath:" + x_path_name + "\nError:" + err.Error())
	}
	//Xpath vertex inserted, now adding edge from respective measureID to Xpath vertex
	newID := measureID + x_path_name
	_, err = g.AddE("links").From(g.V(measureID)).To(g.V(x_path_name)).Property(gremlingo.T.Id, newID).Next()
	if err != nil {
		return errors.New("failed to insert edge" + newID + "\nError:" + err.Error())
	}

	return nil
}

// fucntion to insert a Process vertex within the graph.
// returns an error if insertion fails. else returns name of process inserted
func insertProcess(g *gremlingo.GraphTraversalSource, processName string) (string, error) {

	_, err := g.AddV("Process").Property("name", processName).Property(gremlingo.T.Id, processName).As(processName).Next()

	if err != nil {
		return "", errors.New("failed to execute query for process:" + processName + "\nError:" + err.Error())
	}

	return processName, nil
}

// fucntion to insert an edge between a stage and process
// returns an error if insertion fails, else returns nil
func edgeProcessStage(g *gremlingo.GraphTraversalSource, processID string, stageID string) error {

	newID := processID + stageID
	_, err := g.AddE("has").From(g.V(processID)).To(g.V(stageID)).Property(gremlingo.T.Id, newID).Next()

	if err != nil {
		return errors.New("failed to create an edge from:" + processID + "to" + stageID + "\nError:" + err.Error())
	}

	return nil
}

// fucntion to insert a Stage vertex within the graph
// returns an error if insertion fails, else returns nil
func insertStage(g *gremlingo.GraphTraversalSource, stage Stage) (string, error) {
	stageName := stage.Stage
	parts := strings.Split(stageName, " ")
	stageID := parts[1]

	_, err := g.AddV("Stage").Property("name", stageName).Property(gremlingo.T.Id, stageID).As(stageID).Next()

	if err != nil {
		return "", errors.New("failed to execute query for stageID:" + stageID + "\nError:" + err.Error())
	}

	//inserting any measures of stage
	stageMeasures := stage.Measures
	for _, measure := range stageMeasures {
		measureID, err := insertMeasure(g, measure)
		if err != nil {
			return "", err
		}
		err = edgeStageMeasure(g, stageID, measureID)
		if err != nil {
			return "", err
		}
	}

	return stageID, nil
}

// fucntion to insert an edge between a stage and operation
// returns an error if insertion fails, else returns nil
func edgeStageOperation(g *gremlingo.GraphTraversalSource, stageID string, operationID string) error {

	newID := stageID + operationID
	_, err := g.AddE("has").From(g.V(stageID)).To(g.V(operationID)).Property(gremlingo.T.Id, newID).Next()

	if err != nil {
		return errors.New("failed to create an edge from:" + stageID + "to" + operationID + "\nError:" + err.Error())
	}

	return nil
}

// fucntion to insert a Measure vertex within the graph
// returns an error if insertion fails, else returns nil
func insertOperation(g *gremlingo.GraphTraversalSource, operation Operation) (string, error) {
	operationName := operation.Operation
	parts := strings.Split(operationName, " ")
	operationID := parts[1]

	_, err := g.AddV("Operation").Property("name", operationName).Property(gremlingo.T.Id, operationID).As(operationID).Next()

	if err != nil {
		return "", errors.New("failed to execute query for operationID:" + operationID + "\nError:" + err.Error())
	}

	//inserting any measures of stage
	operationMeasures := operation.Measures
	for _, measure := range operationMeasures {
		measureID, err := insertMeasure(g, measure)
		if err != nil {
			return "", err
		}
		err = edgeOperationMeasure(g, operationID, measureID)
		if err != nil {
			return "", err
		}
	}

	return operationID, nil
}

// fucntion to insert an edge between a operation and action
// returns an error if insertion fails, else returns nil
func edgeOperationAction(g *gremlingo.GraphTraversalSource, operationID string, actionID string) error {

	newID := operationID + actionID
	_, err := g.AddE("has").From(g.V(operationID)).To(g.V(actionID)).Property(gremlingo.T.Id, newID).Next()

	if err != nil {
		return errors.New("failed to create an edge from:" + operationID + "to" + actionID + "\nError:" + err.Error())
	}

	return nil
}

func insertAction(g *gremlingo.GraphTraversalSource, action Action) (string, error) {
	actionName := action.Action
	parts := strings.Split(actionName, " ")
	actionID := parts[1]

	_, err := g.AddV("Action").Property("name", actionName).Property(gremlingo.T.Id, actionID).As(actionID).Next()

	if err != nil {
		return "", errors.New("failed to execute query for stageID:" + actionID + "\nError:" + err.Error())
	}

	return actionID, nil
}

// fucntion to insert an edge between a stage and a measure
// returns an error if insertion fails, else returns nil
func edgeStageMeasure(g *gremlingo.GraphTraversalSource, stageID string, measureID string) error {

	newID := stageID + measureID
	_, err := g.AddE("has").From(g.V(stageID)).To(g.V(measureID)).Property(gremlingo.T.Id, newID).Next()

	if err != nil {
		return errors.New("failed to create an edge from:" + stageID + "to" + measureID + "\nError:" + err.Error())
	}

	return nil
}

// fucntion to insert an edge between a measure and a stage
// returns an error if insertion fails, else returns nil
func edgeOperationMeasure(g *gremlingo.GraphTraversalSource, operationID string, measureID string) error {

	//TODO ADD AN APPROPRIATE ID FOR EDGE
	newID := operationID + measureID
	_, err := g.AddE("has").From(g.V(operationID)).To(g.V(measureID)).Property(gremlingo.T.Id, newID).Next()

	if err != nil {
		return errors.New("failed to create an edge from:" + operationID + "to" + measureID + "\nError:" + err.Error())
	}

	return nil
}

// fucntion to insert an edge between a stage and a measure
// returns an error if insertion fails, else returns nil
func edgeActionMeasure(g *gremlingo.GraphTraversalSource, actionID string, measureID string) error {

	newID := actionID + measureID
	_, err := g.AddE("has").From(g.V(actionID)).To(g.V(measureID)).Property(gremlingo.T.Id, newID).Next()

	if err != nil {
		return errors.New("failed to create an edge from:" + actionID + "to" + measureID + "\nError:" + err.Error())
	}

	return nil
}

// fucntion to insert a Measure vertex within the graph
// returns an error if insertion fails, else returns nil
func insertMeasure(g *gremlingo.GraphTraversalSource, measure Measure) (string, error) {
	measureName := measure.Measure
	measureID := measure.MeasureID

	_, err := g.AddV("Measure").Property("name", measureName).Property(gremlingo.T.Id, measureID).As(measureID).Next()

	if err != nil {
		return "", errors.New("failed to execute query for measureID:" + measureID + "\nError:" + err.Error())
	}

	return measureID, nil
}
