package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsService "github.com/florianwoelki/uber-movement-speed/aws"
)

type Book struct {
	Id    string `json:"id" dynamodbav:"id"`
	Title string `json:"title" dynamodbav:"title"`
}

var tableName string
var client *awsService.DynamoDB

func init() {
	// tableName = os.Getenv("TABLE_NAME")
	tableName = "books"
	if tableName == "" {
		log.Fatal("missing environment variable`TABLE_NAME`")
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           fmt.Sprintf("http://%s:4566", os.Getenv("LOCALSTACK_HOSTNAME")),
				SigningRegion: "us-east-1",
			}, nil
		}

		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           "http://localhost:4566",
			SigningRegion: "us-east-1",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		log.Fatal(err)
	}

	client = awsService.NewDynamoDB(cfg)
}

func handleRequest(ctx context.Context, event events.KinesisEvent) error {
	for _, record := range event.Records {
		log.Printf("Received message from kinesis. partition key: %s\n", record.Kinesis.PartitionKey)
		log.Printf("Storing information to dynamodb table: %s\n", tableName)

		kinesisRecord := record.Kinesis
		dataBytes := kinesisRecord.Data

		var book Book
		err := json.Unmarshal(dataBytes, &book)
		if err != nil {
			return err
		}

		item, err := attributevalue.MarshalMap(book)
		if err != nil {
			return err
		}

		err = client.PutItem(tableName, item)
		if err != nil {
			return err
		}

		log.Println("Successfully put item into dynamodb table")
	}
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
