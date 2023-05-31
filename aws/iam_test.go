package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type mockIAMAPI struct {
	createRoleFn       func(ctx context.Context, params *iam.CreateRoleInput, opts ...func(*iam.Options)) (*iam.CreateRoleOutput, error)
	createPolicyFn     func(ctx context.Context, params *iam.CreatePolicyInput, opts ...func(*iam.Options)) (*iam.CreatePolicyOutput, error)
	attachRolePolicyFn func(ctx context.Context, params *iam.AttachRolePolicyInput, opts ...func(*iam.Options)) (*iam.AttachRolePolicyOutput, error)
}

func (m *mockIAMAPI) CreateRole(ctx context.Context, params *iam.CreateRoleInput, opts ...func(*iam.Options)) (*iam.CreateRoleOutput, error) {
	return m.createRoleFn(ctx, params, opts...)
}

func (m *mockIAMAPI) CreatePolicy(ctx context.Context, params *iam.CreatePolicyInput, opts ...func(*iam.Options)) (*iam.CreatePolicyOutput, error) {
	return m.createPolicyFn(ctx, params, opts...)
}

func (m *mockIAMAPI) AttachRolePolicy(ctx context.Context, params *iam.AttachRolePolicyInput, opts ...func(*iam.Options)) (*iam.AttachRolePolicyOutput, error) {
	return m.attachRolePolicyFn(ctx, params, opts...)
}

func TestIAM_CreateRoleWithPolicy(t *testing.T) {
	mockClient := &mockIAMAPI{
		createRoleFn: func(ctx context.Context, params *iam.CreateRoleInput, opts ...func(*iam.Options)) (*iam.CreateRoleOutput, error) {
			return &iam.CreateRoleOutput{
				Role: &types.Role{
					RoleName: aws.String(""),
				},
			}, nil
		},
		createPolicyFn: func(ctx context.Context, params *iam.CreatePolicyInput, opts ...func(*iam.Options)) (*iam.CreatePolicyOutput, error) {
			return &iam.CreatePolicyOutput{
				Policy: &types.Policy{
					Arn: aws.String(""),
				},
			}, nil
		},
		attachRolePolicyFn: func(ctx context.Context, params *iam.AttachRolePolicyInput, opts ...func(*iam.Options)) (*iam.AttachRolePolicyOutput, error) {
			return &iam.AttachRolePolicyOutput{}, nil
		},
	}

	iam := &IAM{iamClient: mockClient}

	_, err := iam.CreateRoleWithPolicy("test-role", "test-service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
