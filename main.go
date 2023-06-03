package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsService "github.com/florianwoelki/uber-movement-speed/aws"
)

type IAMRoles struct {
	s3       *aws.CredentialsCache
	kinesis  *aws.CredentialsCache
	lambda   *aws.CredentialsCache
	dynamodb *aws.CredentialsCache
	glue     *aws.CredentialsCache
	aurora   *aws.CredentialsCache
}

func createIAMRoles(iam *awsService.IAM) (*IAMRoles, error) {
	iamRoles := &IAMRoles{}

	s3Creds, err := iam.CreateRoleWithPolicy("s3-role", "s3")
	if err != nil {
		return nil, err
	}
	iamRoles.s3 = s3Creds

	kinesisCreds, err := iam.CreateRoleWithPolicy("kinesis-role", "kinesis")
	if err != nil {
		return nil, err
	}
	iamRoles.kinesis = kinesisCreds

	lambdaCreds, err := iam.CreateRoleWithPolicy("lambda-role", "lambda")
	if err != nil {
		return nil, err
	}
	iamRoles.lambda = lambdaCreds

	dynamodbCreds, err := iam.CreateRoleWithPolicy("dynamodb-role", "dynamodb")
	if err != nil {
		return nil, err
	}
	iamRoles.dynamodb = dynamodbCreds

	glueCreds, err := iam.CreateRoleWithPolicy("glue-role", "glue")
	if err != nil {
		return nil, err
	}
	iamRoles.glue = glueCreds

	auroraCreds, err := iam.CreateRoleWithPolicy("rds-role", "rds")
	if err != nil {
		return nil, err
	}
	iamRoles.aurora = auroraCreds

	return iamRoles, nil
}

func main() {
	log.Println("Starting setup...")
	defer log.Println("Finished setup")

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

	log.Println("Creating IAM roles...")
	iam := awsService.NewIAM(cfg)
	iamRoles, err := createIAMRoles(iam)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created IAM roles")

	cfg.Credentials = iamRoles.kinesis
	kinesis := awsService.NewKinesis(cfg)
	cfg.Credentials = iamRoles.lambda
	lambda := awsService.NewLambda(cfg)
	cfg.Credentials = iamRoles.s3
	s3 := awsService.NewS3(cfg)
	cfg.Credentials = iamRoles.dynamodb
	dynamodb := awsService.NewDynamoDB(cfg)
	cfg.Credentials = iamRoles.glue
	glue := awsService.NewGlue(cfg)
	cfg.Credentials = iamRoles.aurora
	aurora := awsService.NewAurora(cfg)
	apiGateway := awsService.NewAPIGateway(cfg)

	// Creates the lambda S3 bucket.
	log.Println("Creating S3 bucket for lambda...")
	err = s3.CreateBucket("lambda-bucket")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created S3 bucket for lambda")

	// Creates the S3 bucket for the raw data that is being sent from the `preprocessing`
	// service.
	log.Println("Creating S3 bucket for raw data...")
	err = s3.CreateBucket("raw-data")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created S3 bucket for raw data")

	// Loads the python PySpark script.
	log.Println("Uploading PySpark script to `raw-data` S3 bucket...")
	file, err := os.Open("services/glue/raw_data_etl.py")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scriptFileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Uploads the python PySpark script to the S3 bucket.
	err = s3.PutObject("raw-data", "scripts/raw_data_etl.py", scriptFileBytes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Uploaded PySpark script to `raw-data` S3 bucket")

	// Creates the S3 bucket for the transformed data that is being sent from the glue job
	log.Println("Creating S3 bucket for transformed data...")
	err = s3.CreateBucket("transformed-data")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created S3 bucket for transformed data")

	// Creates the glue job.
	log.Println("Creating glue job...")
	err = glue.CreateJob("raw-data-etl", "s3://raw-data/scripts/raw_data_etl.py")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created glue job")

	// Loads the lambda zip file.
	log.Println("Uploading `Preprocessing` lambda zip file to S3 bucket...")
	file, err = os.Open("services/preprocessing/preprocessing.zip")
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
	log.Println("Uploaded `Preprocessing` lambda zip file to S3 bucket")

	// Creates the lambda function.
	log.Println("Creating `Preprocessing` lambda function...")
	_, err = lambda.CreateGo("Preprocessing", "lambda-bucket", "preprocessing.zip")
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Wait for lambda function to be created.
	log.Println("Created `Preprocessing` lambda function")

	// Loads the lambda zip file.
	log.Println("Uploading `KinesisDataForwarder` lambda zip file to S3 bucket...")
	file, err = os.Open("services/kinesis_data_forwarder/dist/kinesis_data_forwarder.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	zipFileBytes, err = ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Uploads the lambda zip file to the S3 bucket.
	err = s3.PutObject("lambda-bucket", "kinesis_data_forwarder.zip", zipFileBytes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Uploaded `KinesisDataForwarder` lambda zip file to S3 bucket")

	// Creates the lambda function.
	log.Println("Creating `KinesisDataForwarder` lambda function...")
	kinesisDataForwarderARN, err := lambda.CreateNode("KinesisDataForwarder", "lambda-bucket", "kinesis_data_forwarder.zip")
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Wait for lambda function to be created.
	log.Println("Created `KinesisDataForwarder` lambda function")

	// Creates the API Gateway.
	log.Println("Creating API Gateway...")
	apiName := "my-kinesis-api"
	apiGatewayId, err := apiGateway.Create(apiName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created API Gateway")

	// Creates the API Gateway endpoint for the kinesis data forwarder lambda function.
	log.Println("Creating API Gateway endpoint for `KinesisDataForwarder` lambda function...")
	err = apiGateway.CreateEndpoint(apiGatewayId, awsService.EndpointOptions{
		Path:   "/kinesis",
		Method: "POST",
		Uri:    kinesisDataForwarderARN,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created API Gateway endpoint for `KinesisDataForwarder` lambda function")
	log.Println("API Gateway ID:", apiGatewayId)

	// Creates the kinesis stream.
	log.Println("Creating kinesis stream...")
	streamName := "my-kinesis-stream"
	err = kinesis.Create(streamName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created kinesis stream")

	// Gets the ARN from the kinesis stream.
	log.Println("Binding lambda function to kinesis stream...")
	streamARN, err := kinesis.GetARN(streamName)
	if err != nil {
		log.Fatal(err)
	}

	// Creates the lambda event source mapping.
	err = lambda.BindToService("Preprocessing", streamARN)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Bound lambda function to kinesis stream")

	// Create dynamodb table.
	log.Println("Creating dynamodb table...")
	err = dynamodb.CreateTable("books")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created dynamodb table")

	// Creates the Aurora DB Cluster.
	log.Println("Creating Aurora DB Cluster...")
	// Changing `dbpass`, `db1`, or `test` requires a change in
	// `services/glue/raw_data_etl.py` as well.
	cluster, secret, err := aurora.CreateDBCluster("db1", "test", "dbpass", "test")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created Aurora DB Cluster")

	clusterARN := aws.ToString(cluster.DBCluster.DBClusterArn)
	secretARN := aws.ToString(secret.ARN)

	status := cluster.DBCluster.Status
	for aws.ToString(status) != "available" {
		log.Println("Waiting for Aurora DB Cluster to be available...")
		c, err := aurora.GetDBCluster("db1")
		if err != nil {
			log.Fatal(err)
		}

		status = c.Status
		time.Sleep(2 * time.Second)
	}
	log.Println("Aurora DB Cluster is available")

	log.Println("Creating Aurora DB Cluster Endpoint...")
	// Creates the table for the Aurora DB. Changing `books` requires a change in
	// `services/glue/raw_data_etl.py` as well.
	_, err = aurora.ExecuteStatement("test", clusterARN, secretARN, "CREATE TABLE books (id SERIAL PRIMARY KEY, title VARCHAR(100))")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created Aurora DB Cluster Endpoint")

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
