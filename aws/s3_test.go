package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type mockS3Client struct {
	createBucketFunc func(context.Context, *s3.CreateBucketInput, ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	putObjectFunc    func(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	deleteBucketFunc func(context.Context, *s3.DeleteBucketInput, ...func(*s3.Options)) (*s3.DeleteBucketOutput, error)
}

func (m *mockS3Client) CreateBucket(ctx context.Context, input *s3.CreateBucketInput, opts ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return m.createBucketFunc(ctx, input, opts...)
}

func (m *mockS3Client) PutObject(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return m.putObjectFunc(ctx, input, opts...)
}

func (m *mockS3Client) DeleteBucket(ctx context.Context, input *s3.DeleteBucketInput, opts ...func(*s3.Options)) (*s3.DeleteBucketOutput, error) {
	return m.deleteBucketFunc(ctx, input, opts...)
}

func TestS3_CreateBucket(t *testing.T) {
	mockClient := &mockS3Client{
		createBucketFunc: func(ctx context.Context, input *s3.CreateBucketInput, opts ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
			if aws.ToString(input.Bucket) != "test-bucket" {
				t.Errorf("unexpected bucket name: %s", aws.ToString(input.Bucket))
			}
			return &s3.CreateBucketOutput{}, nil
		},
	}

	s3Client := &S3{
		client: mockClient,
	}

	err := s3Client.CreateBucket("test-bucket")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestS3_DeleteBucket(t *testing.T) {
	mockClient := &mockS3Client{
		deleteBucketFunc: func(ctx context.Context, input *s3.DeleteBucketInput, opts ...func(*s3.Options)) (*s3.DeleteBucketOutput, error) {
			if aws.ToString(input.Bucket) != "test-bucket" {
				t.Errorf("unexpected bucket name: %s", aws.ToString(input.Bucket))
			}
			return &s3.DeleteBucketOutput{}, nil
		},
	}

	s3Client := &S3{
		client: mockClient,
	}

	err := s3Client.DeleteBucket("test-bucket")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestS3_PutObject(t *testing.T) {
	mockClient := &mockS3Client{
		putObjectFunc: func(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			if aws.ToString(input.Bucket) != "test-bucket" {
				t.Errorf("unexpected bucket name: %s", aws.ToString(input.Bucket))
			}
			if aws.ToString(input.Key) != "test-key" {
				t.Errorf("unexpected key name: %s", aws.ToString(input.Key))
			}
			return &s3.PutObjectOutput{}, nil
		},
	}

	s3Client := &S3{
		client: mockClient,
	}

	err := s3Client.PutObject("test-bucket", "test-key", []byte("test-data"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
