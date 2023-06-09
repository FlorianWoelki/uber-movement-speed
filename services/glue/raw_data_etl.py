import sys
import time
import boto3
from datetime import datetime
from pyspark.context import SparkContext
from pyspark.sql.functions import col
from pyspark.sql.types import IntegerType, FloatType, LongType
from awsglue.utils import getResolvedOptions
from awsglue.dynamicframe import DynamicFrame
from awsglue.context import GlueContext
from awsglue.job import Job

# Endpoint for localstack.
endpoint_url = "http://localhost.localstack.cloud:4566"


# TODO: Probably remove in the future when `JOB_RUN_ID` is passed to the `pythonshell` job.
def get_running_job_id(job_name: str) -> str:
    """
    Gets the ID of a running glue job with the given name.

    Args:
        job_name (str): The name of the glue job.

    Returns:
        str: The ID of the glue job run.

    Raises:
        Exception: If there is an error in the boto3 client.
    """
    glue_client = boto3.client("glue", endpoint_url=endpoint_url)
    try:
        response = glue_client.get_job_runs(JobName=job_name)
        for res in response["JobRuns"]:
            if res.get("JobRunState") == "RUNNING":
                return res.get("Id")
        else:
            return None
    except boto3.ClientError as e:
        raise Exception(
            "boto3 client error in get_status_of_job_all_runs: " + e.__str__()
        )
    except Exception as e:
        raise Exception(
            "Unexpected error in get_status_of_job_all_runs: " + e.__str__()
        )


def create_log_group(log_group_name: str, log_stream_name: str) -> boto3.client:
    """
    Creates a log group and stream with the given names.

    Args:
        log_group_name (str): The name of the log group.
        log_stream_name (str): The name of the log stream.

    Returns:
        boto3.client: The logs client.
    """
    logs = boto3.client("logs", endpoint_url=endpoint_url)
    logs.create_log_group(logGroupName=log_group_name)
    logs.create_log_stream(logGroupName=log_group_name, logStreamName=log_stream_name)
    return logs


def log(logs: boto3.client, log_group_name: str, log_stream_name: str, message: str):
    """
    Logs a message to the given log group and stream.

    Args:
        logs (boto3.client): The logs client.
        log_group_name (str): The name of the log group.
        log_stream_name (str): The name of the log stream.
        message (str): The message to log.
    """
    timestamp = int(round(time.time() * 1000))
    logs.put_log_events(
        logGroupName=log_group_name,
        logStreamName=log_stream_name,
        logEvents=[
            {
                "timestamp": timestamp,
                "message": f"{time.strftime('%Y/%m/%d %H:%M:%S')} {message}",
            },
        ],
    )


def get_db_and_secret_arns(identifier: str, secret_name: str) -> tuple[str, str]:
    db_client = boto3.client("rds", endpoint_url=endpoint_url)
    clusters = db_client.describe_db_clusters(DBClusterIdentifier=identifier)
    db_cluster_arn = clusters["DBClusters"][0]["DBClusterArn"]

    secret_client = boto3.client("secretsmanager", endpoint_url=endpoint_url)
    secret = secret_client.describe_secret(SecretId=secret_name)

    return db_cluster_arn, secret["ARN"]


def main():
    sc = SparkContext.getOrCreate()
    glue_context = GlueContext(sc)
    job = Job(glue_context)

    args = getResolvedOptions(sys.argv, ["JOB_NAME"])

    LOG_GROUP_NAME = "/aws/glue/jobs"
    LOG_STREAM_NAME = get_running_job_id(args["JOB_NAME"])
    logs = create_log_group(LOG_GROUP_NAME, LOG_STREAM_NAME)

    job.init(args["JOB_NAME"], args)

    now = datetime.now()
    year = now.year
    month = now.strftime("%m")
    day = now.strftime("%d")

    source_bucket = "raw-data"
    source_key = f"year={year}/month={month}/day={day}"
    source_path = f"s3://{source_bucket}/{source_key}"
    df = glue_context.create_dynamic_frame.from_options(
        connection_type="s3",
        connection_options={"paths": [source_path], "recurse": True},
        format="csv",
        format_options={"withHeader": True},
    ).toDF()

    log(
        logs,
        LOG_GROUP_NAME,
        LOG_STREAM_NAME,
        f"Read {df.count()} rows from {source_path}.",
    )

    df_transformed = (
        df.withColumn("id", col("id"))
        .withColumn("year", col("year").cast(IntegerType()))
        .withColumn("month", col("month").cast(IntegerType()))
        .withColumn("day", col("day").cast(IntegerType()))
        .withColumn("hour", col("hour").cast(IntegerType()))
        .withColumn("utc_timestamp", col("utc_timestamp"))
        .withColumn("start_junction_id", col("start_junction_id"))
        .withColumn("end_junction_id", col("end_junction_id"))
        .withColumn("osm_way_id", col("osm_way_id").cast(LongType()))
        .withColumn("osm_start_node_id", col("osm_start_node_id").cast(LongType()))
        .withColumn("osm_end_node_id", col("osm_end_node_id").cast(LongType()))
        .withColumn("speed_mph_mean", col("speed_mph_mean").cast(FloatType()))
        .withColumn("speed_mph_stddev", col("speed_mph_stddev").cast(FloatType()))
    )

    log(
        logs,
        LOG_GROUP_NAME,
        LOG_STREAM_NAME,
        f"Transformed {df_transformed.count()} rows.",
    )

    df_transformed = DynamicFrame.fromDF(df_transformed, glue_context, "df_transformed")

    target_bucket = "transformed-data"
    target_key = f"year={year}/month={month}/day={day}/"
    target_path = f"s3://{target_bucket}/{target_key}"
    glue_context.write_dynamic_frame.from_options(
        frame=df_transformed,
        connection_type="s3",
        connection_options={"path": target_path},
        format="csv",
        format_options={"quoteChar": -1},
    )

    log(
        logs,
        LOG_GROUP_NAME,
        LOG_STREAM_NAME,
        f"Wrote {df_transformed.count()} rows to {target_path}.",
    )

    # Gets the ARNs for the RDS Aurora database and the secret.
    cluster_arn, secret_arn = get_db_and_secret_arns("db1", "dbpass")

    # Saves data to the RDS Aurora database.
    aurora_client = boto3.client("rds-data", endpoint_url=endpoint_url)

    rows = df_transformed.toDF().collect()
    for row in rows:
        parameters = [
            {"name": "year", "value": {"longValue": row["year"]}},
            {"name": "month", "value": {"longValue": row["month"]}},
            {"name": "day", "value": {"longValue": row["day"]}},
            {"name": "hour", "value": {"longValue": row["hour"]}},
            {"name": "utc_timestamp", "value": {"stringValue": row["utc_timestamp"]}},
            {
                "name": "start_junction_id",
                "value": {"stringValue": row["start_junction_id"]},
            },
            {
                "name": "end_junction_id",
                "value": {"stringValue": row["end_junction_id"]},
            },
            {"name": "osm_way_id", "value": {"longValue": row["osm_way_id"]}},
            {
                "name": "osm_start_node_id",
                "value": {"longValue": row["osm_start_node_id"]},
            },
            {"name": "osm_end_node_id", "value": {"longValue": row["osm_end_node_id"]}},
            {"name": "speed_mph_mean", "value": {"doubleValue": row["speed_mph_mean"]}},
            {
                "name": "speed_mph_stddev",
                "value": {"doubleValue": row["speed_mph_stddev"]},
            },
        ]
        aurora_client.execute_statement(
            resourceArn=cluster_arn,
            secretArn=secret_arn,
            database="test",
            sql="INSERT INTO street_segment_speeds (year, month, day, hour, utc_timestamp, start_junction_id, end_junction_id, osm_way_id, osm_start_node_id, osm_end_node_id, speed_mph_mean, speed_mph_stddev) VALUES (:year, :month, :day, :hour, :utc_timestamp, :start_junction_id, :end_junction_id, :osm_way_id, :osm_start_node_id, :osm_end_node_id, :speed_mph_mean, :speed_mph_stddev)",
            parameters=parameters,
        )

    log(
        logs,
        LOG_GROUP_NAME,
        LOG_STREAM_NAME,
        f"Inserted {len(rows)} rows into the Aurora database.",
    )
    job.commit()


main()
