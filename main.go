package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsService "github.com/florianwoelki/uber-movement-speed/aws"
)

type IAMRoles struct {
	s3         *aws.CredentialsCache
	kinesis    *aws.CredentialsCache
	lambda     *aws.CredentialsCache
	dynamodb   *aws.CredentialsCache
	glue       *aws.CredentialsCache
	aurora     *aws.CredentialsCache
	apiGateway *aws.CredentialsCache
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

	apiGatewayCreds, err := iam.CreateRoleWithPolicy("apigatewayv2-role", "apigatewayv2")
	if err != nil {
		return nil, err
	}
	iamRoles.apiGateway = apiGatewayCreds

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
	cfg.Credentials = iamRoles.apiGateway
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
	scriptFileBytes, err := io.ReadAll(file)
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
	zipFileBytes, err := io.ReadAll(file)
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
	zipFileBytes, err = io.ReadAll(file)
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

	// Loads the lambda zip file.
	log.Println("Uploading `DynamoGetter` lambda zip file to S3 bucket...")
	file, err = os.Open("services/dynamo_getter/dynamo_getter.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	zipFileBytes, err = io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Uploads the lambda zip file to the S3 bucket.
	err = s3.PutObject("lambda-bucket", "dynamo_getter.zip", zipFileBytes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Uploaded `DynamoGetter` lambda zip file to S3 bucket")

	// Creates the lambda function.
	log.Println("Creating `DynamoGetter` lambda function...")
	dynamoGetterARN, err := lambda.CreateGo("DynamoGetter", "lambda-bucket", "dynamo_getter.zip")
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Wait for lambda function to be created.
	log.Println("Created `DynamoGetter` lambda function")

	// Creates the websocket API Gateway.
	log.Println("Creating websocket API Gateway...")
	apiName := "my-kinesis-api"
	websocketApiGatewayId, err := apiGateway.CreateWebSocketApi(apiName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created websocket API Gateway")

	// Creates the HTTP API Gateway.
	log.Println("Creating http API Gateway...")
	apiName = "dynamo-getter"
	httpApiGatewayId, err := apiGateway.CreateHTTPApi(apiName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created http API Gateway")

	// Creates the websocket API Gateway endpoint for the kinesis data forwarder lambda function.
	log.Println("Creating websocket API Gateway endpoint for `KinesisDataForwarder` lambda function...")
	err = apiGateway.CreateWebSocket(websocketApiGatewayId, awsService.EndpointOptions{
		Path:   "kinesis-data-forwarder",
		Method: "POST",
		Uri:    kinesisDataForwarderARN,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created websocket API Gateway endpoint for `KinesisDataForwarder` lambda function")

	// Creates the REST API Gateway endpoint for the dynamo getter lambda function.
	log.Println("Creating REST API Gateway endpoint for `DynamoGetter` lambda function...")
	err = apiGateway.CreateEndpoint(httpApiGatewayId, awsService.EndpointOptions{
		Path:   "/dynamo-getter",
		Method: "GET",
		Uri:    fmt.Sprintf("arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/%s/invocations", dynamoGetterARN),
		RequestParameters: map[string]string{
			"method.request.querystring.id": "true",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created REST API Gateway endpoint for `DynamoGetter` lambda function")

	log.Println("Deploying API Gateway...")
	err = apiGateway.Deploy(websocketApiGatewayId)
	if err != nil {
		log.Fatal(err)
	}
	err = apiGateway.Deploy(httpApiGatewayId)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Deployed websocket API Gateway with ID:", websocketApiGatewayId)
	log.Println("Deployed http API Gateway with ID:", httpApiGatewayId)

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
	err = dynamodb.CreateTable("street_segment_speeds")
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
	// Creates the table for the Aurora DB. Changing `street_segment_speeds` requires a change in
	// `services/glue/raw_data_etl.py` as well.
	_, err = aurora.ExecuteStatement("test", clusterARN, secretARN, "CREATE TABLE street_segment_speeds (id SERIAL PRIMARY KEY, year INT, month INT, day INT, hour INT, utc_timestamp VARCHAR(100), start_junction_id VARCHAR(200), end_junction_id VARCHAR(200), osm_way_id BIGINT, osm_start_node_id BIGINT, osm_end_node_id BIGINT, speed_mph_mean FLOAT, speed_mph_stddev FLOAT)")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created Aurora DB Cluster Endpoint")
}
