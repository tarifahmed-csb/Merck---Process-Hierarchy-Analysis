package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
)

//This file contains all of the necessary functions to interact with graph DB

// Function to insert a Raw Material vertex within the graph.
// returns an error if insertion fails, else returns nil
func insertRawMat(g *gremlingo.GraphTraversalSource, rawMat RawMaterials, processName string) error {
	childBatchID := rawMat.ChildBatchID
	childMaterialName := rawMat.ChildMaterialName
	//there will be various childMaterialNum as inputs for the given ProcessName
	childMaterialNum := rawMat.ChildMaterialNum
	//the batch of the output
	parentBatchID := rawMat.ParentBatchID
	//remember parentMaterialNum is the output for the given processName
	parentMaterialNum := rawMat.ParentMaterialNum

	//adding a new vertex to g
	res, err := g.AddV("rawMat").Property("name", parentMaterialNum).Property("materialNum", parentMaterialNum).Property("batchID", parentBatchID).Property("inputBatchID", childBatchID).Property("inputMaterialNum", childMaterialNum).Next()
	if err != nil {
		return errors.New("failed to insert raw material:" + childMaterialName + "\nBatch ID:" + childBatchID + "\nError:" + err.Error())
	}

	//ensuring res is a valid vertex
	vertex, err := res.GetVertex()
	if err != nil {
		return errors.New("Not a vertex" + err.Error())
	}

	//result vertex inserted, now adding edge from respective process to raw material it outputs
	//newID := measureID + x_path_name
	_, err = g.AddE("outputs").From(g.V(processName)).To(g.V(vertex.Id)).Next()

	if err != nil {
		return errors.New("failed to insert edge for process with ouput:" + parentMaterialNum + "\nWith batch:" + parentBatchID + " & input batch:" + childBatchID + "\nError:" + err.Error())
	}

	return nil
}

func insertResult(g *gremlingo.GraphTraversalSource, result Results) error {
	x_path_name := result.Xpath
	measureID := result.ResultName
	resultID := result.BatchID

	// TODO: change resultID to include xpathhname + batchID
	//currently resultID is being stored in vertex.ID, so it is based on gremlin
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
	_, err = g.AddE("xpath_out").From(g.V(x_path_name)).To(g.V(vertex.Id)).Next()

	if err != nil {
		return errors.New("failed to insert edge for result with measure:" + measureID + "\nError:" + err.Error())
	}

	return nil
}

// fucntion to insert a XPATH vertex & corresponding edge within the graph
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
	_, err = g.AddE("xpath_in").From(g.V(measureID)).To(g.V(x_path_name)).Property(gremlingo.T.Id, newID).Next()
	if err != nil {
		return errors.New("failed to insert edge" + newID + "\nError:" + err.Error())
	}

	return nil
}

// fucntion to insert a Process vertex within the graph.
// returns an error if insertion fails. else returns name of process inserted
func insertProcess(g *gremlingo.GraphTraversalSource, processName string, inputs []string, output string) (string, error) {

	v := g.AddV("process").Property("name", processName).Property("output", output).Property(gremlingo.T.Id, processName).As(processName)

	// if err != nil {
	// 	return "", errors.New("failed to execute query for process:" + processName + "\nError:" + err.Error())
	// }

	//vertex, err := res.GetVertex()
	// if err != nil {
	// 	return "", errors.New("Not a vertex" + err.Error())
	// }

	for _, input := range inputs {
		v = v.Property("inputs", input)
	}
	//check for errors once executing/finalizing traversal with Next()
	_, err := v.Next()

	if err != nil {
		return "", errors.New("failed to insert process:" + processName + "\nError:" + err.Error())
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

	_, err := g.AddV("stage").Property("name", stageName).Property(gremlingo.T.Id, stageID).As(stageID).Next()

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

	_, err := g.AddV("operation").Property("name", operationName).Property(gremlingo.T.Id, operationID).As(operationID).Next()

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

	_, err := g.AddV("action").Property("name", actionName).Property(gremlingo.T.Id, actionID).As(actionID).Next()

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

	_, err := g.AddV("measure").Property("name", measureName).Property(gremlingo.T.Id, measureID).As(measureID).Next()

	if err != nil {
		return "", errors.New("failed to execute query for measureID:" + measureID + "\nError:" + err.Error())
	}

	return measureID, nil
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
