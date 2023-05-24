package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type mockLambdaClient struct {
	createFunctionFunc func(context.Context, *lambda.CreateFunctionInput, ...func(*lambda.Options)) (*lambda.CreateFunctionOutput, error)
	deleteFunctionFunc func(context.Context, *lambda.DeleteFunctionInput, ...func(*lambda.Options)) (*lambda.DeleteFunctionOutput, error)
}

func (m *mockLambdaClient) CreateFunction(ctx context.Context, input *lambda.CreateFunctionInput, opts ...func(*lambda.Options)) (*lambda.CreateFunctionOutput, error) {
	return m.createFunctionFunc(ctx, input, opts...)
}

func (m *mockLambdaClient) DeleteFunction(ctx context.Context, input *lambda.DeleteFunctionInput, opts ...func(*lambda.Options)) (*lambda.DeleteFunctionOutput, error) {
	return m.deleteFunctionFunc(ctx, input, opts...)
}

func TestLambda_CreateGo(t *testing.T) {
	mockClient := &mockLambdaClient{
		createFunctionFunc: func(ctx context.Context, input *lambda.CreateFunctionInput, opts ...func(*lambda.Options)) (*lambda.CreateFunctionOutput, error) {
			if *input.FunctionName != "test-function" {
				t.Errorf("unexpected function name: %s", *input.FunctionName)
			}
			return &lambda.CreateFunctionOutput{
				FunctionArn: input.FunctionName,
			}, nil
		},
	}

	lambdaClient := &Lambda{
		client: mockClient,
	}

	_, err := lambdaClient.CreateGo("test-function", "test-bucket", "test-key")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLambda_Delete(t *testing.T) {
	mockClient := &mockLambdaClient{
		deleteFunctionFunc: func(ctx context.Context, input *lambda.DeleteFunctionInput, opts ...func(*lambda.Options)) (*lambda.DeleteFunctionOutput, error) {
			if *input.FunctionName != "test-function" {
				t.Errorf("unexpected function name: %s", *input.FunctionName)
			}
			return &lambda.DeleteFunctionOutput{}, nil
		},
	}

	lambdaClient := &Lambda{
		client: mockClient,
	}

	err := lambdaClient.Delete("test-function")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
