package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsService "github.com/florianwoelki/uber-movement-speed/aws"
)

func main() {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           "http://s3.localhost.localstack.cloud:4566",
				SigningRegion: "us-east-1",
			}, nil
		}

		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           "http://localhost:4566",
			SigningRegion: "us-east-1",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		log.Fatal(err)
	}

	kinesis := awsService.NewKinesis(cfg)
	lambda := awsService.NewLambda(cfg)
	s3 := awsService.NewS3(cfg)
	apiGateway := awsService.NewAPIGateway(cfg)

	// Creates the lambda S3 bucket.
	err = s3.CreateBucket("lambda-bucket")
	if err != nil {
		log.Fatal(err)
	}

	// Loads the lambda zip file.
	file, err := os.Open("services/preprocessing/preprocessing.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	zipFileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Uploads the lambda zip file to the S3 bucket.
	err = s3.PutObject("lambda-bucket", "preprocessing.zip", zipFileBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Creates the lambda function.
	kinesisDataForwarderARN, err := lambda.CreateGo("Preprocessing", "lambda-bucket", "preprocessing.zip")
	if err != nil {
		log.Fatal(err)
	}

	// Creates the kinesis stream.
	streamName := "my-kinesis-stream"
	err = kinesis.Create(streamName)
	if err != nil {
		log.Fatal(err)
	}

	// Gets the ARN from the kinesis stream.
	streamARN, err := kinesis.GetARN(streamName)
	if err != nil {
		log.Fatal(err)
	}

	// Creates the lambda event source mapping.
	err = lambda.BindToService("KinesisDataForwarder", streamARN)
	if err != nil {
		log.Fatal(err)
	}

	// Creates the API Gateway.
	apiName := "my-kinesis-api"
	apiGatewayId, err := apiGateway.Create(apiName)
	if err != nil {
		log.Fatal(err)
	}

	// Creates the API Gateway endpoint for the kinesis data forwarder lambda function.
	err = apiGateway.CreateEndpoint(apiGatewayId, awsService.EndpointOptions{
		Path:            "kinesis",
		Method:          "POST",
		Uri:             fmt.Sprintf("arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/%s/invocations", kinesisDataForwarderARN),
		IntegrationType: "AWS_PROXY",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(apiGatewayId)

	// aurora := awsService.NewAurora(cfg)

	// _, secret, err := aurora.CreateDBCluster("db1", "test", "test", "test")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(aws.ToString(secret.ARN))

	// cluster, err := aurora.GetDBCluster("db1")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// clusterArn := aws.ToString(cluster.DBClusterArn)
	// secretArn := "arn:aws:secretsmanager:us-east-1:000000000000:secret:test-abQbyQ"
	// err = aurora.ExecuteStatement("test", clusterArn, secretArn, "SELECT 123")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// apiGatewayClient := apigateway.NewFromConfig(cfg)
	// s3client := s3.NewFromConfig(cfg)
	// lambdaClient := lambda.NewFromConfig(cfg)

	// _, err = s3client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
	// 	Bucket: aws.String("my-bucket"),
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// zipFilePath := "my-code.zip"
	// zipFileKey := "my-code.zip"
	// file, err := os.Open(zipFilePath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// zipFileBytes, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = s3client.PutObject(context.TODO(), &s3.PutObjectInput{
	// 	Bucket: aws.String("my-bucket"),
	// 	Key:    aws.String(zipFileKey),
	// 	Body:   bytes.NewReader(zipFileBytes),
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// createFunctionOutput, err := lambdaClient.CreateFunction(context.TODO(), &lambda.CreateFunctionInput{
	// 	Code: &lambdaTypes.FunctionCode{
	// 		S3Bucket: aws.String("my-bucket"),
	// 		S3Key:    aws.String("my-code.zip"),
	// 	},
	// 	FunctionName: aws.String("MyLambdaFunction"),
	// 	Handler:      aws.String("main"),
	// 	Runtime:      lambdaTypes.RuntimeGo1x,
	// 	Role:         aws.String("arn:aws:iam::123456789012:role/lambda-role"),
	// 	Timeout:      aws.Int32(60),
	// 	MemorySize:   aws.Int32(128),
	// 	Publish:      true,
	// 	Environment:  &lambdaTypes.Environment{},
	// })

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// createAPIOutput, err := apiGatewayClient.CreateRestApi(context.TODO(), &apigateway.CreateRestApiInput{
	// 	Name: aws.String("MyApi"),
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// getResourcesOutput, err := apiGatewayClient.GetResources(context.TODO(), &apigateway.GetResourcesInput{
	// 	RestApiId: createAPIOutput.Id,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var parentId string
	// for _, resource := range getResourcesOutput.Items {
	// 	if aws.ToString(resource.Path) == "/" {
	// 		parentId = aws.ToString(resource.Id)
	// 		break
	// 	}
	// }

	// if parentId == "" {
	// 	log.Fatal("Root resource not found")
	// }

	// createResourceOutput, err := apiGatewayClient.CreateResource(context.TODO(), &apigateway.CreateResourceInput{
	// 	RestApiId: createAPIOutput.Id,
	// 	ParentId:  aws.String(parentId),
	// 	PathPart:  aws.String("hello"),
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// resourceId := aws.ToString(createResourceOutput.Id)

	// _, err = apiGatewayClient.PutMethod(context.TODO(), &apigateway.PutMethodInput{
	// 	RestApiId:         createAPIOutput.Id,
	// 	ResourceId:        aws.String(resourceId),
	// 	HttpMethod:        aws.String("GET"),
	// 	RequestParameters: map[string]bool{"method.request.path.hello": true},
	// 	AuthorizationType: aws.String("NONE"),
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = apiGatewayClient.PutIntegration(context.TODO(), &apigateway.PutIntegrationInput{
	// 	RestApiId:             createAPIOutput.Id,
	// 	ResourceId:            aws.String(resourceId),
	// 	HttpMethod:            aws.String("GET"),
	// 	Type:                  apiGatewayTypes.IntegrationTypeAwsProxy,
	// 	IntegrationHttpMethod: aws.String("POST"),
	// 	Uri:                   aws.String(fmt.Sprintf("arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/%s/invocations", *createFunctionOutput.FunctionArn)),
	// 	PassthroughBehavior:   aws.String("WHEN_NO_MATCH"),
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = apiGatewayClient.CreateDeployment(context.TODO(), &apigateway.CreateDeploymentInput{
	// 	RestApiId: createAPIOutput.Id,
	// 	StageName: aws.String("dev"),
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(aws.ToString(createAPIOutput.Id))
}
