package main

// import (
// 	"log"

// 	"github.com/northwesternmutual/grammes"
// )

// func main() {
// 	// Creates a new client with the localhost IP.
// 	client, err := grammes.DialWithWebSocket("ws://172.31.64.105:8182")
// 	if err != nil {
// 		log.Fatalf("Error while creating client: %s\n", err.Error())
// 	}

// 	_ = client
// }

import (
	"errors"
	"fmt"
	"strings"

	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
)

func main() {

	// Creating the connection to the server.
	driverRemoteConnection, err := gremlingo.NewDriverRemoteConnection("wss://neptunedbinstance-6sb274y70mqt.c7a08meoiuv1.us-east-2.neptune.amazonaws.com:8182/gremlin",
		func(settings *gremlingo.DriverRemoteConnectionSettings) {
			settings.TraversalSource = "g"
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	// Cleanup; Deffered Cleanup
	defer driverRemoteConnection.Close()

	// Creating graph traversal, this traverses the graph and initializes a traversal object
	//which is used to navigate & query the graph data stored remotely
	g := gremlingo.Traversal_().WithRemote(driverRemoteConnection)

	//MODIFYING TO USE API DIRECTLY FROM THE CODE AS OPPOSED TO HTTP

	//Goal 1: single out only the hierarchy (process) from the output
	//note output only contains one hierarchy/process
	name := "test"
	output, err := ModelData(name)
	if err != nil {
		fmt.Println("error")
	}
	// fmt.Println(output.Hierarchy)

	//Goal 2: Store the name of the process, the stages, and any measures
	processName := output.Hierarchy.Process
	fmt.Println("Here is process name:" + processName)
	//insertProcess(g, processName)
	process_stages := output.Hierarchy.Stages

	//this outer loop traverses stages
	for _, stage := range process_stages {
		//insertStage(g, stage)
		operations := stage.Operations

		//this loop traverses operations
		for _, operation := range operations {
			//insertOperation(g,operation)
			fmt.Println(operation.Operation)
			actions := operation.Actions

			//this loop traverses operations
			for _, action := range actions {
				//fmt.Println("Element at index", i, ":", value)
				//insertAction(action)
				measures := action.Measures

				for _, measure := range measures {
					fmt.Println(measure)
					//insertMeasure(g, measure)
				}

			}

		}
	}

	result, err := g.AddE("has").From(g.V("test-1")).To(g.V("test-1-1-1-M1")).Next()

	fmt.Println("here is result:")

	fmt.Println(result)
	fmt.Println("here is err")

	fmt.Println(err)

	// Perform traversal
	results, err := g.V().Limit(5).ToList()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Print results
	fmt.Println("printing here")
	for _, r := range results {
		fmt.Println(r.GetString())
	}
}

// fucntion to insert a Process vertex within the graph
// returns an error if insertion fails, else returns nil
func insertProcess(g *gremlingo.GraphTraversalSource, processName string) error {

	_, err := g.AddV("Process").Property("name", processName).Property(gremlingo.T.Id, processName).As(processName).Next()

	if err != nil {
		return errors.New("failed to execute query for process:" + processName + "\nError:" + err.Error())
	}

	return nil
}

// fucntion to insert a Measure vertex within the graph
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
		edgeMeasure(g, stageID, measureID)
	}

	return stageID, nil
}

// fucntion to insert a Measure vertex within the graph
// returns an error if insertion fails, else returns nil
func insertOperation(g *gremlingo.GraphTraversalSource, operation Operation) (string, error) {
	operationName := operation.Operation
	parts := strings.Split(operationName, " ")
	operationID := parts[1]

	_, err := g.AddV("Stage").Property("name", operationName).Property(gremlingo.T.Id, operationID).As(operationID).Next()

	if err != nil {
		return "", errors.New("failed to execute query for operationID:" + operationID + "\nError:" + err.Error())
	}

	return operationID, nil
}

// fucntion to insert an edge between a measure and a stage
// returns an error if insertion fails, else returns nil
func edgeMeasure(g *gremlingo.GraphTraversalSource, stageID string, measureID string) error {

	//TODO ADD AN APPROPRIATE ID FOR EDGE
	_, err := g.AddE("has").From(stageID).To(measureID).Next()

	if err != nil {
		return errors.New("failed to create an edge from:" + stageID + "to" + measureID + "\nError:" + err.Error())
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
