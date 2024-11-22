package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
)

//This file contains all of the necessary functions to interact with Neptune DB with a gremlingo.GraphTraversalSource

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
	res, err := g.AddV("rawMat").Property("name", parentMaterialNum).Property("materialNum", parentMaterialNum).Property("batchID", parentBatchID).Property("inputBatchID", childBatchID).Property("inputMaterialNum", childMaterialNum).Property("inputMaterialName", childMaterialName).Next()
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

func insertInputEdgeForProcess(g *gremlingo.GraphTraversalSource, newProcessId string, originalProcessId string) error {
	//we have measureId corresponding to the xpath to the result, so we can just form an edge between measure & result
	_, err := g.AddE("input").From(g.V(newProcessId)).To(g.V(originalProcessId)).Next()
	if err != nil {
		return errors.New("failed to insert input edge from process: " + newProcessId + "to process:  " + originalProcessId + "\nError:" + err.Error())
	}

	return nil

}

// func insertResult(g *gremlingo.GraphTraversalSource, result Results) error {
// 	x_path_name := result.Xpath
// 	measureID := result.ResultName
// 	resultID := result.BatchID

// 	// TODO: change resultID to include xpathhname + batchID
// 	//currently resultID is being stored in vertex.ID, so it is based on gremlin
// 	res, err := g.AddV("result").Property("name", "result:"+resultID).Property("materialNum", result.MaterialNum).Property("batchID", resultID).Property("xpath", x_path_name).Property("result", result.Result).Property("DOM", result.DOM).Property("site", result.Site).Property("measureID", resultID).Next()
// 	if err != nil {
// 		return errors.New("failed to execute query for xpath:" + x_path_name + "\nError:" + err.Error())
// 	}

// 	vertex, err := res.GetVertex()
// 	if err != nil {
// 		return errors.New("Not a vertex" + err.Error())
// 	}

// 	//result vertex inserted, now adding edge from respective measureID to Xpath vertex
// 	//newID := measureID + x_path_name
// 	_, err = g.AddE("xpath_out").From(g.V(x_path_name)).To(g.V(vertex.Id)).Next()

// 	if err != nil {
// 		return errors.New("failed to insert edge for result with measure:" + measureID + "\nError:" + err.Error())
// 	}

// 	return nil
// }

func connectMeasureToResultViaXpath(g *gremlingo.GraphTraversalSource, processId string, result Results, x_path []Xpath) error {
	//measure has a related xpath
	//xpath has an xpath from measure to result (which has the same xpath
	//so we are creating a new vertex of type results
	//then appeninding an edge from measure to results
	//note there will be many measures to results
	res_x_path := result.Xpath
	resultName := result.ResultName
	batchId := result.BatchID
	measureId := ""
	for _, x_p := range x_path {
		if x_p.Xpath == res_x_path {
			measureId = x_p.MeasureID
			break
		}
	}

	// TODO: change resultID to include xpathhname + batchID
	//currently resultID is being stored in vertex.ID, so it is based on gremlin
	res, err := g.AddV("result").Property("name", "result:"+resultName).Property("materialNum", result.MaterialNum).Property("batchID", batchId).Property("xpath", res_x_path).Property("result", result.Result).Property("DOM", result.DOM).Property("site", result.Site).Property("result name", resultName).Property("measureID", measureId).Next()
	if err != nil {
		return errors.New("failed to insert result with batchID" + batchId + "\nError:" + err.Error())
	}
	//in order to draw an edge from measure to results, we must find the correct measure.
	//This can be found by traversing the process with the given id
	//so we must have the processID, then traverse to find the measure with the corresponding measureID

	vertex, err := res.GetVertex()
	if err != nil {
		return errors.New("Not a vertex" + err.Error())
	}

	resMeasure, err := g.V(processId).Repeat(gremlingo.T__.Out("has")).Until(gremlingo.T__.Has("measureId", measureId)).Next()
	if err != nil {
		return errors.New("Cannot find measure with ID: " + measureId + " to connect with result: " + resultName + err.Error())
	}

	measureVertex, err := resMeasure.GetVertex()

	//ensuring res is a valid vertex
	if err != nil {
		return errors.New("Measure queried is not a vertex" + err.Error())
	}
	//we have measureId corresponding to the xpath to the result, so we can just form an edge between measure & result
	_, err = g.AddE("xpath").From(measureVertex).To(vertex).Next()
	if err != nil {
		return errors.New("failed to insert edge from measureID" + measureId + "to result with batchID " + batchId + "\nError:" + err.Error())
	}

	return nil

}

// fucntion to insert a XPATH vertex & corresponding edge within the graph
// returns an error if insertion fails, else returns nil
// func insertXPathnEdge(g *gremlingo.GraphTraversalSource, x_path Xpath) (int64, error) {
// 	x_path_name := x_path.Xpath
// 	measureID := x_path.MeasureID
// 	res, err := g.AddV("xpath").Property("name", x_path_name).Property("site", x_path.Site).Property("measureID", x_path.MeasureID).Next()

// 	if err != nil {
// 		return errors.New("failed to execute query for xpath:" + x_path_name + "\nError:" + err.Error())
// 	}
// 	vertex, err := res.GetVertex()
// 	if err != nil {
// 		return 0, errors.New("Not a vertex" + err.Error())
// 	}
// 	//Xpath vertex inserted, now adding edge from respective measureID to Xpath vertex
// 	_, err = g.AddE("xpath_in").From(g.V(measureID)).To(g.V(x_path_name)).Next()
// 	if err != nil {
// 		return errors.New("failed to insert edgefrom xpath to edge " + "\nError:" + err.Error())
// 	}

// 	return nil
// }

// fucntion to insert a Process vertex within the graph.
// returns an error if insertion fails. else returns name of process inserted
func insertProcessDB(g *gremlingo.GraphTraversalSource, processName string, inputs []string, output string) (string, error) {

	processId := processName + output
	v := g.AddV("process").Property("name", processName).Property("output", output).Property(gremlingo.T.Id, processId)

	// if err != nil {
	// 	return "", errors.New("failed to execute query for process:" + processName + "\nError:" + err.Error())
	// }

	//vertex, err := res.GetVertex()
	// if err != nil {
	// 	return "", errors.New("Not a vertex" + err.Error())
	// }

	//keep entering more input property attributes (there may be multiple inputs)
	for _, input := range inputs {
		v = v.Property("inputs", input)
	}
	//check for errors once executing/finalizing traversal with Next()
	_, err := v.Next()

	if err != nil {
		return "", errors.New("failed to insert process:" + processName + "\nError:" + err.Error())
	}

	// vertex, err := res.GetVertex()
	// if err != nil {
	// 	return "", errors.New("Not a vertex" + err.Error())
	// }

	//returns the process ID
	return processId, nil
}

// fucntion to insert an edge between a stage and process
// returns an error if insertion fails, else returns nil
func edgeProcessStage(g *gremlingo.GraphTraversalSource, processID string, stageID int64) error {

	_, err := g.AddE("has").From(g.V(processID)).To(g.V(stageID)).Next()

	if err != nil {
		return errors.New("failed to create an edge from process: to stage:" + "\nError:" + err.Error())
	}

	return nil
}

// fucntion to insert a Stage vertex within the graph
// returns an error if insertion fails, else returns nil
func insertStage(g *gremlingo.GraphTraversalSource, processId string, stage Stage) (int64, error) {
	stageName := stage.Stage
	parts := strings.Split(stageName, " ")
	stageID := parts[1]

	res, err := g.AddV("stage").Property("name", stageName).Next()

	if err != nil {
		return 0, errors.New("failed to execute query for stageID:" + stageID + "\nError:" + err.Error())
	}

	vertex, err := res.GetVertex()
	if err != nil {
		return 0, errors.New("Not a vertex" + err.Error())
	}
	stageVertexId := vertex.Id.(int64)

	//inserting any measures of stage
	stageMeasures := stage.Measures

	for _, measure := range stageMeasures {
		measureID, err := insertMeasure(g, stageVertexId, measure)
		if err != nil {
			return 0, err
		}
		err = edgeStageMeasure(g, stageVertexId, measureID)
		if err != nil {
			return 0, err
		}
	}
	return stageVertexId, nil
}

// fucntion to insert an edge between a stage and operation
// returns an error if insertion fails, else returns nil
func edgeStageOperation(g *gremlingo.GraphTraversalSource, stageID int64, operationID int64) error {

	_, err := g.AddE("has").From(g.V(stageID)).To(g.V(operationID)).Next()

	if err != nil {
		return errors.New("failed to create an edge from stage to operation" + "\nError:" + err.Error())
	}

	return nil
}

// fucntion to insert a Measure vertex within the graph
// returns an error if insertion fails, else returns nil
func insertOperation(g *gremlingo.GraphTraversalSource, stageID int64, operation Operation) (int64, error) {
	operationName := operation.Operation
	parts := strings.Split(operationName, " ")
	operationID := parts[1]

	res, err := g.AddV("operation").Property("name", operationName).Next()

	if err != nil {
		return 0, errors.New("failed to execute query for operationID:" + operationID + "\nError:" + err.Error())
	}

	vertex, err := res.GetVertex()
	if err != nil {
		return 0, errors.New("Not a vertex" + err.Error())
	}
	operationVertexId := vertex.Id.(int64)

	//inserting any measures of stage
	operationMeasures := operation.Measures
	for _, measure := range operationMeasures {
		measureID, err := insertMeasure(g, operationVertexId, measure)
		if err != nil {
			return 0, err
		}
		err = edgeOperationMeasure(g, operationVertexId, measureID)
		if err != nil {
			return 0, err
		}
	}

	return operationVertexId, nil
}

// fucntion to insert an edge between a operation and action
// returns an error if insertion fails, else returns nil
func edgeOperationAction(g *gremlingo.GraphTraversalSource, operationID int64, actionID int64) error {

	_, err := g.AddE("has").From(g.V(operationID)).To(g.V(actionID)).Next()

	if err != nil {
		return errors.New("failed to create an edge from operation to  action" + "\nError:" + err.Error())
	}

	return nil
}

func insertAction(g *gremlingo.GraphTraversalSource, operationID int64, action Action) (int64, error) {
	actionName := action.Action
	parts := strings.Split(actionName, " ")
	actionID := parts[1]

	res, err := g.AddV("action").Property("name", actionName).Next()

	if err != nil {
		return 0, errors.New("failed to execute query for stageID:" + actionID + "\nError:" + err.Error())
	}

	vertex, err := res.GetVertex()
	if err != nil {
		return 0, errors.New("Not a vertex" + err.Error())
	}
	actionVertexId := vertex.Id.(int64)

	measures := action.Measures

	//this loop traverses action measures
	for _, measure := range measures {

		measureID, err := insertMeasure(g, actionVertexId, measure)
		if err != nil {
			log.Fatal(err)
		}
		err = edgeActionMeasure(g, actionVertexId, measureID)
		if err != nil {
			log.Fatal(err)
		}
	}

	return actionVertexId, nil
}

// fucntion to insert an edge between a stage and a measure
// returns an error if insertion fails, else returns nil
func edgeStageMeasure(g *gremlingo.GraphTraversalSource, stageID int64, measureID int64) error {

	_, err := g.AddE("has").From(g.V(stageID)).To(g.V(measureID)).Next()

	if err != nil {
		return errors.New("failed to create an edge from stage to measure" + "\nError:" + err.Error())
	}

	return nil
}

// fucntion to insert an edge between a measure and a stage
// returns an error if insertion fails, else returns nil
func edgeOperationMeasure(g *gremlingo.GraphTraversalSource, operationID int64, measureID int64) error {

	_, err := g.AddE("has").From(g.V(operationID)).To(g.V(measureID)).Next()

	if err != nil {
		return errors.New("failed to create an edge from operations to  measureID" + "\nError:" + err.Error())
	}

	return nil
}

// fucntion to insert an edge between a measure and a stage
// returns an error if insertion fails, else returns nil
func test(g *gremlingo.GraphTraversalSource) (string, error) {

	res, err := g.AddV("test").Property("name", "helloooo").Next()

	if err != nil {
		return "", errors.New("failed to create an edge from operations:" + "to measure" + "\nError:" + err.Error())
	}
	vertex, err := res.GetVertex()
	if err != nil {
		return "", errors.New("Not a vertex" + err.Error())
	}
	println(vertex.String())
	//println(vertex.Element)
	println(vertex.Label)
	println(vertex.Id.(int64))
	println(vertex.Properties)

	result, _ := g.V(vertex.Id.(int64)).ToList()

	println(result[0].GetString())

	return vertex.String(), nil

}

// fucntion to insert an edge between a stage and a measure
// returns an error if insertion fails, else returns nil
func edgeActionMeasure(g *gremlingo.GraphTraversalSource, actionID int64, measureID int64) error {

	_, err := g.AddE("has").From(g.V(actionID)).To(g.V(measureID)).Next()

	if err != nil {
		return errors.New("failed to create an edge from action measure" + "\nError:" + err.Error())
	}

	return nil
}

// fucntion to insert a Measure vertex within the graph
// returns an error if insertion fails, else returns nil
func insertMeasure(g *gremlingo.GraphTraversalSource, connectID int64, measure Measure) (int64, error) {
	measureName := measure.Measure
	measureID := measure.MeasureID

	res, err := g.AddV("measure").Property("name", measureName).Property("measureId", measureID).Next()

	if err != nil {
		return 0, errors.New("failed to execute query for measureID:" + measureID + "\nError:" + err.Error())
	}

	vertex, err := res.GetVertex()
	if err != nil {
		return 0, errors.New("Not a vertex" + err.Error())
	}

	return vertex.Id.(int64), nil
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
		fmt.Println(vertex.GetString())
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
		fmt.Println(vertex.GetString())
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
		fmt.Println(vertex.GetString())
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
		fmt.Println(vertex.GetString())
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
		fmt.Println(vertex.GetString())
	}
	elapsedTime := end.Sub(start)
	fmt.Println("Elapsed Time:", elapsedTime)
}
