package main

//Importing AWS DynamoDB packages
import (
    "fmt" 
    "os"

    "github.com/aws/aws-sdk-go/aws" //credentials, regions, and settings
    "github.com/aws/aws-sdk-go/aws/session" //create and manag a session
    "github.com/aws/aws-sdk-go/service/dynamodb" //DynamoDB client and API
)

//function just in case I need to create a new table
func DeleteAllItemsFromTable(tableName string) {
    //Session
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-2")},
    )

	//Handling error
    if err != nil {
        fmt.Println("Got error creating session:")
        fmt.Println(err.Error())
        os.Exit(1)
    }

    // DynamoDB client
    svc := dynamodb.New(sess)

    //Scanning table 
    scanInput := &dynamodb.ScanInput{
        TableName: aws.String(tableName),
    }

    //Last Evaluated Key
    var lastEvaluatedKey map[string]*dynamodb.AttributeValue
    for {
        if lastEvaluatedKey != nil {
            scanInput.ExclusiveStartKey = lastEvaluatedKey
        }

		//Scanning input
        result, err := svc.Scan(scanInput)

		//Error handling
        if err != nil {
            fmt.Println("Got error scanning table:")
            fmt.Println(err.Error())
            os.Exit(1)
        }

        //Iterating over each item 
        for _, item := range result.Items {
            // Key map for item to be deleted
            key := map[string]*dynamodb.AttributeValue{
                "PK": {S: item["PK"].S},
                "SK": {S: item["SK"].S},
            }

			//Input structure for deleting
            deleteInput := &dynamodb.DeleteItemInput{
                TableName: aws.String(tableName),
                Key:       key,
            }

			//DELETION
            _, err := svc.DeleteItem(deleteInput)

			//Error handling
            if err != nil {
                fmt.Println("Got error deleting item:")
                fmt.Println(err.Error())
                os.Exit(1)
            }
        }

        //Evaluating remaining data
        if result.LastEvaluatedKey == nil {
            break
        }
        lastEvaluatedKey = result.LastEvaluatedKey
    }

    fmt.Println("Successfully deleted all items from the DynamoDB table")
}