package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type mockAurora struct {
	createDBClusterFn    func(context.Context, *rds.CreateDBClusterInput, ...func(*rds.Options)) (*rds.CreateDBClusterOutput, error)
	describeDBClustersFn func(context.Context, *rds.DescribeDBClustersInput, ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
}

func (m *mockAurora) CreateDBCluster(ctx context.Context, params *rds.CreateDBClusterInput, optFns ...func(*rds.Options)) (*rds.CreateDBClusterOutput, error) {
	return m.createDBClusterFn(ctx, params, optFns...)
}

func (m *mockAurora) DescribeDBClusters(ctx context.Context, params *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error) {
	return m.describeDBClustersFn(ctx, params, optFns...)
}

type mockRDSData struct {
	executeStatementFn func(context.Context, *rdsdata.ExecuteStatementInput, ...func(*rdsdata.Options)) (*rdsdata.ExecuteStatementOutput, error)
}

func (m *mockRDSData) ExecuteStatement(ctx context.Context, params *rdsdata.ExecuteStatementInput, optFns ...func(*rdsdata.Options)) (*rdsdata.ExecuteStatementOutput, error) {
	return m.executeStatementFn(ctx, params, optFns...)
}

type mockSecretsManager struct {
	createSecretFn func(context.Context, *secretsmanager.CreateSecretInput, ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error)
}

func (m *mockSecretsManager) CreateSecret(ctx context.Context, params *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error) {
	return m.createSecretFn(ctx, params, optFns...)
}

func TestAurora_CreateDBCluster(t *testing.T) {
	mockClient := &mockAurora{
		createDBClusterFn: func(ctx context.Context, params *rds.CreateDBClusterInput, optFns ...func(*rds.Options)) (*rds.CreateDBClusterOutput, error) {
			return &rds.CreateDBClusterOutput{}, nil
		},
	}

	mockRDSDataClient := &mockRDSData{
		executeStatementFn: func(ctx context.Context, params *rdsdata.ExecuteStatementInput, optFns ...func(*rdsdata.Options)) (*rdsdata.ExecuteStatementOutput, error) {
			return &rdsdata.ExecuteStatementOutput{}, nil
		},
	}

	mockSecretsManager := &mockSecretsManager{
		createSecretFn: func(ctx context.Context, params *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error) {
			return &secretsmanager.CreateSecretOutput{}, nil
		},
	}

	aurora := &Aurora{
		rdsClient:            mockClient,
		rdsDataClient:        mockRDSDataClient,
		secretsManagerClient: mockSecretsManager,
	}

	_, _, err := aurora.CreateDBCluster("identifier", "databaseName", "username", "password")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAurora_ExecuteStatement(t *testing.T) {
	mockClient := &mockAurora{
		createDBClusterFn: func(ctx context.Context, params *rds.CreateDBClusterInput, optFns ...func(*rds.Options)) (*rds.CreateDBClusterOutput, error) {
			return &rds.CreateDBClusterOutput{}, nil
		},
	}

	mockRDSDataClient := &mockRDSData{
		executeStatementFn: func(ctx context.Context, params *rdsdata.ExecuteStatementInput, optFns ...func(*rdsdata.Options)) (*rdsdata.ExecuteStatementOutput, error) {
			return &rdsdata.ExecuteStatementOutput{}, nil
		},
	}

	mockSecretsManager := &mockSecretsManager{
		createSecretFn: func(ctx context.Context, params *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error) {
			return &secretsmanager.CreateSecretOutput{}, nil
		},
	}

	aurora := &Aurora{
		rdsClient:            mockClient,
		rdsDataClient:        mockRDSDataClient,
		secretsManagerClient: mockSecretsManager,
	}

	_, _, err := aurora.CreateDBCluster("identifier", "databaseName", "username", "password")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = aurora.ExecuteStatement("databaseName", "clusterArn", "secretArn", "sql")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAurora_GetDBCluster(t *testing.T) {
	mockClient := &mockAurora{
		createDBClusterFn: func(ctx context.Context, params *rds.CreateDBClusterInput, optFns ...func(*rds.Options)) (*rds.CreateDBClusterOutput, error) {
			return &rds.CreateDBClusterOutput{}, nil
		},
		describeDBClustersFn: func(ctx context.Context, params *rds.DescribeDBClustersInput, optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error) {
			return &rds.DescribeDBClustersOutput{}, nil
		},
	}

	mockRDSDataClient := &mockRDSData{
		executeStatementFn: func(ctx context.Context, params *rdsdata.ExecuteStatementInput, optFns ...func(*rdsdata.Options)) (*rdsdata.ExecuteStatementOutput, error) {
			return &rdsdata.ExecuteStatementOutput{}, nil
		},
	}

	mockSecretsManager := &mockSecretsManager{
		createSecretFn: func(ctx context.Context, params *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error) {
			return &secretsmanager.CreateSecretOutput{}, nil
		},
	}

	aurora := &Aurora{
		rdsClient:            mockClient,
		rdsDataClient:        mockRDSDataClient,
		secretsManagerClient: mockSecretsManager,
	}

	_, _, err := aurora.CreateDBCluster("identifier", "databaseName", "username", "password")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = aurora.GetDBCluster("identifier")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
