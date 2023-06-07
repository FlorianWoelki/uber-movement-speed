package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
)

type apiGatewayAPI interface {
	CreateApi(ctx context.Context, params *apigatewayv2.CreateApiInput, optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateApiOutput, error)
	DeleteApi(ctx context.Context, params *apigatewayv2.DeleteApiInput, optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.DeleteApiOutput, error)
	CreateDeployment(ctx context.Context, params *apigatewayv2.CreateDeploymentInput, optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateDeploymentOutput, error)
	CreateRoute(ctx context.Context, params *apigatewayv2.CreateRouteInput, optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateRouteOutput, error)
	CreateIntegration(ctx context.Context, params *apigatewayv2.CreateIntegrationInput, optFns ...func(*apigatewayv2.Options)) (*apigatewayv2.CreateIntegrationOutput, error)
}

// APIGateway is a wrapper around the AWS API Gateway client.
type APIGateway struct {
	client apiGatewayAPI
}

// NewAPIGateway creates a new API Gateway client with the given configuration.
func NewAPIGateway(config aws.Config) *APIGateway {
	return &APIGateway{
		client: apigatewayv2.NewFromConfig(config),
	}
}

// Create creates an API Gateway with the given name and returns the ID of the API Gateway
// that was created.
func (a *APIGateway) Create(name string) (string, error) {
	createOutput, err := a.client.CreateApi(context.TODO(), &apigatewayv2.CreateApiInput{
		Name:                     aws.String(name),
		ProtocolType:             types.ProtocolTypeWebsocket,
		RouteSelectionExpression: aws.String("$request.body.action"),
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(createOutput.ApiId), nil
}

// Delete deletes the API Gateway with the given ID.
func (a *APIGateway) Delete(id string) error {
	_, err := a.client.DeleteApi(context.TODO(), &apigatewayv2.DeleteApiInput{
		ApiId: aws.String(id),
	})
	if err != nil {
		return err
	}

	return nil
}

// EndpointOptions are the options for creating an endpoint.
type EndpointOptions struct {
	// Path is the path of the resource where the endpoint will be created.
	Path string
	// Method is the HTTP method of the endpoint.
	Method string
	// Uri is the URI of the endpoint.
	Uri string
	// RequestParameters are the request parameters of the endpoint.
	RequestParameters map[string]types.ParameterConstraints
}

// CreateEndpoint creates an endpoint for the given API Gateway ID with the given options.
// It creates a resource with the given path, a method with the given HTTP method,
// an integration with the given URI, and a deployment.
func (a *APIGateway) CreateEndpoint(id string, options EndpointOptions) error {
	integrationOutput, err := a.client.CreateIntegration(context.TODO(), &apigatewayv2.CreateIntegrationInput{
		ApiId:             aws.String(id),
		IntegrationType:   types.IntegrationTypeAwsProxy,
		IntegrationMethod: aws.String(options.Method),
		IntegrationUri:    aws.String(options.Uri),
		// PayloadFormatVersion: aws.String("1.0"),
		// PassthroughBehavior:  types.PassthroughBehaviorWhenNoMatch,
	})
	if err != nil {
		return err
	}

	_, err = a.client.CreateRoute(context.TODO(), &apigatewayv2.CreateRouteInput{
		ApiId:                            aws.String(id),
		RouteKey:                         aws.String(options.Path),
		Target:                           integrationOutput.IntegrationId,
		RouteResponseSelectionExpression: aws.String("$default"),
		// RouteKey:          aws.String(fmt.Sprintf("%s %s", options.Method, options.Path)),
		// RequestParameters: options.RequestParameters,
		// AuthorizationType: types.AuthorizationTypeNone,
	})
	if err != nil {
		return err
	}

	// Create a route for the `$connect` websocket event which is used for opening the
	// connection.
	_, err = a.client.CreateRoute(context.TODO(), &apigatewayv2.CreateRouteInput{
		ApiId:    aws.String(id),
		RouteKey: aws.String("$connect"),
		Target:   integrationOutput.IntegrationId,
	})
	if err != nil {
		return err
	}

	// Create a route for the `$disconnect` websocket event which is used for closing the
	// connection.
	_, err = a.client.CreateRoute(context.TODO(), &apigatewayv2.CreateRouteInput{
		ApiId:    aws.String(id),
		RouteKey: aws.String("$disconnect"),
		Target:   integrationOutput.IntegrationId,
	})
	if err != nil {
		return err
	}

	// Create a route for the `$default` websocket event which is used for data transfer.
	_, err = a.client.CreateRoute(context.TODO(), &apigatewayv2.CreateRouteInput{
		ApiId:                            aws.String(id),
		RouteKey:                         aws.String("$default"),
		Target:                           integrationOutput.IntegrationId,
		RouteResponseSelectionExpression: aws.String("$default"),
	})
	if err != nil {
		return err
	}

	_, err = a.client.CreateDeployment(context.TODO(), &apigatewayv2.CreateDeploymentInput{
		ApiId:     aws.String(id),
		StageName: aws.String("dev"),
	})
	if err != nil {
		return err
	}

	return nil
}
