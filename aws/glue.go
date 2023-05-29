package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"
)

type glueAPI interface {
	CreateJob(ctx context.Context, params *glue.CreateJobInput, optFns ...func(*glue.Options)) (*glue.CreateJobOutput, error)
}

// Glue is a wrapper around the AWS Glue API.
type Glue struct {
	client glueAPI
}

// NewGlue creates a new Glue client.
func NewGlue(config aws.Config) *Glue {
	return &Glue{client: glue.NewFromConfig(config)}
}

// CreateJob creates a new Glue job with the given name and script location.
func (g *Glue) CreateJob(jobName, scriptLocation string) error {
	_, err := g.client.CreateJob(context.TODO(), &glue.CreateJobInput{
		Name: aws.String(jobName),
		Role: aws.String("arn:aws:iam::000000000000:role/glue-role"),
		Command: &types.JobCommand{
			Name:           aws.String("pythonshell"),
			ScriptLocation: aws.String(scriptLocation),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
