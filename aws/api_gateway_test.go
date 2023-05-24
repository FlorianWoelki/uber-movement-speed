package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

type mockAPIGatewayClient struct {
	createRestAPIFunc    func(context.Context, *apigateway.CreateRestApiInput, ...func(*apigateway.Options)) (*apigateway.CreateRestApiOutput, error)
	deleteRestAPIFunc    func(context.Context, *apigateway.DeleteRestApiInput, ...func(*apigateway.Options)) (*apigateway.DeleteRestApiOutput, error)
	getResourcesFunc     func(context.Context, *apigateway.GetResourcesInput, ...func(*apigateway.Options)) (*apigateway.GetResourcesOutput, error)
	createResourceFunc   func(context.Context, *apigateway.CreateResourceInput, ...func(*apigateway.Options)) (*apigateway.CreateResourceOutput, error)
	putMethodFunc        func(context.Context, *apigateway.PutMethodInput, ...func(*apigateway.Options)) (*apigateway.PutMethodOutput, error)
	putIntegrationFunc   func(context.Context, *apigateway.PutIntegrationInput, ...func(*apigateway.Options)) (*apigateway.PutIntegrationOutput, error)
	createDeploymentFunc func(context.Context, *apigateway.CreateDeploymentInput, ...func(*apigateway.Options)) (*apigateway.CreateDeploymentOutput, error)
}

func (m *mockAPIGatewayClient) CreateRestApi(ctx context.Context, input *apigateway.CreateRestApiInput, opts ...func(*apigateway.Options)) (*apigateway.CreateRestApiOutput, error) {
	return m.createRestAPIFunc(ctx, input, opts...)
}

func (m *mockAPIGatewayClient) DeleteRestApi(ctx context.Context, input *apigateway.DeleteRestApiInput, opts ...func(*apigateway.Options)) (*apigateway.DeleteRestApiOutput, error) {
	return m.deleteRestAPIFunc(ctx, input, opts...)
}

func (m *mockAPIGatewayClient) GetResources(ctx context.Context, input *apigateway.GetResourcesInput, opts ...func(*apigateway.Options)) (*apigateway.GetResourcesOutput, error) {
	return m.getResourcesFunc(ctx, input, opts...)
}

func (m *mockAPIGatewayClient) CreateResource(ctx context.Context, input *apigateway.CreateResourceInput, opts ...func(*apigateway.Options)) (*apigateway.CreateResourceOutput, error) {
	return m.createResourceFunc(ctx, input, opts...)
}

func (m *mockAPIGatewayClient) PutMethod(ctx context.Context, input *apigateway.PutMethodInput, opts ...func(*apigateway.Options)) (*apigateway.PutMethodOutput, error) {
	return m.putMethodFunc(ctx, input, opts...)
}

func (m *mockAPIGatewayClient) PutIntegration(ctx context.Context, input *apigateway.PutIntegrationInput, opts ...func(*apigateway.Options)) (*apigateway.PutIntegrationOutput, error) {
	return m.putIntegrationFunc(ctx, input, opts...)
}

func (m *mockAPIGatewayClient) CreateDeployment(ctx context.Context, input *apigateway.CreateDeploymentInput, opts ...func(*apigateway.Options)) (*apigateway.CreateDeploymentOutput, error) {
	return m.createDeploymentFunc(ctx, input, opts...)
}

func TestAPIGateway_Create(t *testing.T) {
	mockClient := &mockAPIGatewayClient{
		createRestAPIFunc: func(ctx context.Context, input *apigateway.CreateRestApiInput, opts ...func(*apigateway.Options)) (*apigateway.CreateRestApiOutput, error) {
			if *input.Name != "test-api" {
				t.Errorf("unexpected api name: %s", *input.Name)
			}
			return &apigateway.CreateRestApiOutput{
				Id: input.Name,
			}, nil
		},
	}

	apiGatewayClient := &APIGateway{
		client: mockClient,
	}

	id, err := apiGatewayClient.Create("test-api")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if id != "test-api" {
		t.Errorf("unexpected api id: %s", id)
	}
}

func TestAPIGateway_Delete(t *testing.T) {
	mockClient := &mockAPIGatewayClient{
		deleteRestAPIFunc: func(ctx context.Context, input *apigateway.DeleteRestApiInput, opts ...func(*apigateway.Options)) (*apigateway.DeleteRestApiOutput, error) {
			if *input.RestApiId != "test-api" {
				t.Errorf("unexpected api id: %s", *input.RestApiId)
			}
			return &apigateway.DeleteRestApiOutput{}, nil
		},
	}

	apiGatewayClient := &APIGateway{
		client: mockClient,
	}

	err := apiGatewayClient.Delete("test-api")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAPIGateway_getRootId(t *testing.T) {
	mockClient := &mockAPIGatewayClient{
		getResourcesFunc: func(ctx context.Context, input *apigateway.GetResourcesInput, opts ...func(*apigateway.Options)) (*apigateway.GetResourcesOutput, error) {
			if *input.RestApiId != "test-api" {
				t.Errorf("unexpected api id: %s", *input.RestApiId)
			}
			return &apigateway.GetResourcesOutput{
				Items: []types.Resource{
					{
						Id:   aws.String("root-id"),
						Path: aws.String("/"),
					},
				},
			}, nil
		},
	}

	apiGatewayClient := &APIGateway{
		client: mockClient,
	}

	id, err := apiGatewayClient.getRootId("test-api")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if id != "root-id" {
		t.Errorf("unexpected root id: %s", id)
	}
}

func TestAPIGateway_CreateEndpoint(t *testing.T) {
	mockClient := &mockAPIGatewayClient{
		createResourceFunc: func(ctx context.Context, input *apigateway.CreateResourceInput, opts ...func(*apigateway.Options)) (*apigateway.CreateResourceOutput, error) {
			if *input.PathPart != "hello" {
				t.Errorf("unexpected path part: %s", *input.PathPart)
			}
			return &apigateway.CreateResourceOutput{
				Id: input.PathPart,
			}, nil
		},
		putMethodFunc: func(ctx context.Context, input *apigateway.PutMethodInput, opts ...func(*apigateway.Options)) (*apigateway.PutMethodOutput, error) {
			if *input.HttpMethod != "GET" {
				t.Errorf("unexpected http method: %s", *input.HttpMethod)
			}
			return &apigateway.PutMethodOutput{}, nil
		},
		putIntegrationFunc: func(ctx context.Context, input *apigateway.PutIntegrationInput, opts ...func(*apigateway.Options)) (*apigateway.PutIntegrationOutput, error) {
			if *input.HttpMethod != "GET" {
				t.Errorf("unexpected http method: %s", *input.HttpMethod)
			}
			return &apigateway.PutIntegrationOutput{}, nil
		},
		createDeploymentFunc: func(ctx context.Context, input *apigateway.CreateDeploymentInput, opts ...func(*apigateway.Options)) (*apigateway.CreateDeploymentOutput, error) {
			if *input.RestApiId != "test-api" {
				t.Errorf("unexpected api id: %s", *input.RestApiId)
			}
			return &apigateway.CreateDeploymentOutput{}, nil
		},
		getResourcesFunc: func(ctx context.Context, input *apigateway.GetResourcesInput, opts ...func(*apigateway.Options)) (*apigateway.GetResourcesOutput, error) {
			if *input.RestApiId != "test-api" {
				t.Errorf("unexpected api id: %s", *input.RestApiId)
			}
			return &apigateway.GetResourcesOutput{
				Items: []types.Resource{
					{
						Id:   aws.String("root-id"),
						Path: aws.String("/"),
					},
				},
			}, nil
		},
	}

	apiGatewayClient := &APIGateway{
		client: mockClient,
	}

	err := apiGatewayClient.CreateEndpoint("test-api", EndpointOptions{
		Path:   "hello",
		Method: "GET",
		Uri:    "http://example.com",
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
