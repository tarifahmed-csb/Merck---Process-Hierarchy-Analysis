// package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"strconv"
// 	"time"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/dynamodb"
// 	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
// )

// type Item struct {
// 	PROCESSID    string
// 	ProcessName  string
// 	MeasuresP    []Measure
// 	Stage        string
// 	MeasuresS    []Measure
// 	Operation    string
// 	Action       string
// 	Measure      string
// 	MeasureID    string
// 	Site         string
// 	XPath        string
// 	MaterialNum  string
// 	DOM          string
// 	Result       string
// 	BatchID      string
// 	ParentMatNum string
// 	ChildMatName string
// 	ChildBatchID string
// 	ChildMatNum  string
// 	Key          string
// 	Value        string
// }

// type Measure struct {
// 	Measure   string
// 	MeasureID string
// }

// func main() {
// 	// Initialize a session that the SDK will use to load credentials from the shared credentials file
// 	// and region from the shared configuration file
// 	sess := session.Must(session.NewSessionWithOptions(session.Options{
// 		SharedConfigState: session.SharedConfigEnable,
// 	}))

// 	// Create DynamoDB client
// 	svc := dynamodb.New(sess)

// 	// Generate and insert mock data into the DynamoDB table
// 	err := generateAndInsertMockData(svc)
// 	if err != nil {
// 		fmt.Printf("Error generating and inserting mock data: %v\n", err)
// 		return
// 	}

// 	fmt.Println("Mock data insertion completed successfully.")
// }

// func generateAndInsertMockData(svc *dynamodb.DynamoDB) error {
// 	// Generate a mock item
// 	mockItem := generateMockItem()

// 	// Marshal the mock item into a DynamoDB attribute value map
// 	av, err := dynamodbattribute.MarshalMap(mockItem)
// 	if err != nil {
// 		return fmt.Errorf("error marshalling item: %v", err)
// 	}

// 	// Specify the table name
// 	tableName := "Merck"

// 	// Create a PutItemInput object
// 	input := &dynamodb.PutItemInput{
// 		Item:      av,
// 		TableName: aws.String(tableName),
// 	}

// 	// Insert the item into the DynamoDB table
// 	_, err = svc.PutItem(input)
// 	if err != nil {
// 		return fmt.Errorf("error inserting item: %v", err)
// 	}

// 	return nil
// }

// func generateMockItem() Item {
// 	rand.Seed(time.Now().UnixNano())

// 	// Generate random values for the mock item
// 	processID := "6789"
// 	processName := "Sample Process"
// 	stage := "Sample Stage"
// 	operation := "Sample Operation"
// 	action := "Sample Action"
// 	measure := "Sample Measure"
// 	measureID := "Sample Measure ID"
// 	site := "Sample Site"
// 	xPath := "Sample XPath"
// 	materialNum := "Sample MaterialNum"
// 	dom := time.Now().Format("2006-01-02")
// 	result := strconv.FormatFloat(rand.Float64()*100, 'f', 2, 64) // Generating a random float between 0 and 100 for the result, then converting to string
// 	batchID := "Sample BatchID"
// 	parentMatNum := "Sample ParentMatNum"
// 	childMatName := "Sample ChildMatName"
// 	childBatchID := "Sample ChildBatchID"
// 	childMatNum := "Sample ChildMatNum"
// 	key := "Sample Key"
// 	value := "Sample Value"

// 	// Create the mock item
// 	mockItem := Item{
// 		PROCESSID:    processID,
// 		ProcessName:  processName,
// 		Stage:        stage,
// 		Operation:    operation,
// 		Action:       action,
// 		Measure:      measure,
// 		MeasureID:    measureID,
// 		Site:         site,
// 		XPath:        xPath,
// 		MaterialNum:  materialNum,
// 		DOM:          dom,
// 		Result:       result,
// 		BatchID:      batchID,
// 		ParentMatNum: parentMatNum,
// 		ChildMatName: childMatName,
// 		ChildBatchID: childBatchID,
// 		ChildMatNum:  childMatNum,
// 		Key:          key,
// 		Value:        value,
// 	}

// 	return mockItem
// }
