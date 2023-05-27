package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Book struct {
	Title string `json:"title"`
}

func handleRequest(ctx context.Context, book Book) (map[string]interface{}, error) {
	return map[string]interface{}{
		"statusCode": 200,
		"body":       fmt.Sprintf("Book: %s", book.Title),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
