package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
)

type mockAPIGatewayClient struct {
	createApiFunc         func(ctx context.Context, input *apigatewayv2.CreateApiInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateApiOutput, error)
	deleteApiFunc         func(ctx context.Context, input *apigatewayv2.DeleteApiInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.DeleteApiOutput, error)
	createDeploymentFunc  func(ctx context.Context, input *apigatewayv2.CreateDeploymentInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateDeploymentOutput, error)
	createRouteFunc       func(ctx context.Context, input *apigatewayv2.CreateRouteInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateRouteOutput, error)
	createIntegrationFunc func(ctx context.Context, input *apigatewayv2.CreateIntegrationInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateIntegrationOutput, error)
}

func (m *mockAPIGatewayClient) CreateApi(ctx context.Context, input *apigatewayv2.CreateApiInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateApiOutput, error) {
	return m.createApiFunc(ctx, input, opts...)
}

func (m *mockAPIGatewayClient) DeleteApi(ctx context.Context, input *apigatewayv2.DeleteApiInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.DeleteApiOutput, error) {
	return m.deleteApiFunc(ctx, input, opts...)
}

func (m *mockAPIGatewayClient) CreateDeployment(ctx context.Context, input *apigatewayv2.CreateDeploymentInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateDeploymentOutput, error) {
	return m.createDeploymentFunc(ctx, input, opts...)
}

func (m *mockAPIGatewayClient) CreateRoute(ctx context.Context, input *apigatewayv2.CreateRouteInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateRouteOutput, error) {
	return m.createRouteFunc(ctx, input, opts...)
}

func (m *mockAPIGatewayClient) CreateIntegration(ctx context.Context, input *apigatewayv2.CreateIntegrationInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateIntegrationOutput, error) {
	return m.createIntegrationFunc(ctx, input, opts...)
}

func TestAPIGateway_CreateWebSocketApi(t *testing.T) {
	mockClient := &mockAPIGatewayClient{
		createApiFunc: func(ctx context.Context, input *apigatewayv2.CreateApiInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateApiOutput, error) {
			if aws.ToString(input.Name) != "test-api" {
				t.Errorf("unexpected api name: %s", aws.ToString(input.Name))
			}
			return &apigatewayv2.CreateApiOutput{
				ApiId: aws.String("test-api"),
			}, nil
		},
	}

	apiGatewayClient := &APIGateway{
		client: mockClient,
	}

	id, err := apiGatewayClient.CreateWebSocketApi("test-api")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if id != "test-api" {
		t.Errorf("unexpected api id: %s", id)
	}
}

func TestAPIGateway_CreateHTTPApi(t *testing.T) {
	mockClient := &mockAPIGatewayClient{
		createApiFunc: func(ctx context.Context, input *apigatewayv2.CreateApiInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateApiOutput, error) {
			if aws.ToString(input.Name) != "test-api" {
				t.Errorf("unexpected api name: %s", aws.ToString(input.Name))
			}
			return &apigatewayv2.CreateApiOutput{
				ApiId: aws.String("test-api"),
			}, nil
		},
	}

	apiGatewayClient := &APIGateway{
		client: mockClient,
	}

	id, err := apiGatewayClient.CreateHTTPApi("test-api")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if id != "test-api" {
		t.Errorf("unexpected api id: %s", id)
	}
}

func TestAPIGateway_Delete(t *testing.T) {
	mockClient := &mockAPIGatewayClient{
		deleteApiFunc: func(ctx context.Context, input *apigatewayv2.DeleteApiInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.DeleteApiOutput, error) {
			if aws.ToString(input.ApiId) != "test-api" {
				t.Errorf("unexpected api id: %s", aws.ToString(input.ApiId))
			}
			return &apigatewayv2.DeleteApiOutput{}, nil
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

func TestAPIGateway_CreateEndpoint(t *testing.T) {
	mockClient := &mockAPIGatewayClient{
		createDeploymentFunc: func(ctx context.Context, input *apigatewayv2.CreateDeploymentInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateDeploymentOutput, error) {
			if aws.ToString(input.ApiId) != "test-api" {
				t.Errorf("unexpected api id: %s", aws.ToString(input.ApiId))
			}
			return &apigatewayv2.CreateDeploymentOutput{}, nil
		},
		createRouteFunc: func(ctx context.Context, input *apigatewayv2.CreateRouteInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateRouteOutput, error) {
			if aws.ToString(input.ApiId) != "test-api" {
				t.Errorf("unexpected api id: %s", aws.ToString(input.ApiId))
			}
			if aws.ToString(input.RouteKey) != "GET /hello" {
				t.Errorf("unexpected route key: %s", aws.ToString(input.RouteKey))
			}
			return &apigatewayv2.CreateRouteOutput{
				RouteId: aws.String("test-route"),
			}, nil
		},
		createIntegrationFunc: func(ctx context.Context, input *apigatewayv2.CreateIntegrationInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateIntegrationOutput, error) {
			if aws.ToString(input.ApiId) != "test-api" {
				t.Errorf("unexpected api id: %s", aws.ToString(input.ApiId))
			}
			if input.IntegrationType != "AWS_PROXY" {
				t.Errorf("unexpected integration type: %s", input.IntegrationType)
			}
			if aws.ToString(input.IntegrationUri) != "http://example.com" {
				t.Errorf("unexpected integration uri: %s", aws.ToString(input.IntegrationUri))
			}
			if aws.ToString(input.IntegrationMethod) != "GET" {
				t.Errorf("unexpected integration method: %s", aws.ToString(input.IntegrationMethod))
			}
			return &apigatewayv2.CreateIntegrationOutput{}, nil
		},
	}

	apiGatewayClient := &APIGateway{
		client: mockClient,
	}

	err := apiGatewayClient.CreateEndpoint("test-api", EndpointOptions{
		Path:   "/hello",
		Method: "GET",
		Uri:    "http://example.com",
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAPIGateway_CreateWebSocket(t *testing.T) {
	mockClient := &mockAPIGatewayClient{
		createDeploymentFunc: func(ctx context.Context, input *apigatewayv2.CreateDeploymentInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateDeploymentOutput, error) {
			if aws.ToString(input.ApiId) != "test-api" {
				t.Errorf("unexpected api id: %s", aws.ToString(input.ApiId))
			}
			return &apigatewayv2.CreateDeploymentOutput{}, nil
		},
		createRouteFunc: func(ctx context.Context, input *apigatewayv2.CreateRouteInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateRouteOutput, error) {
			if aws.ToString(input.ApiId) != "test-api" {
				t.Errorf("unexpected api id: %s", aws.ToString(input.ApiId))
			}
			if aws.ToString(input.RouteKey) != "$connect" && aws.ToString(input.RouteKey) != "$disconnect" && aws.ToString(input.RouteKey) != "$default" && aws.ToString(input.RouteKey) != "hello" {
				t.Errorf("unexpected route key: %s", aws.ToString(input.RouteKey))
			}
			return &apigatewayv2.CreateRouteOutput{
				RouteId: aws.String("test-route"),
			}, nil
		},
		createIntegrationFunc: func(ctx context.Context, input *apigatewayv2.CreateIntegrationInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateIntegrationOutput, error) {
			if aws.ToString(input.ApiId) != "test-api" {
				t.Errorf("unexpected api id: %s", aws.ToString(input.ApiId))
			}
			if input.IntegrationType != "AWS_PROXY" {
				t.Errorf("unexpected integration type: %s", input.IntegrationType)
			}
			if aws.ToString(input.IntegrationUri) != "http://example.com" {
				t.Errorf("unexpected integration uri: %s", aws.ToString(input.IntegrationUri))
			}
			if aws.ToString(input.IntegrationMethod) != "POST" {
				t.Errorf("unexpected integration method: %s", aws.ToString(input.IntegrationMethod))
			}
			return &apigatewayv2.CreateIntegrationOutput{}, nil
		},
	}

	apiGatewayClient := &APIGateway{
		client: mockClient,
	}

	err := apiGatewayClient.CreateWebSocket("test-api", EndpointOptions{
		Path:   "hello",
		Method: "POST",
		Uri:    "http://example.com",
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAPIGateway_CreateDeployment(t *testing.T) {
	mockClient := &mockAPIGatewayClient{
		createDeploymentFunc: func(ctx context.Context, input *apigatewayv2.CreateDeploymentInput, opts ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateDeploymentOutput, error) {
			if aws.ToString(input.ApiId) != "test-api" {
				t.Errorf("unexpected api id: %s", aws.ToString(input.ApiId))
			}
			return &apigatewayv2.CreateDeploymentOutput{}, nil
		},
	}

	apiGatewayClient := &APIGateway{
		client: mockClient,
	}

	err := apiGatewayClient.Deploy("test-api")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
