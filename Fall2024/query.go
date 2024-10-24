package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Entity struct {
	PK        string `json:"PK"`
	SK        string `json:"SK"`
	EntityID  string `json:"EntityID"`
	EntityType string `json:"EntityType"` 
}

func query(tableName string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), 
	})
	if err != nil {
		fmt.Println("Failed to create session,", err)
		return
	}
	svc := dynamodb.New(sess)

	// Prompt user 
	fmt.Print("Enter the starting value for EntityID: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	prefix := strings.TrimSpace(scanner.Text())

	// input
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
		FilterExpression: aws.String("begins_with(EntityID, :prefix) AND EntityType = :entityType"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":prefix": {
				S: aws.String(prefix),
			},
			":entityType": {
				S: aws.String("Measure"), 
			},
		},
	}

	// Start timing 
	startTime := time.Now()

	// Execute the scan
	result, err := svc.Scan(input)
	if err != nil {
		fmt.Println("Failed to scan the table,", err)
		return
	}

	// Stop timing
	duration := time.Since(startTime)

	// Unmarshal the result into a slice of Entity structs
	var entities []Entity
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &entities)
	if err != nil {
		fmt.Println("Failed to unmarshal results,", err)
		return
	}

	// Sort the results 
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].EntityID < entities[j].EntityID
	})

	//Results
	if len(entities) == 0 {
		fmt.Println("No items found.")
	} else {
		fmt.Println("Items found:")
		for _, entity := range entities {
			fmt.Printf("Measure: %s\n", entity.EntityID)
		}
	}

	// Execution time
	fmt.Printf("Query executed in: %v\n", duration)
}
