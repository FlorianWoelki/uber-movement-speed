package aws

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3API interface {
	CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	DeleteBucket(ctx context.Context, params *s3.DeleteBucketInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketOutput, error)
}

// S3 is a wrapper around the AWS S3 client.
type S3 struct {
	client s3API
}

// NewS3 creates a new S3 client with the given configuration.
func NewS3(config aws.Config) *S3 {
	return &S3{
		client: s3.NewFromConfig(config),
	}
}

// CreateBucket creates a S3 bucket with the given name.
func (s *S3) CreateBucket(name string) error {
	_, err := s.client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteBucket deletes a S3 bucket with the given name.
func (s *S3) DeleteBucket(name string) error {
	_, err := s.client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		return err
	}

	return nil
}

// PutObject puts an object into a S3 bucket with the given name and key.
func (s *S3) PutObject(bucket, key string, data []byte) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return err
	}

	return nil
}
