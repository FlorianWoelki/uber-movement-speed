package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type auroraAPI interface {
	CreateDBCluster(ctx context.Context, params *rds.CreateDBClusterInput, optFns ...func(*rds.Options)) (*rds.CreateDBClusterOutput, error)
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
	cluster              *rds.CreateDBClusterOutput
	clusterName          string
	secret               *secretsmanager.CreateSecretOutput
}

func NewAurora(config aws.Config) *Aurora {
	return &Aurora{
		rdsClient:            rds.NewFromConfig(config),
		rdsDataClient:        rdsdata.NewFromConfig(config),
		secretsManagerClient: secretsmanager.NewFromConfig(config),
	}
}

func (a *Aurora) CreateDBCluster(identifier, databaseName, username, password string) (*rds.CreateDBClusterOutput, error) {
	// Creates the database cluster.
	cluster, err := a.rdsClient.CreateDBCluster(context.TODO(), &rds.CreateDBClusterInput{
		DBClusterIdentifier: aws.String(identifier),
		Engine:              aws.String("aurora-postgresql"),
		DatabaseName:        aws.String(databaseName),
	})
	if err != nil {
		return nil, err
	}
	a.cluster = cluster
	a.clusterName = databaseName

	// Creates the secret for the database cluster.
	secret, err := a.secretsManagerClient.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
		Name:         aws.String(username),
		SecretString: aws.String(password),
	})
	if err != nil {
		return nil, err
	}
	a.secret = secret

	fmt.Println(cluster.DBCluster)
	fmt.Println(secret)
	return cluster, nil
}

func (a *Aurora) ExecuteStatement(sql string) error {
	_, err := a.rdsDataClient.ExecuteStatement(context.TODO(), &rdsdata.ExecuteStatementInput{
		Database:              aws.String(a.clusterName),
		ResourceArn:           a.cluster.DBCluster.DBClusterArn,
		SecretArn:             a.secret.ARN,
		IncludeResultMetadata: true,
		Sql:                   aws.String(sql),
	})
	if err != nil {
		return err
	}

	return nil
}
