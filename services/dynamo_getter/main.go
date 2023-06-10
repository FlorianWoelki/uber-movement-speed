package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	awsService "github.com/florianwoelki/uber-movement-speed/aws"
)

var (
	tableName string
)

// Used clients for the AWS services.
var (
	dynamodbClient *awsService.DynamoDB
)

type DynamoGetterResponse struct {
	Item map[string]types.AttributeValue `json:"item"`
}

func init() {
	tableName = "street_segment_speeds"

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
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
}

func handleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (DynamoGetterResponse, error) {
	id := event.QueryStringParameters["id"]
	if id == "" {
		return DynamoGetterResponse{}, fmt.Errorf("id is required")
	}

	item, err := dynamodbClient.GetItemById(tableName, id)
	if err != nil {
		return DynamoGetterResponse{}, err
	}

	if item == nil {
		return DynamoGetterResponse{}, fmt.Errorf("item with id %s not found", id)
	}

	return DynamoGetterResponse{
		Item: item,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
