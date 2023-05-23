package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

type mockKinesisClient struct {
	createStreamFunc func(context.Context, *kinesis.CreateStreamInput, ...func(*kinesis.Options)) (*kinesis.CreateStreamOutput, error)
	deleteStreamFunc func(context.Context, *kinesis.DeleteStreamInput, ...func(*kinesis.Options)) (*kinesis.DeleteStreamOutput, error)
	putRecordFunc    func(context.Context, *kinesis.PutRecordInput, ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error)
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

}

func TestKinesis_PutRecord(t *testing.T) {

}
