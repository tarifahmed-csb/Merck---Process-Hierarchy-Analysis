package main

import (
	"fmt"
	//"os"
	"time"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// PopulateDatabase populates the database for a given process name
func PopulateDatabase(processName string) (status string, duration string, errMsg string) {
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

	// Return success with status, duration, and no error message
	return "OK", fmt.Sprintf("%.2f ms", float64(durationTime.Milliseconds())), ""
}

// func PopulateDatabase(tableName string) {
// 	// Record the start time
// 	startTime := time.Now()

// 	//Session 
// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String("us-east-1")},
// 	)

// 	//Handling error
// 	if err != nil {
// 		fmt.Println("Got error creating session:")
// 		fmt.Println(err.Error())
// 		os.Exit(1)
// 	}

// 	//DynamoDB client
// 	svc := dynamodb.New(sess)

// 	//Error handling
// 	if err != nil {
// 		fmt.Println("Got error connecting to DynamoDB:")
// 		fmt.Println(err.Error())
// 		os.Exit(1)
// 	}

// 	fmt.Println("Successfully connected to Merck-Fall2024")

// 	// Genertating data
// 	data, err := ModelData("Process1")

// 	//Error handling
// 	if err != nil {
// 		fmt.Println("Got error generating model data:")
// 		fmt.Println(err.Error())
// 		os.Exit(1)
// 	}

// 	/*
// 	"Attribute Name" {dataType: aws.String("     " + data.Hierarcy.... )}

// 	aws.String -  pointer to a string
// 	"   " + _____ - forming a string value for attribute 
// 	*/

// 	//PROCESS
// 	processItem := map[string]*dynamodb.AttributeValue{
// 		//"PK":         {S: aws.String("Process#" + data.Hierarchy.Process)}, //PRIMARY KEY
// 		//"SK":         {S: aws.String("Process#" + data.Hierarchy.Process)}, //SORT KEY
// 		"EntityType": {S: aws.String("Process")},
// 		"EntityID":   {S: aws.String(data.Hierarchy.Process)},
// 	}

// 	//Inserting Process
// 	_, err = svc.PutItem(&dynamodb.PutItemInput{
// 		TableName: aws.String(tableName),
// 		Item:      processItem,
// 	})
// 	//Error Handling: Process
// 	if err != nil {
// 		fmt.Println("Got error putting process item:")
// 		fmt.Println(err.Error())
// 		os.Exit(1)
// 	}

// 	// Iterating over the stages
// 	for _, stage := range data.Hierarchy.Stages {
// 		//STAGE
// 		stageItem := map[string]*dynamodb.AttributeValue{
// 			//"PK":         {S: aws.String(stage.Stage)},
// 			//"SK":         {S: aws.String(stage.Stage)},
// 			"EntityType": {S: aws.String("Stage")},
// 			"EntityID":   {S: aws.String(stage.Stage)},
// 		}

// 		//Inserting the stage 
// 		_, err = svc.PutItem(&dynamodb.PutItemInput{
// 			TableName: aws.String(tableName),
// 			Item:      stageItem,
// 		})
// 		//Error Handling: Stage 
// 		if err != nil {
// 			fmt.Println("Got error putting stage item:")
// 			fmt.Println(err.Error())
// 			os.Exit(1)
// 		}

// 		//Iterating over the operations
// 		for _, operation := range stage.Operations {
// 			//OPERATION
// 			operationItem := map[string]*dynamodb.AttributeValue{
// 				//"PK":         {S: aws.String("Operation#" + operation.Operation)},
// 				//"SK":         {S: aws.String("Operation#" + operation.Operation)},
// 				"EntityType": {S: aws.String("Operation")},
// 				"EntityID":   {S: aws.String(operation.Operation)},
// 				"ParentID":   {S: aws.String(stage.Stage)},
// 			}
// 			//Inserting the operation
// 			_, err = svc.PutItem(&dynamodb.PutItemInput{
// 				TableName: aws.String(tableName),
// 				Item:      operationItem,
// 			})
// 			//Error handling: operation
// 			if err != nil {
// 				fmt.Println("Got error putting operation item:")
// 				fmt.Println(err.Error())
// 				os.Exit(1)
// 			}

// 			//Iterating over the actions
// 			for _, action := range operation.Actions {
// 				//ACTION
// 				actionItem := map[string]*dynamodb.AttributeValue{
// 					//"PK":         {S: aws.String(action.Action)},
// 					//"SK":         {S: aws.String(action.Action)},
// 					"EntityType": {S: aws.String("Action")},
// 					"EntityID":   {S: aws.String(action.Action)},
// 					"ParentID":   {S: aws.String(operation.Operation)},
// 				}
// 				//Inserting the action
// 				_, err = svc.PutItem(&dynamodb.PutItemInput{
// 					TableName: aws.String(tableName),
// 					Item:      actionItem,
// 				})
// 				//Error Handling: action
// 				if err != nil {
// 					fmt.Println("Got error putting action item:")
// 					fmt.Println(err.Error())
// 					os.Exit(1)
// 				}
// 				//Iterating over the measures
// 				for _, measure := range action.Measures {
// 					//fmt.Println("Made it in to the for loop")
// 					//MEASURE
// 					measureItem := map[string]*dynamodb.AttributeValue{
// 						//"PK":         {S: aws.String("Measure#" + measure.MeasureID)},
// 						//"SK":         {S: aws.String("Measure#" + measure.MeasureID)},
// 						"EntityType": {S: aws.String("Measure")},
// 						"EntityID":   {S: aws.String(measure.MeasureID)},
// 						"Measure": 	  {S: aws.String(strings.TrimPrefix(measure.Measure, "Measure: "))},
// 						"ParentID":   {S: aws.String(action.Action)},
// 					}
// 					//Inserting the measure 
// 					_, err = svc.PutItem(&dynamodb.PutItemInput{
// 						TableName: aws.String(tableName),
// 						Item:      measureItem,
// 					})
// 					//Error handling: measure 
// 					if err != nil {
// 						fmt.Println("Got error putting measure item:")
// 						fmt.Println(err.Error())
// 						os.Exit(1)
// 					}
// 				}
// 			}
// 		}
// 	}

// 	//Iterating over xpath 
// 	for _, xpath := range data.Xpath {
// 		//XPATH
// 		xpathItem := map[string]*dynamodb.AttributeValue{
// 			//"PK":         {S: aws.String("Xpath#" + xpath.Xpath)},
// 			//"SK":         {S: aws.String("Xpath#" + xpath.Xpath)},
// 			"EntityType": {S: aws.String("Xpath")},
// 			"EntityID":   {S: aws.String(xpath.Xpath)},
// 			"MeasureID":  {S: aws.String(xpath.MeasureID)},
// 			"Site":       {S: aws.String(strings.TrimPrefix((xpath).Site, "Site- "))},
// 		}
// 		//Inserting the xpath
// 		_, err = svc.PutItem(&dynamodb.PutItemInput{
// 			TableName: aws.String(tableName),
// 			Item:      xpathItem,
// 		})
// 		//Error handling: xpath 
// 		if err != nil {
// 			fmt.Println("Got error putting xpath item:")
// 			fmt.Println(err.Error())
// 			os.Exit(1)
// 		}
// 	}

// 	//Iterating over the metadata
// 	for _, metadata := range data.Metadata {
// 		//METADATA
// 		metadataItem := map[string]*dynamodb.AttributeValue{
// 			//"PK":         {S: aws.String("Metadata#" + metadata.MeasureID + "#" + metadata.Key)},
// 			//"SK":         {S: aws.String("Metadata#" + metadata.MeasureID + "#" + metadata.Key)},
// 			"EntityType": {S: aws.String("Metadata")},
// 			"EntityID":   {S: aws.String(metadata.MeasureID)},
// 			"Key":        {S: aws.String(strings.TrimPrefix(metadata.Key, "Key "))},
// 			"Value":      {S: aws.String(strings.TrimPrefix(metadata.Value, "Value "))},
// 		}
// 		//Inserting the metadata
// 		_, err = svc.PutItem(&dynamodb.PutItemInput{
// 			TableName: aws.String(tableName),
// 			Item:      metadataItem,
// 		})
// 		//Error handling: metadata
// 		if err != nil {
// 			fmt.Println("Got error putting metadata item:")
// 			fmt.Println(err.Error())
// 			os.Exit(1)
// 		}
// 	}

// 	//Iterating over the results
// 	for _, result := range data.Results {
// 		//RESULT
// 		resultItem := map[string]*dynamodb.AttributeValue{
// 			//"PK":         {S: aws.String("Result#" + result.Xpath + "#" + result.BatchID)},
// 			//"SK":         {S: aws.String("Result#" + result.Xpath + "#" + result.BatchID)},
// 			"EntityType": {S: aws.String("Result")},
// 			"EntityID":   {S: aws.String(result.Xpath + "#" + result.BatchID)},
// 			"Xpath":      {S: aws.String(result.Xpath)},
// 			"Site":       {S: aws.String(strings.TrimPrefix(result.Site, "Site- "))},
// 			"BatchID":    {N: aws.String(result.BatchID)},
// 			"DOM":        {S: aws.String(result.DOM)},
// 			"ResultName": {S: aws.String(result.ResultName)},
// 			"Result":     {N: aws.String(fmt.Sprintf("%f", result.Result))},
// 		}
// 		//Inserting the result 
// 		_, err = svc.PutItem(&dynamodb.PutItemInput{
// 			TableName: aws.String(tableName),
// 			Item:      resultItem,
// 		})
// 		//Error hanlding: result
// 		if err != nil {
// 			fmt.Println("Got error putting result item:")
// 			fmt.Println(err.Error())
// 			os.Exit(1)
// 		}
// 	}

// 	//Iterating over rawmaterials
// 	for _, rawMaterial := range data.RawMaterials {
// 		fmt.Println("made it to raw materials")
// 		//RAWMATERIALS
// 		rawMaterialItem := map[string]*dynamodb.AttributeValue{
// 			//"PK":                 {S: aws.String("RawMaterial#" + rawMaterial.ParentBatchID + "#" + rawMaterial.ChildMaterialName)},
// 			//"SK":                 {S: aws.String("RawMaterial#" + rawMaterial.ParentBatchID + "#" + rawMaterial.ChildMaterialName)},
// 			"EntityType":         {S: aws.String("RawMaterial")},
// 			"EntityID":           {S: aws.String(rawMaterial.ParentBatchID + "#" + rawMaterial.ChildMaterialName)},
// 			"ParentBatchID":      {S: aws.String(rawMaterial.ParentBatchID)},
// 			"ParentMaterialNum":  {S: aws.String(rawMaterial.ParentMaterialNum)},
// 			"ChildMaterialName":  {S: aws.String(rawMaterial.ChildMaterialName)},
// 			"ChildBatchID":       {S: aws.String(rawMaterial.ChildBatchID)},
// 			"ChildMaterialNum":   {S: aws.String(rawMaterial.ChildMaterialNum)},
// 		}
// 		//Inserting rawmaterials
// 		_, err = svc.PutItem(&dynamodb.PutItemInput{
// 			TableName: aws.String(tableName),
// 			Item:      rawMaterialItem,
// 		})
// 		//Error handling: rawmaterials
// 		if err != nil {
// 			fmt.Println("Got error putting raw material item:")
// 			fmt.Println(err.Error())
// 			os.Exit(1)
// 		}
// 	}

// 	//End timer
// 	endTime := time.Now()
// 	duration := endTime.Sub(startTime) //total time 

// 	//Success
// 	fmt.Printf("Successfully added model data to the Merck-Fall2024 table\n")
// 	fmt.Printf("Time taken: %.2f ms\n", float64(duration.Milliseconds()))
// }
