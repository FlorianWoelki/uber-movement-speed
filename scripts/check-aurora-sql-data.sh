#!/bin/bash

# This script checks the data in the Aurora MySQL database.
# Gets the database cluster.
CLUSTER=$(aws --endpoint-url=http://localhost:4566 rds describe-db-clusters --db-cluster-identifier db1)
# Gets the ARN of the database cluster.
CLUSTER_ARN=$(echo $CLUSTER | jq -r '.DBClusters[0].DBClusterArn')

# Gets the secret of the database cluster.
SECRET=$(aws --endpoint-url=http://localhost:4566 secretsmanager describe-secret --secret-id dbpass)
# Gets the ARN of the secret.
SECRET_ARN=$(echo $SECRET | jq -r '.ARN')

# Selects all the data in the table `street_segment_speeds` in the database.
aws --endpoint-url=http://localhost:4566 rds-data execute-statement \
    --resource-arn $CLUSTER_ARN \
    --secret-arn $SECRET_ARN \
    --database test \
    --sql "SELECT * FROM street_segment_speeds" \
    --output json
