package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsService "github.com/florianwoelki/uber-movement-speed/aws"
)

type Book struct {
	Id    string `json:"id" dynamodbav:"id"`
	Title string `json:"title" dynamodbav:"title"`
}

var (
	tableName    string
	s3BucketName string
)

var (
	dataBatch []Book
	batchSize = 1000
)

// Used clients for the AWS services.
var (
	dynamodbClient *awsService.DynamoDB
	s3Client       *awsService.S3
)

func init() {
	// tableName = os.Getenv("TABLE_NAME")
	tableName = "books"
	s3BucketName = "raw-data"
	if tableName == "" {
		log.Fatal("missing environment variable`TABLE_NAME`")
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           "http://s3.localhost.localstack.cloud:4566",
				SigningRegion: "us-east-1",
			}, nil
		}

		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           fmt.Sprintf("http://%s:4566", os.Getenv("LOCALSTACK_HOSTNAME")),
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

	dynamodbClient = awsService.NewDynamoDB(cfg)
	s3Client = awsService.NewS3(cfg)
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

		// Prepare the item to be stored in the DynamoDB table.
		item, err := attributevalue.MarshalMap(book)
		if err != nil {
			return err
		}

		// Stores the item in the DynamoDB table.
		err = dynamodbClient.PutItem(tableName, item)
		if err != nil {
			return err
		}

		log.Println("Successfully put item into dynamodb table")

		// Accumulate data for batch upload to S3.
		dataBatch = append(dataBatch, book)

		if len(dataBatch) > batchSize {
			if err := uploadToS3(dataBatch); err != nil {
				return fmt.Errorf("failed to upload data to S3: %v", err)
			}
			dataBatch = nil
		}
	}

	flushBatch()
	return nil
}

// flushBatch flushes the current batch to the S3 bucket.
func flushBatch() {
	if len(dataBatch) > 0 {
		if err := uploadToS3(dataBatch); err != nil {
			log.Printf("Failed to upload data to S3: %v", err)
		}
		dataBatch = nil
	}
}

// uploadToS3 uploads the given data to the S3 bucket as a CSV file.
func uploadToS3(data []Book) error {
	log.Println("Uploading data to S3")
	csvBuffer := new(bytes.Buffer)
	writer := csv.NewWriter(csvBuffer)

	// Write headers to CSV file.
	headers := []string{"id", "title"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, book := range data {
		if err := writer.Write([]string{book.Id, book.Title}); err != nil {
			return err
		}
	}
	writer.Flush()

	key := fmt.Sprintf("batch-from-%s-to-%s.csv", data[0].Id, data[len(data)-1].Id)
	err := s3Client.PutObject(s3BucketName, key, csvBuffer.Bytes())
	if err != nil {
		return err
	}
	log.Println("Successfully uploaded data to S3")
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
