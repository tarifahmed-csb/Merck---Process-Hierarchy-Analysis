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
	"log"
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
	name := "bye"
	output, err := ModelData(name)
	if err != nil {
		log.Fatal(err)
	}

	//Goal 2: Store process, the stages, operations, & actions and any measures into variables

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

	ressy := output.Results

	for _, result := range ressy {
		err = insertResult(g, result)
		fmt.Println("in here")
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(result.Result)
		// fmt.Println(result.ResultName)
		//RESULT NAME HOLDS MEAURE ID

		//insertXPath(g, x_path)
	}

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
	fmt.Println("here is xpath")
	fmt.Println(x_path_name)
	fmt.Println("here is vertex id")
	fmt.Println(resultID)

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

// fucntion to insert a Process vertex within the graph
// returns an error if insertion fails, else returns nil
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
