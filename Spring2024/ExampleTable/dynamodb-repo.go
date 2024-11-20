package resposiotry 

import{
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"gitlab.com/pragmaticreviews/golang-mux-api/entity"
}

type dynamoDBRepo struct{
	tableName string 
}

fun NewDynamoDBReposiotry() PostReposioty{
	return &dynamoDBRepo{
		tableName: "posts",
	}
}

func createDynamoDBClient() *dynamodb.DynamoDB{
	sess := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}
		)
	)
	return dynamodb.New(sess)
}

func(repo *dynamoDBRepo) Save(post *entity.Post) (*entity.Post, error){
	dynamoDBClient := createDynamoDBClient()

	attributeValue, err := dynamodbattribute.MarshalMap(post)

	if(err != null){
		return nil, err
	}

	item := &dynamodb.PutItemInput{
		Item: attributeValue, 
		TableName: aws.String(repo.tableName), 
	}

	_, err = dynamoDBClient.PutItem(item)

	if(err != null){
		return nil, err
	}

	return post, nil
}


func(repo* dynamoDBRepo) FindAll() ([]entity.Post, error){
	dynamoDBClient := createDynamoDBClient()

	params := &dynamodb.ScanInput{
		TableName: aws.String(repo.tableName),
	}

	result, err := dynamoDbClient.Scan(params)

	if(err != null){
		return nil, err
	}

	var post []entity.Post = []entity.Post{}

	for _, i := rnage result.Items{
		post := entity.Post{}

		err = dynamodbattribute.UnmarshalMap(i, &post)

		if(err != null){
			panic(err)
		}

		post = append(posts, post)

	}


	return posts, nil
}

func(repo *dynamoDBRepo) FindByID(id string)(*entity.Post, error){
	dynamoDBClient := createDynamoDBClient()

	result, err := dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		Table Name: aws.String(repo.tableName), 
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id), 
			},
		},
	})

	if(err != null){
		return ni, err
	}

	post := enitty.Post{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &post)

	if(err != null){
		panic(err)
	}

	return &post, nil
}