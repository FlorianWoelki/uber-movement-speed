package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/glue"
)

type mockGlueAPI struct {
	createJobFn func(ctx context.Context, params *glue.CreateJobInput, optFns ...func(*glue.Options)) (*glue.CreateJobOutput, error)
}

func (m *mockGlueAPI) CreateJob(ctx context.Context, params *glue.CreateJobInput, optFns ...func(*glue.Options)) (*glue.CreateJobOutput, error) {
	return m.createJobFn(ctx, params, optFns...)
}

func TestGlue_CreateJob(t *testing.T) {
	mockClient := &mockGlueAPI{
		createJobFn: func(ctx context.Context, params *glue.CreateJobInput, optFns ...func(*glue.Options)) (*glue.CreateJobOutput, error) {
			return &glue.CreateJobOutput{}, nil
		},
	}

	g := &Glue{client: mockClient}

	err := g.CreateJob("test-job", "s3://test-bucket/test-job.py")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
