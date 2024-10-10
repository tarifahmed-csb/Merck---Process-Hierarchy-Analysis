package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

//Response Structure
type Item struct {
	ParentID string `json:"ParentID"`
	EntityID string `json:"EntityID"`
}

func accessPattern(tableName string) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n\t1. Filter operation by stage")
		fmt.Println("\t2. Filter all actions by operation")
		fmt.Println("\t3. Filter all measures by action")
		fmt.Println("\t4. Exit")

		// Prompt for the option
		fmt.Print("\n\tEnter your choice (1-4): ")
		scanner.Scan()
		choice := scanner.Text()

		// Exit if user chooses option 4
		if choice == "4" {
			fmt.Println("Exiting...")
			break
		}

		// Prompt for appropriate
		var idType string
		switch choice {
		case "1":
			idType = "Stage ID"
		case "2":
			idType = "Operation ID"
		case "3":
			idType = "Action ID"
		default:
			fmt.Println("Invalid choice, please enter a number between 1 and 4.")
			continue
		}

		fmt.Printf("\tEnter the %s: ", idType)
		scanner.Scan()
		id := scanner.Text()

		// New session 
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-2"),
		})
		if err != nil {
			return fmt.Errorf("failed to create session: %v", err)
		}

		// DynamoDB client 
		svc := dynamodb.New(sess)

		// QUERY PREP
		params := &dynamodb.QueryInput{
			TableName: aws.String(tableName),
			IndexName: aws.String("GSI2"), 
			KeyConditionExpression: aws.String("ParentID = :id"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":id": {
					S: aws.String(id),
				},
			},
		}

		// Start timing
		startTime := time.Now()

		//QUERY
		result, err := svc.Query(params)
		if err != nil {
			return fmt.Errorf("failed to query items: %v", err)
		}

		// Stop timing
		duration := time.Since(startTime)

		// Parse result
		var entityIDs []string
		for _, item := range result.Items {
			var i Item
			err = dynamodbattribute.UnmarshalMap(item, &i)
			if err != nil {
				return fmt.Errorf("failed to unmarshal item: %v", err)
			}
			entityIDs = append(entityIDs, i.EntityID)
		}

		// Print results
		fmt.Printf("\n\t%s '%s':\n", idType, id)
		for _, entityID := range entityIDs {
			fmt.Printf("\t\t%s\n", entityID)
		}

		// Print the duration
		fmt.Printf("\n\tQuery completed in %s\n", duration)
	}

	return nil
}
