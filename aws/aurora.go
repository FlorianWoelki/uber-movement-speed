package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type auroraAPI interface {
	CreateDBCluster(ctx context.Context, params *rds.CreateDBClusterInput, optFns ...func(*rds.Options)) (*rds.CreateDBClusterOutput, error)
	DescribeDBClusters(ctx context.Context, params *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
}

type rdsDataAPI interface {
	ExecuteStatement(ctx context.Context, params *rdsdata.ExecuteStatementInput, optFns ...func(*rdsdata.Options)) (*rdsdata.ExecuteStatementOutput, error)
}

type secretsManagerAPI interface {
	CreateSecret(ctx context.Context, params *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error)
}

type Aurora struct {
	rdsClient            auroraAPI
	rdsDataClient        rdsDataAPI
	secretsManagerClient secretsManagerAPI
}

func NewAurora(config aws.Config) *Aurora {
	return &Aurora{
		rdsClient:            rds.NewFromConfig(config),
		rdsDataClient:        rdsdata.NewFromConfig(config),
		secretsManagerClient: secretsmanager.NewFromConfig(config),
	}
}

func (a *Aurora) CreateDBCluster(identifier, databaseName, username, password string) (*rds.CreateDBClusterOutput, *secretsmanager.CreateSecretOutput, error) {
	// Creates the database cluster.
	cluster, err := a.rdsClient.CreateDBCluster(context.TODO(), &rds.CreateDBClusterInput{
		DBClusterIdentifier: aws.String(identifier),
		Engine:              aws.String("aurora-postgresql"),
		DatabaseName:        aws.String(databaseName),
		SourceRegion:        aws.String("us-east-1"),
	})
	if err != nil {
		return nil, nil, err
	}

	// Creates the secret for the database cluster.
	secret, err := a.secretsManagerClient.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
		Name:         aws.String(username),
		SecretString: aws.String(password),
	})
	if err != nil {
		return nil, nil, err
	}

	return cluster, secret, nil
}

func (a *Aurora) GetDBCluster(clusterIdentifier string) (*types.DBCluster, error) {
	clusters, err := a.rdsClient.DescribeDBClusters(context.TODO(), &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(clusterIdentifier),
	})
	if err != nil {
		return nil, err
	}

	if len(clusters.DBClusters) == 0 {
		return nil, nil
	}

	return &clusters.DBClusters[0], nil
}

func (a *Aurora) ExecuteStatement(databaseName, clusterArn, secretArn, sql string) error {
	_, err := a.rdsDataClient.ExecuteStatement(context.TODO(), &rdsdata.ExecuteStatementInput{
		Database:              aws.String(databaseName),
		ResourceArn:           aws.String(clusterArn),
		SecretArn:             aws.String(secretArn),
		IncludeResultMetadata: true,
		Sql:                   aws.String(sql),
	})
	if err != nil {
		return err
	}

	return nil
}
