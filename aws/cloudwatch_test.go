package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

type mockCloudWatchClient struct {
	putMetricAlarmFunc func(ctx context.Context, params *cloudwatch.PutMetricAlarmInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricAlarmOutput, error)
	putMetricDataFunc  func(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error)
}

func (m *mockCloudWatchClient) PutMetricAlarm(ctx context.Context, params *cloudwatch.PutMetricAlarmInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricAlarmOutput, error) {
	return m.putMetricAlarmFunc(ctx, params, optFns...)
}

func (m *mockCloudWatchClient) PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error) {
	return m.putMetricDataFunc(ctx, params, optFns...)
}

func TestCloudWatch_PutMetricAlarm(t *testing.T) {
	mockClient := &mockCloudWatchClient{
		putMetricAlarmFunc: func(ctx context.Context, params *cloudwatch.PutMetricAlarmInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricAlarmOutput, error) {
			return &cloudwatch.PutMetricAlarmOutput{}, nil
		},
	}

	cw := &CloudWatch{
		client: mockClient,
	}

	err := cw.PutMetricAlarm("test", "test", "test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCloudWatch_PutMetricData(t *testing.T) {
	mockClient := &mockCloudWatchClient{
		putMetricDataFunc: func(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error) {
			return &cloudwatch.PutMetricDataOutput{}, nil
		},
	}

	cw := &CloudWatch{
		client: mockClient,
	}

	err := cw.PutMetricData("test", "test", 1.0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
