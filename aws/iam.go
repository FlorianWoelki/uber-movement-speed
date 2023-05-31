package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var policies = map[string]string{
	"s3": `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"s3:PutObject",
					"s3:GetObject",
					"s3:CreateBucket"
				],
				"Resource": [
					"arn:aws:s3:::*/*",
					"arn:aws:s3:::*"
				]
			}
		]
	}`,
}

type iamAPI interface {
	CreateRole(context.Context, *iam.CreateRoleInput, ...func(*iam.Options)) (*iam.CreateRoleOutput, error)
	CreatePolicy(context.Context, *iam.CreatePolicyInput, ...func(*iam.Options)) (*iam.CreatePolicyOutput, error)
	AttachRolePolicy(context.Context, *iam.AttachRolePolicyInput, ...func(*iam.Options)) (*iam.AttachRolePolicyOutput, error)
}

// IAM is a wrapper around the AWS IAM client.
type IAM struct {
	iamClient iamAPI
	stsClient *sts.Client
}

// NewIAM creates a new IAM client with the given configuration.
func NewIAM(config aws.Config) *IAM {
	return &IAM{
		iamClient: iam.NewFromConfig(config),
		stsClient: sts.NewFromConfig(config),
	}
}

// CreateRoleWithPolicy creates a role with the given name and service. It first creates
// a role with the given name and then attaches a policy to it. The policy is defined
// in the policies map. The service is the AWS service that will assume the role.
// It returns a credentials cache that can be used to assume the role.
func (i *IAM) CreateRoleWithPolicy(name, service string) (*aws.CredentialsCache, error) {
	serviceForRole := fmt.Sprintf("%s.amazonaws.com", service)
	// Creates a role for the given service.
	role, err := i.iamClient.CreateRole(context.TODO(), &iam.CreateRoleInput{
		RoleName: aws.String(name),
		AssumeRolePolicyDocument: aws.String(fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Action": "sts:AssumeRole",
					"Principal": {
						"Service": "%s"
					},
					"Effect": "Allow",
				}
			]
		}`, serviceForRole)),
	})
	if err != nil {
		return nil, err
	}

	// Creates a policy for the given service.
	policyOutput, err := i.iamClient.CreatePolicy(context.TODO(), &iam.CreatePolicyInput{
		PolicyDocument: aws.String(policies[service]),
		PolicyName:     aws.String(fmt.Sprintf("%s-policy", name)),
	})
	if err != nil {
		return nil, err
	}

	// Attaches the policy to the role.
	_, err = i.iamClient.AttachRolePolicy(context.TODO(), &iam.AttachRolePolicyInput{
		PolicyArn: policyOutput.Policy.Arn,
		RoleName:  role.Role.RoleName,
	})
	if err != nil {
		return nil, err
	}

	provider := stscreds.NewAssumeRoleProvider(i.stsClient, aws.ToString(role.Role.Arn))
	return aws.NewCredentialsCache(provider), nil
}
