package main

import (
	"fmt"
	//"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// PopulateDatabase populates the database for a given process name
func PopulateDatabase(processName string) (status string, duration string, errMsg string) {
	fmt.Print("Inserting " + processName)

	// Define the table name inside the function
	tableName := "Merck-Fall2024-Final-2"

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return "Failed", "0ms", fmt.Sprintf("Got error creating session: %v", err)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Record the start time
	startTime := time.Now()

	// Generate data using the process name
	//TODO: call ModelDataParent() instead
	data, err := ModelData(processName)

	if err != nil {
		return "Failed", "0ms", fmt.Sprintf("Got error generating model data: %v", err)
	}

	// PROCESS
	processItem := map[string]*dynamodb.AttributeValue{
		"EntityType": {S: aws.String("Process")},
		"EntityID":   {S: aws.String(data.Hierarchy.Process)},
	}

	// Insert the process
	_, err = svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      processItem,
	})
	if err != nil {
		return "Failed", "0ms", fmt.Sprintf("Got error putting process item: %v", err)
	}

	// Iterate over the stages
	for _, stage := range data.Hierarchy.Stages {
		// STAGE
		stageItem := map[string]*dynamodb.AttributeValue{
			"EntityType": {S: aws.String("Stage")},
			"EntityID":   {S: aws.String(stage.Stage)},
		}

		// Insert the stage
		_, err = svc.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      stageItem,
		})
		if err != nil {
			return "Failed", "0ms", fmt.Sprintf("Got error putting stage item: %v", err)
		}

		// Iterate over the operations
		for _, operation := range stage.Operations {
			// OPERATION
			operationItem := map[string]*dynamodb.AttributeValue{
				"EntityType": {S: aws.String("Operation")},
				"EntityID":   {S: aws.String(operation.Operation)},
				"ParentID":   {S: aws.String(stage.Stage)},
			}

			// Insert the operation
			_, err = svc.PutItem(&dynamodb.PutItemInput{
				TableName: aws.String(tableName),
				Item:      operationItem,
			})
			if err != nil {
				return "Failed", "0ms", fmt.Sprintf("Got error putting operation item: %v", err)
			}

			// Iterate over the actions
			for _, action := range operation.Actions {
				// ACTION
				actionItem := map[string]*dynamodb.AttributeValue{
					"EntityType": {S: aws.String("Action")},
					"EntityID":   {S: aws.String(action.Action)},
					"ParentID":   {S: aws.String(operation.Operation)},
				}

				// Insert the action
				_, err = svc.PutItem(&dynamodb.PutItemInput{
					TableName: aws.String(tableName),
					Item:      actionItem,
				})
				if err != nil {
					return "Failed", "0ms", fmt.Sprintf("Got error putting action item: %v", err)
				}

				// Iterate over the measures
				for _, measure := range action.Measures {
					// MEASURE
					measureItem := map[string]*dynamodb.AttributeValue{
						"EntityType": {S: aws.String("Measure")},
						"EntityID":   {S: aws.String(measure.MeasureID)},
						"Measure":    {S: aws.String(strings.TrimPrefix(measure.Measure, "Measure: "))},
						"ParentID":   {S: aws.String(action.Action)},
					}

					// Insert the measure
					_, err = svc.PutItem(&dynamodb.PutItemInput{
						TableName: aws.String(tableName),
						Item:      measureItem,
					})
					if err != nil {
						return "Failed", "0ms", fmt.Sprintf("Got error putting measure item: %v", err)
					}
				}
			}
		}
	}

	// Iterate over xpath
	for _, xpath := range data.Xpath {
		// XPATH
		xpathItem := map[string]*dynamodb.AttributeValue{
			"EntityType": {S: aws.String("Xpath")},
			"EntityID":   {S: aws.String(xpath.Xpath)},
			"MeasureID":  {S: aws.String(xpath.MeasureID)},
			"Site":       {S: aws.String(strings.TrimPrefix(xpath.Site, "Site- "))},
		}

		// Inserting the xpath
		_, err = svc.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      xpathItem,
		})
		if err != nil {
			return "Failed", "0ms", fmt.Sprintf("Got error putting xpath item: %v", err)
		}
	}

	// Iterate over metadata
	for _, metadata := range data.Metadata {
		// METADATA
		metadataItem := map[string]*dynamodb.AttributeValue{
			"EntityType": {S: aws.String("Metadata")},
			"EntityID":   {S: aws.String(metadata.MeasureID)},
			"Key":        {S: aws.String(strings.TrimPrefix(metadata.Key, "Key "))},
			"Value":      {S: aws.String(strings.TrimPrefix(metadata.Value, "Value "))},
		}

		// Inserting the metadata
		_, err = svc.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      metadataItem,
		})
		if err != nil {
			return "Failed", "0ms", fmt.Sprintf("Got error putting metadata item: %v", err)
		}
	}

	// Iterate over the results
	for _, result := range data.Results {
		// RESULT
		resultItem := map[string]*dynamodb.AttributeValue{
			"EntityType": {S: aws.String("Result")},
			"EntityID":   {S: aws.String(result.Xpath + "#" + result.BatchID)},
			"Xpath":      {S: aws.String(result.Xpath)},
			"Site":       {S: aws.String(strings.TrimPrefix(result.Site, "Site- "))},
			"BatchID":    {N: aws.String(result.BatchID)},
			"DOM":        {S: aws.String(result.DOM)},
			"ResultName": {S: aws.String(result.ResultName)},
			"Result":     {N: aws.String(fmt.Sprintf("%f", result.Result))},
		}

		// Inserting the result
		_, err = svc.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      resultItem,
		})
		if err != nil {
			return "Failed", "0ms", fmt.Sprintf("Got error putting result item: %v", err)
		}
	}

	// Iterate over raw materials
	for _, rawMaterial := range data.RawMaterials {
		// RAWMATERIALS
		rawMaterialItem := map[string]*dynamodb.AttributeValue{
			"EntityType":        {S: aws.String("RawMaterial")},
			"EntityID":          {S: aws.String(rawMaterial.ParentBatchID + "#" + rawMaterial.ChildMaterialName)},
			"ParentBatchID":     {S: aws.String(rawMaterial.ParentBatchID)},
			"ParentMaterialNum": {S: aws.String(rawMaterial.ParentMaterialNum)},
			"ChildMaterialName": {S: aws.String(rawMaterial.ChildMaterialName)},
			"ChildBatchID":      {S: aws.String(rawMaterial.ChildBatchID)},
			"ChildMaterialNum":  {S: aws.String(rawMaterial.ChildMaterialNum)},
		}

		// Inserting raw materials
		_, err = svc.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      rawMaterialItem,
		})
		if err != nil {
			return "Failed", "0ms", fmt.Sprintf("Got error putting raw material item: %v", err)
		}
	}

	// Record the end time
	endTime := time.Now()
	durationTime := endTime.Sub(startTime)

	// Log the success message
	fmt.Printf("Process '%s' successfully loaded into the database in %.2f ms\n", processName, float64(durationTime.Milliseconds()))

	fmt.Print("Completed Inserting" + processName)

	// Return success with status, duration, and no error message
	return "OK", fmt.Sprintf("%.2f ms", float64(durationTime.Milliseconds())), ""
}
