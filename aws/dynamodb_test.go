package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type mockDynamoDBClient struct {
	createTableFunc   func(context.Context, *dynamodb.CreateTableInput, ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	updateTableFunc   func(context.Context, *dynamodb.UpdateTableInput, ...func(*dynamodb.Options)) (*dynamodb.UpdateTableOutput, error)
	deleteTableFunc   func(context.Context, *dynamodb.DeleteTableInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error)
	putItemFunc       func(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	deleteItemFunc    func(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	describeTableFunc func(context.Context, *dynamodb.DescribeTableInput, ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
}

func (m *mockDynamoDBClient) CreateTable(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	return m.createTableFunc(ctx, params, optFns...)
}

func (m *mockDynamoDBClient) UpdateTable(ctx context.Context, params *dynamodb.UpdateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateTableOutput, error) {
	return m.updateTableFunc(ctx, params, optFns...)
}

func (m *mockDynamoDBClient) DeleteTable(ctx context.Context, params *dynamodb.DeleteTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error) {
	return m.deleteTableFunc(ctx, params, optFns...)
}

func (m *mockDynamoDBClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m.putItemFunc(ctx, params, optFns...)
}

func (m *mockDynamoDBClient) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return m.deleteItemFunc(ctx, params, optFns...)
}

func (m *mockDynamoDBClient) DescribeTable(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
	return m.describeTableFunc(ctx, params, optFns...)
}

func TestDynamoDB_CreateTable(t *testing.T) {
	mockClient := &mockDynamoDBClient{
		createTableFunc: func(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
			return &dynamodb.CreateTableOutput{}, nil
		},
	}

	dynamoDB := &DynamoDB{
		client: mockClient,
	}

	err := dynamoDB.CreateTable("test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDynamoDB_UpdateReplicas(t *testing.T) {
	mockClient := &mockDynamoDBClient{
		updateTableFunc: func(ctx context.Context, params *dynamodb.UpdateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateTableOutput, error) {
			return &dynamodb.UpdateTableOutput{}, nil
		},
	}

	dynamoDB := &DynamoDB{
		client: mockClient,
	}

	err := dynamoDB.UpdateReplicas("test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDynamoDB_DeleteTable(t *testing.T) {
	mockClient := &mockDynamoDBClient{
		deleteTableFunc: func(ctx context.Context, params *dynamodb.DeleteTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error) {
			return &dynamodb.DeleteTableOutput{}, nil
		},
	}

	dynamoDB := &DynamoDB{
		client: mockClient,
	}

	err := dynamoDB.DeleteTable("test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDynamoDB_DescribeTable(t *testing.T) {
	mockClient := &mockDynamoDBClient{
		describeTableFunc: func(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
			return &dynamodb.DescribeTableOutput{
				Table: &types.TableDescription{
					TableName: aws.String("test"),
				},
			}, nil
		},
	}

	dynamoDB := &DynamoDB{
		client: mockClient,
	}

	table, err := dynamoDB.DescribeTable("test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if aws.ToString(table.Table.TableName) != "test" {
		t.Errorf("unexpected table name: %v", table.Table.TableName)
	}
}

func TestDynamoDB_PutItem(t *testing.T) {
	mockClient := &mockDynamoDBClient{
		putItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	}

	dynamoDB := &DynamoDB{
		client: mockClient,
	}

	err := dynamoDB.PutItem("test", map[string]types.AttributeValue{"test": &types.AttributeValueMemberS{Value: "test"}})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDynamoDB_DeleteItem(t *testing.T) {
	mockClient := &mockDynamoDBClient{
		deleteItemFunc: func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
			return &dynamodb.DeleteItemOutput{}, nil
		},
	}

	dynamoDB := &DynamoDB{
		client: mockClient,
	}

	err := dynamoDB.DeleteItem("test", map[string]types.AttributeValue{"test": &types.AttributeValueMemberS{Value: "test"}})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
