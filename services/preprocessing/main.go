package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Book struct {
	Title string `json:"title"`
}

func handleRequest(ctx context.Context, event events.KinesisEvent) error {
	for _, record := range event.Records {
		kinesisRecord := record.Kinesis
		dataBytes := kinesisRecord.Data

		var book Book
		jsonParsingErr := json.Unmarshal(dataBytes, &book)
		if jsonParsingErr != nil {
			return jsonParsingErr
		}

		fmt.Printf("%s Data = %s \n", record.EventName, book.Title)
	}
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
