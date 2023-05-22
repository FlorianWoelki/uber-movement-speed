package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Payload struct {
	Body       string `json:"body"`
	StatusCode int    `json:"statusCode"`
}

func main() {
	lambda.Start(func(ctx context.Context) (events.APIGatewayProxyResponse, error) {
		headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept"}
		p := Payload{
			Body:       "Hello, World!",
			StatusCode: 200,
		}
		return events.APIGatewayProxyResponse{
			StatusCode:      p.StatusCode,
			Body:            p.Body,
			Headers:         headers,
			IsBase64Encoded: false,
		}, nil
	})
}
