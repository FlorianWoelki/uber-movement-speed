package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsService "github.com/florianwoelki/uber-movement-speed/aws"
)

type SegmentSpeed struct {
	Id              string  `json:"id" dynamodbav:"id"`
	Year            int     `json:"year" dynamodbav:"year"`
	Month           int     `json:"month" dynamodbav:"month"`
	Day             int     `json:"day" dynamodbav:"day"`
	Hour            int     `json:"hour" dynamodbav:"hour"`
	UtcTimestamp    string  `json:"utc_timestamp" dynamodbav:"utc_timestamp"`
	StartJunctionId string  `json:"start_junction_id" dynamodbav:"start_junction_id"`
	EndJunctionId   string  `json:"end_junction_id" dynamodbav:"end_junction_id"`
	OsmWayId        int64   `json:"osm_way_id" dynamodbav:"osm_way_id"`
	OsmStartNodeId  int64   `json:"osm_start_node_id" dynamodbav:"osm_start_node_id"`
	OsmEndNodeId    int64   `json:"osm_end_node_id" dynamodbav:"osm_end_node_id"`
	SpeedMphMean    float32 `json:"speed_mph_mean" dynamodbav:"speed_mph_mean"`
	SpeedMphStddev  float32 `json:"speed_mph_stddev" dynamodbav:"speed_mph_stddev"`
}

var (
	tableName    string
	s3BucketName string
)

var (
	dataBatch []SegmentSpeed
	batchSize = 1000
)

// Used clients for the AWS services.
var (
	dynamodbClient *awsService.DynamoDB
	s3Client       *awsService.S3
)

func init() {
	tableName = "street_segment_speeds"
	s3BucketName = "raw-data"

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

		var segmentSpeed SegmentSpeed
		err := json.Unmarshal(dataBytes, &segmentSpeed)
		if err != nil {
			return err
		}

		// Prepare the item to be stored in the DynamoDB table.
		item, err := attributevalue.MarshalMap(segmentSpeed)
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
		dataBatch = append(dataBatch, segmentSpeed)

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

// uploadToS3 uploads the given data to the S3 bucket as a CSV file and partitions it by
// the current time.
func uploadToS3(data []SegmentSpeed) error {
	// Gets the current time to construct the partition path.
	currentTime := time.Now()
	year := currentTime.Format("2006")
	month := currentTime.Format("01")
	day := currentTime.Format("02")

	// Construct the partition path based on the timestamp.
	partitionPath := fmt.Sprintf("year=%s/month=%s/day=%s/",
		year, month, day)

	log.Println("Uploading data to S3")
	csvBuffer := new(bytes.Buffer)
	writer := csv.NewWriter(csvBuffer)

	// Write headers to CSV file.
	headers := []string{"id", "year", "month", "day", "hour", "utc_timestamp", "start_junction_id", "end_junction_id", "osm_way_id", "osm_start_node_id", "osm_end_node_id", "speed_mph_mean", "speed_mph_stddev"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, segmentSpeed := range data {
		if err := writer.Write(
			[]string{
				segmentSpeed.Id,
				fmt.Sprint(segmentSpeed.Year),
				fmt.Sprint(segmentSpeed.Month),
				fmt.Sprint(segmentSpeed.Day),
				fmt.Sprint(segmentSpeed.Hour),
				segmentSpeed.UtcTimestamp,
				segmentSpeed.StartJunctionId,
				segmentSpeed.EndJunctionId,
				fmt.Sprint(segmentSpeed.OsmWayId),
				fmt.Sprint(segmentSpeed.OsmStartNodeId),
				fmt.Sprint(segmentSpeed.OsmEndNodeId),
				fmt.Sprint(segmentSpeed.SpeedMphMean),
				fmt.Sprint(segmentSpeed.SpeedMphStddev),
			},
		); err != nil {
			return err
		}
	}
	writer.Flush()

	key := fmt.Sprintf("batch-from-%s-to-%s.csv", data[0].Id, data[len(data)-1].Id)
	keyWithPartition := partitionPath + key
	err := s3Client.PutObject(s3BucketName, keyWithPartition, csvBuffer.Bytes())
	if err != nil {
		return err
	}
	log.Println("Successfully uploaded data to S3")
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
