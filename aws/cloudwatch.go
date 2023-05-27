package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type cloudWatchAPI interface {
	PutMetricAlarm(ctx context.Context, params *cloudwatch.PutMetricAlarmInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricAlarmOutput, error)
	PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error)
}

// CloudWatch is a wrapper around the AWS CloudWatch client.
type CloudWatch struct {
	client cloudWatchAPI
}

// NewCloudWatch creates a new CloudWatch client with the given configuration.
func NewCloudWatch(config aws.Config) *CloudWatch {
	return &CloudWatch{
		client: cloudwatch.NewFromConfig(config),
	}
}

// PutMetricAlarm creates a CloudWatch alarm with the given name, metric name and namespace.
// The alarm is set to evaluate the metric every `30` seconds and the alarm will be
// triggered if the metric value is less than `1.0` for `1` evaluation period.
func (c *CloudWatch) PutMetricAlarm(alarmName, metricName, namespace string) error {
	_, err := c.client.PutMetricAlarm(context.TODO(), &cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(alarmName),
		MetricName:         aws.String(metricName),
		Namespace:          aws.String(namespace),
		Threshold:          aws.Float64(1.0),
		ComparisonOperator: types.ComparisonOperatorLessThanThreshold,
		EvaluationPeriods:  aws.Int32(1),
		Period:             aws.Int32(30),
		Statistic:          types.StatisticMinimum,
		TreatMissingData:   aws.String("notBreaching"),
	})
	if err != nil {
		return err
	}

	return nil
}

// PutMetricData puts a metric into CloudWatch with the given name and namespace.
func (c *CloudWatch) PutMetricData(metricName, namespace string, value float64) error {
	_, err := c.client.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
		Namespace: aws.String(namespace),
		MetricData: []types.MetricDatum{
			{
				MetricName: aws.String(metricName),
				Unit:       types.StandardUnitCount,
				Value:      aws.Float64(value),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
