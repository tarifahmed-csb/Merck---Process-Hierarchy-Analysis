package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws" //credentials, regions, and settings
    "github.com/aws/aws-sdk-go/aws/session" //create and manag a session
    "github.com/aws/aws-sdk-go/service/dynamodb" //DynamoDB client and API
)

//function just in case I need to create a new table
func GetAllStages(tableName string) {
	// Start the timer
	startTime := time.Now()

	//Session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)

	// Handling error
	if err != nil {
		fmt.Println("Got error creating session:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//DynamoDB client
	svc := dynamodb.New(sess)

	//Query to get all stages
	params := &dynamodb.ScanInput{ 
		TableName: aws.String(tableName), 
		FilterExpression: aws.String("EntityType = :entityType"), //Filtering the results based on attribute value
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{ //map of attribute values
			":entityType": {S: aws.String("Stage")}, //filtering by stage
		},
	}

	//Results of the filter 
	result, err := svc.Scan(params)

	//Error handling
	if err != nil {
		fmt.Println("Got error scanning table:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// Calculate the duration
	duration := time.Since(startTime)

	// Print the list of stages
	fmt.Println("Retrieved stages:")
	for _, item := range result.Items {
		fmt.Printf("Stage: %v\n", item)
	}

	//Printing duration
	fmt.Printf("Time taken to retrieve data: %v\n", duration)
}
