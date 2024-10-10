package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// GetAllOperations retrieves all operations from the specified DynamoDB table.
func GetAllOperations(tableName string) {
	// Start the timer
	startTime := time.Now()

	// Session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)

	// Handling error
	if err != nil {
		fmt.Println("Got error creating session:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// DynamoDB client
	svc := dynamodb.New(sess)

	// Query to get all operations
	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
		FilterExpression: aws.String("EntityType = :entityType"), // Filtering the results based on attribute value
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":entityType": {S: aws.String("Operation")}, // Filtering by operation
		},
	}

	// Results of the filter
	result, err := svc.Scan(params)

	// Error handling
	if err != nil {
		fmt.Println("Got error scanning table:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Calculate the duration
	duration := time.Since(startTime)

	// Print the list of operations
	fmt.Println("Retrieved operations:")
	for _, item := range result.Items {
		fmt.Printf("Operation: %v\n", item)
	}

	// Printing duration
	fmt.Printf("Time taken to retrieve data: %v\n", duration)
}
