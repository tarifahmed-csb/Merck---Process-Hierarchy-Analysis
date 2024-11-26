package main

import (
	//"bufio"
	"fmt"
	//"os"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Entity struct {
	PK         string `json:"PK"`
	SK         string `json:"SK"`
	EntityID   string `json:"EntityID"`
	EntityType string `json:"EntityType"` 
}

func query(prefix, entityType string) (status string, duration string, results string, errMsg string) {
	// Establish the table name
	tableName := "Merck-Fall2024-Final-2"

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return "Failed", "0ms", "", fmt.Sprintf("Failed to create session: %v", err)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Trim the input prefix
	prefix = strings.TrimSpace(prefix)

	// Input for Scan
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
		FilterExpression: aws.String("begins_with(EntityID, :prefix) AND EntityType = :entityType"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":prefix": {
				S: aws.String(prefix),
			},
			":entityType": {
				S: aws.String(entityType),
			},
		},
	}

	// Start timing
	startTime := time.Now()

	// Execute the scan
	result, err := svc.Scan(input)
	if err != nil {
		return "Failed", "0ms", "", fmt.Sprintf("Failed to scan the table: %v", err)
	}

	// Stop timing
	durationTime := time.Since(startTime)

	// Unmarshal the result into a slice of Entity structs
	var entities []Entity
	var res string

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &entities)
	if err != nil {
		return "Failed", "0ms", "", fmt.Sprintf("Failed to unmarshal results: %v", err)
	}

	// Sort the results
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].EntityID < entities[j].EntityID
	})

	// Check if any entities are found
	if len(entities) == 0 {
		return "OK", fmt.Sprintf("%.2f ms", float64(durationTime.Milliseconds())), "", "No items found"
	} else {
		// Output the results if needed (optional)
		// For example, just for debugging purposes, you could return these as part of the response
		// This step can be skipped if not needed
		for _, entity := range entities {
			fmt.Printf("%s: %s\n", entityType, entity.EntityID)
			res += fmt.Sprintf("%s: %s", entityType, entity.EntityID)
		}
	}

	// Return the successful execution details
	return "OK", fmt.Sprintf("%.2f ms", float64(durationTime.Milliseconds())), res, ""
}

// func query(tableName, entityType string) {
// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String("us-east-1"),
// 	})
// 	if err != nil {
// 		fmt.Println("Failed to create session,", err)
// 		return
// 	}
// 	svc := dynamodb.New(sess)

// 	// Prompt user
// 	fmt.Printf("Enter the starting value for EntityID for %s: ", entityType)
// 	scanner := bufio.NewScanner(os.Stdin)
// 	scanner.Scan()
// 	prefix := strings.TrimSpace(scanner.Text())

// 	// Input for Scan
// 	input := &dynamodb.ScanInput{
// 		TableName: aws.String(tableName),
// 		FilterExpression: aws.String("begins_with(EntityID, :prefix) AND EntityType = :entityType"),
// 		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
// 			":prefix": {
// 				S: aws.String(prefix),
// 			},
// 			":entityType": {
// 				S: aws.String(entityType),
// 			},
// 		},
// 	}

// 	// Start timing
// 	startTime := time.Now()

// 	// Execute the scan
// 	result, err := svc.Scan(input)
// 	if err != nil {
// 		fmt.Println("Failed to scan the table,", err)
// 		return
// 	}

// 	// Stop timing
// 	duration := time.Since(startTime)

// 	// Unmarshal the result into a slice of Entity structs
// 	var entities []Entity
// 	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &entities)
// 	if err != nil {
// 		fmt.Println("Failed to unmarshal results,", err)
// 		return
// 	}

// 	// Sort the results
// 	sort.Slice(entities, func(i, j int) bool {
// 		return entities[i].EntityID < entities[j].EntityID
// 	})

// 	// Results
// 	if len(entities) == 0 {
// 		fmt.Printf("No items found for %s.\n", entityType)
// 	} else {
// 		fmt.Printf("Items found for %s:\n", entityType)
// 		for _, entity := range entities {
// 			fmt.Printf("%s: %s\n", entityType, entity.EntityID)
// 		}
// 	}

// 	// Execution time
// 	fmt.Printf("Query executed in: %v\n", duration)
// }