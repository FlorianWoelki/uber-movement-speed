package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type LambdaAPI interface {
	CreateFunction(ctx context.Context, params *lambda.CreateFunctionInput, optFns ...func(*lambda.Options)) (*lambda.CreateFunctionOutput, error)
	DeleteFunction(ctx context.Context, params *lambda.DeleteFunctionInput, optFns ...func(*lambda.Options)) (*lambda.DeleteFunctionOutput, error)
}

// Lambda is a wrapper around the AWS Lambda client.
type Lambda struct {
	client LambdaAPI
}

// NewLambda creates a new Lambda client with the given configuration.
func NewLambda(config aws.Config) *Lambda {
	return &Lambda{
		client: lambda.NewFromConfig(config),
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
