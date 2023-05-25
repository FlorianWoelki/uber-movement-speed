package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBAPI interface {
	CreateTable(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	DeleteTable(ctx context.Context, params *dynamodb.DeleteTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	UpdateTable(ctx context.Context, params *dynamodb.UpdateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateTableOutput, error)
	DescribeTable(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
}

type DynamoDB struct {
	client DynamoDBAPI
}

func NewDynamoDB(config aws.Config) *DynamoDB {
	return &DynamoDB{
		client: dynamodb.NewFromConfig(config),
	}
}

func (d *DynamoDB) CreateTable(name string) error {
	_, err := d.client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(name),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DynamoDB) UpdateReplicas(name string) error {
	_, err := d.client.UpdateTable(context.TODO(), &dynamodb.UpdateTableInput{
		TableName: aws.String(name),
		ReplicaUpdates: []types.ReplicationGroupUpdate{
			{
				Create: &types.CreateReplicationGroupMemberAction{
					RegionName: aws.String("eu-central-1"),
				},
			},
			{
				Create: &types.CreateReplicationGroupMemberAction{
					RegionName: aws.String("us-west-1"),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DynamoDB) DeleteTable(name string) error {
	_, err := d.client.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(name),
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DynamoDB) DescribeTable(name string) (*dynamodb.DescribeTableOutput, error) {
	return d.client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	})
}

func (d *DynamoDB) PutItem(name string, item map[string]types.AttributeValue) error {
	_, err := d.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(name),
		Item:      item,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DynamoDB) DeleteItem(name string, key map[string]types.AttributeValue) error {
	_, err := d.client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(name),
		Key:       key,
	})
	if err != nil {
		return err
	}

	return nil
}
