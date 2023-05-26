package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type mockAurora struct {
	createDBClusterFn func(context.Context, *rds.CreateDBClusterInput, ...func(*rds.Options)) (*rds.CreateDBClusterOutput, error)
}

func (m *mockAurora) CreateDBCluster(ctx context.Context, params *rds.CreateDBClusterInput, optFns ...func(*rds.Options)) (*rds.CreateDBClusterOutput, error) {
	return m.createDBClusterFn(ctx, params, optFns...)
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

	_, err := aurora.CreateDBCluster("identifier", "databaseName", "username", "password")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if aurora.cluster == nil {
		t.Errorf("expected cluster to be set, got: %v", aurora.cluster)
	}

	if aurora.clusterName != "databaseName" {
		t.Errorf("expected clusterName to be %s, got: %v", "databaseName", aurora.clusterName)
	}

	if aurora.secret == nil {
		t.Errorf("expected secret to be set, got: %v", aurora.cluster)
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

	_, err := aurora.CreateDBCluster("identifier", "databaseName", "username", "password")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = aurora.ExecuteStatement("sql")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
