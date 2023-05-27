package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

type kinesisAPI interface {
	CreateStream(ctx context.Context, params *kinesis.CreateStreamInput, optFns ...func(*kinesis.Options)) (*kinesis.CreateStreamOutput, error)
	DeleteStream(ctx context.Context, params *kinesis.DeleteStreamInput, optFns ...func(*kinesis.Options)) (*kinesis.DeleteStreamOutput, error)
	PutRecord(ctx context.Context, params *kinesis.PutRecordInput, optFns ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error)
}

// Kinesis is a wrapper around the AWS Kinesis client.
type Kinesis struct {
	client kinesisAPI
}

// NewKinesis creates a new Kinesis client with the given configuration.
func NewKinesis(config aws.Config) *Kinesis {
	return &Kinesis{
		client: kinesis.NewFromConfig(config),
	}
}

// Create creates a Kinesis stream with the given name and sets the shard count to `1`.
func (k *Kinesis) Create(name string) error {
	_, err := k.client.CreateStream(context.TODO(), &kinesis.CreateStreamInput{
		ShardCount: aws.Int32(1),
		StreamName: aws.String(name),
	})
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a Kinesis stream with the given name.
func (k *Kinesis) Delete(name string) error {
	_, err := k.client.DeleteStream(context.TODO(), &kinesis.DeleteStreamInput{
		StreamName: aws.String(name),
	})
	if err != nil {
		return err
	}

	return nil
}

// PutRecord puts a record into a Kinesis stream with the given name and partition key.
func (k *Kinesis) PutRecord(name, partitionKey string, data []byte) error {
	_, err := k.client.PutRecord(context.TODO(), &kinesis.PutRecordInput{
		Data:         data,
		PartitionKey: aws.String(partitionKey),
		StreamName:   aws.String(name),
	})
	if err != nil {
		return err
	}

	return nil
}
