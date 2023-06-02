package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type lambdaAPI interface {
	CreateFunction(ctx context.Context, params *lambda.CreateFunctionInput, optFns ...func(*lambda.Options)) (*lambda.CreateFunctionOutput, error)
	DeleteFunction(ctx context.Context, params *lambda.DeleteFunctionInput, optFns ...func(*lambda.Options)) (*lambda.DeleteFunctionOutput, error)
	CreateEventSourceMapping(ctx context.Context, params *lambda.CreateEventSourceMappingInput, optFns ...func(*lambda.Options)) (*lambda.CreateEventSourceMappingOutput, error)
}

// Lambda is a wrapper around the AWS Lambda client.
type Lambda struct {
	client lambdaAPI
}

// NewLambda creates a new Lambda client with the given configuration.
func NewLambda(config aws.Config) *Lambda {
	return &Lambda{
		client: lambda.NewFromConfig(config.Copy()),
	}
}

// CreateGo creates a Lambda function from a Go binary. The binary must be zipped and
// uploaded to S3. The bucketName and bucketKey parameters are the name of the bucket.
// It will return the ARN of the Lambda function and an error if there is one.
func (l *Lambda) CreateGo(name, bucketName, bucketObjectKey string) (string, error) {
	createOutput, err := l.client.CreateFunction(context.TODO(), &lambda.CreateFunctionInput{
		Code: &types.FunctionCode{
			S3Bucket: aws.String(bucketName),
			S3Key:    aws.String(bucketObjectKey),
		},
		FunctionName: aws.String(name),
		Handler:      aws.String("main"),
		Runtime:      types.RuntimeGo1x,
		Role:         aws.String("arn:aws:iam::000000000000:role/lambda-role"),
		Timeout:      aws.Int32(60),
		MemorySize:   aws.Int32(128),
		Publish:      true,
		Environment:  &types.Environment{},
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(createOutput.FunctionArn), nil
}

// CreateNode creates a Lambda function from a Node.js binary. The binary must be zipped
// and uploaded to S3. The bucketName and bucketKey parameters are the name of the bucket.
// It will return the ARN of the Lambda function and an error if there is one.
func (l *Lambda) CreateNode(name, bucketName, bucketObjecyKey string) (string, error) {
	createOutput, err := l.client.CreateFunction(context.TODO(), &lambda.CreateFunctionInput{
		Code: &types.FunctionCode{
			S3Bucket: aws.String(bucketName),
			S3Key:    aws.String(bucketObjecyKey),
		},
		FunctionName: aws.String(name),
		Handler:      aws.String("index.handler"),
		Runtime:      types.RuntimeNodejs16x,
		Role:         aws.String("arn:aws:iam::123456789012:role/lambda-role"),
		Timeout:      aws.Int32(60),
		MemorySize:   aws.Int32(128),
		Publish:      true,
		Environment:  &types.Environment{},
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(createOutput.FunctionArn), nil
}

// Delete deletes a Lambda function with the given name.
func (l *Lambda) Delete(name string) error {
	_, err := l.client.DeleteFunction(context.TODO(), &lambda.DeleteFunctionInput{
		FunctionName: aws.String(name),
	})
	if err != nil {
		return err
	}

	return nil
}

// BindToService binds a Lambda function to an event source. This can be used to bind a
// Lambda function to an SQS queue or an SNS topic. For instance, if you want to bind a
// Lambda function to Kinesis, you would pass in the ARN of the Kinesis stream as the
// eventSourceArn parameter.
func (l *Lambda) BindToService(name, eventSourceArn string) error {
	_, err := l.client.CreateEventSourceMapping(context.TODO(), &lambda.CreateEventSourceMappingInput{
		FunctionName:     aws.String(name),
		EventSourceArn:   aws.String(eventSourceArn),
		BatchSize:        aws.Int32(100),
		StartingPosition: types.EventSourcePositionLatest,
	})
	if err != nil {
		return err
	}

	return nil
}
