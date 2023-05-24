package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

type APIGatewayAPI interface {
	CreateRestApi(ctx context.Context, params *apigateway.CreateRestApiInput, optFns ...func(*apigateway.Options)) (*apigateway.CreateRestApiOutput, error)
	DeleteRestApi(ctx context.Context, params *apigateway.DeleteRestApiInput, optFns ...func(*apigateway.Options)) (*apigateway.DeleteRestApiOutput, error)
	GetResources(ctx context.Context, params *apigateway.GetResourcesInput, optFns ...func(*apigateway.Options)) (*apigateway.GetResourcesOutput, error)
	CreateResource(ctx context.Context, params *apigateway.CreateResourceInput, optFns ...func(*apigateway.Options)) (*apigateway.CreateResourceOutput, error)
	PutMethod(ctx context.Context, params *apigateway.PutMethodInput, optFns ...func(*apigateway.Options)) (*apigateway.PutMethodOutput, error)
	PutIntegration(ctx context.Context, params *apigateway.PutIntegrationInput, optFns ...func(*apigateway.Options)) (*apigateway.PutIntegrationOutput, error)
	CreateDeployment(ctx context.Context, params *apigateway.CreateDeploymentInput, optFns ...func(*apigateway.Options)) (*apigateway.CreateDeploymentOutput, error)
}

type APIGateway struct {
	client APIGatewayAPI
}

func NewAPIGateway(config aws.Config) *APIGateway {
	return &APIGateway{
		client: apigateway.NewFromConfig(config),
	}
}

func (a *APIGateway) Create(name string) (string, error) {
	createOutput, err := a.client.CreateRestApi(context.TODO(), &apigateway.CreateRestApiInput{
		Name: aws.String(name),
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(createOutput.Id), nil
}

func (a *APIGateway) Delete(id string) error {
	_, err := a.client.DeleteRestApi(context.TODO(), &apigateway.DeleteRestApiInput{
		RestApiId: aws.String(id),
	})
	if err != nil {
		return err
	}

	return nil
}

type EndpointOptions struct {
	// Path is the path of the resource where the endpoint will be created.
	Path string
	// Method is the HTTP method of the endpoint.
	Method string
	// Uri is the URI of the endpoint.
	Uri string
}

// CreateEndpoint creates an endpoint for the given API Gateway ID with the given options.
// It creates a resource with the given path, a method with the given HTTP method,
// an integration with the given URI, and a deployment.
func (a *APIGateway) CreateEndpoint(id string, options EndpointOptions) error {
	parentId, err := a.getRootId(id)
	if err != nil {
		return err
	}

	createResourceOutput, err := a.client.CreateResource(context.TODO(), &apigateway.CreateResourceInput{
		RestApiId: aws.String(id),
		ParentId:  aws.String(parentId),
		PathPart:  aws.String(options.Path),
	})
	if err != nil {
		return err
	}
	resourceId := aws.ToString(createResourceOutput.Id)

	_, err = a.client.PutMethod(context.TODO(), &apigateway.PutMethodInput{
		RestApiId:         aws.String(id),
		ResourceId:        aws.String(resourceId),
		HttpMethod:        aws.String(options.Method),
		RequestParameters: map[string]bool{fmt.Sprintf("method.request.path.%s", options.Path): true},
		AuthorizationType: aws.String("NONE"),
	})
	if err != nil {
		return err
	}

	_, err = a.client.PutIntegration(context.TODO(), &apigateway.PutIntegrationInput{
		RestApiId:             aws.String(id),
		ResourceId:            aws.String(resourceId),
		HttpMethod:            aws.String(options.Method),
		Type:                  types.IntegrationTypeHttpProxy,
		IntegrationHttpMethod: aws.String("POST"),
		Uri:                   aws.String(options.Uri),
		PassthroughBehavior:   aws.String("WHEN_NO_MATCH"),
	})
	if err != nil {
		return err
	}

	_, err = a.client.CreateDeployment(context.TODO(), &apigateway.CreateDeploymentInput{
		RestApiId: aws.String(id),
		StageName: aws.String("dev"),
	})
	if err != nil {
		return err
	}

	return nil
}

// getRootId returns the ID of the root resource for the given API Gateway ID.
// The root resource is the resource with the path `/`.
func (a *APIGateway) getRootId(id string) (string, error) {
	getResourcesOutput, err := a.client.GetResources(context.TODO(), &apigateway.GetResourcesInput{
		RestApiId: aws.String(id),
	})
	if err != nil {
		return "", err
	}

	for _, resource := range getResourcesOutput.Items {
		if aws.ToString(resource.Path) == "/" {
			return aws.ToString(resource.Id), nil
		}
	}

	return "", fmt.Errorf("unable to find root resource for API Gateway ID %s", id)
}
