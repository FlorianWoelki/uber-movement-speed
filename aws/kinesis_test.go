package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
)

type mockKinesisClient struct {
	createStreamFunc   func(context.Context, *kinesis.CreateStreamInput, ...func(*kinesis.Options)) (*kinesis.CreateStreamOutput, error)
	deleteStreamFunc   func(context.Context, *kinesis.DeleteStreamInput, ...func(*kinesis.Options)) (*kinesis.DeleteStreamOutput, error)
	putRecordFunc      func(context.Context, *kinesis.PutRecordInput, ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error)
	describeStreamFunc func(context.Context, *kinesis.DescribeStreamInput, ...func(*kinesis.Options)) (*kinesis.DescribeStreamOutput, error)
}

func (m *mockKinesisClient) CreateStream(ctx context.Context, input *kinesis.CreateStreamInput, opts ...func(*kinesis.Options)) (*kinesis.CreateStreamOutput, error) {
	return m.createStreamFunc(ctx, input, opts...)
}

func (m *mockKinesisClient) DeleteStream(ctx context.Context, input *kinesis.DeleteStreamInput, opts ...func(*kinesis.Options)) (*kinesis.DeleteStreamOutput, error) {
	return m.deleteStreamFunc(ctx, input, opts...)
}

func (m *mockKinesisClient) PutRecord(ctx context.Context, input *kinesis.PutRecordInput, opts ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error) {
	return m.putRecordFunc(ctx, input, opts...)
}

func (m *mockKinesisClient) DescribeStream(ctx context.Context, input *kinesis.DescribeStreamInput, opts ...func(*kinesis.Options)) (*kinesis.DescribeStreamOutput, error) {
	return m.describeStreamFunc(ctx, input, opts...)
}

func TestKinesis_Create(t *testing.T) {
	mockClient := &mockKinesisClient{
		createStreamFunc: func(ctx context.Context, input *kinesis.CreateStreamInput, opts ...func(*kinesis.Options)) (*kinesis.CreateStreamOutput, error) {
			if aws.ToString(input.StreamName) != "test-stream" {
				return nil, errors.New("unexpected stream name")
			}
			return &kinesis.CreateStreamOutput{}, nil
		},
	}

	kinesisClient := &Kinesis{
		client: mockClient,
	}

	err := kinesisClient.Create("test-stream")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestKinesis_Delete(t *testing.T) {
	mockClient := &mockKinesisClient{
		deleteStreamFunc: func(ctx context.Context, input *kinesis.DeleteStreamInput, opts ...func(*kinesis.Options)) (*kinesis.DeleteStreamOutput, error) {
			if aws.ToString(input.StreamName) != "test-stream" {
				return nil, errors.New("unexpected stream name")
			}
			return &kinesis.DeleteStreamOutput{}, nil
		},
	}

	kinesisClient := &Kinesis{
		client: mockClient,
	}

	err := kinesisClient.Delete("test-stream")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestKinesis_PutRecord(t *testing.T) {
	mockClient := &mockKinesisClient{
		putRecordFunc: func(ctx context.Context, input *kinesis.PutRecordInput, opts ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error) {
			if aws.ToString(input.StreamName) != "test-stream" {
				return nil, errors.New("unexpected stream name")
			}
			if aws.ToString(input.PartitionKey) != "test-key" {
				return nil, errors.New("unexpected partition key")
			}
			if string(input.Data) != "test-data" {
				return nil, errors.New("unexpected data")
			}
			return &kinesis.PutRecordOutput{}, nil
		},
	}

	kinesisClient := &Kinesis{
		client: mockClient,
	}

	err := kinesisClient.PutRecord("test-stream", "test-key", []byte("test-data"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestKinesis_GetARN(t *testing.T) {
	mockClient := &mockKinesisClient{
		describeStreamFunc: func(ctx context.Context, input *kinesis.DescribeStreamInput, opts ...func(*kinesis.Options)) (*kinesis.DescribeStreamOutput, error) {
			if aws.ToString(input.StreamName) != "test-stream" {
				return nil, errors.New("unexpected stream name")
			}
			return &kinesis.DescribeStreamOutput{
				StreamDescription: &types.StreamDescription{
					StreamARN: aws.String("test-arn"),
				},
			}, nil
		},
	}

	kinesisClient := &Kinesis{
		client: mockClient,
	}

	arn, err := kinesisClient.GetARN("test-stream")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if arn != "test-arn" {
		t.Errorf("unexpected arn: %s", arn)
	}
}
