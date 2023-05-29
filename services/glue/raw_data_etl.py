import sys
import time
import boto3
from pyspark.context import SparkContext
from pyspark.sql.functions import col, upper
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


def main():
    sc = SparkContext.getOrCreate()
    glue_context = GlueContext(sc)
    job = Job(glue_context)

    args = getResolvedOptions(sys.argv, ["JOB_NAME"])

    LOG_GROUP_NAME = "/aws/glue/jobs"
    LOG_STREAM_NAME = get_running_job_id(args["JOB_NAME"])
    logs = create_log_group(LOG_GROUP_NAME, LOG_STREAM_NAME)

    job.init(args["JOB_NAME"], args)

    source_bucket = "raw-data"
    # TODO: Change to something more dynamic.
    source_key = "year=2023/month=05/day=29/batch-from-d123-to-d123.csv"
    source_path = f"s3://{source_bucket}/{source_key}"
    df = glue_context.create_dynamic_frame.from_options(
        connection_type="s3",
        connection_options={"paths": [source_path]},
        format="csv",
        format_options={"withHeader": True},
    ).toDF()

    log(
        logs,
        LOG_GROUP_NAME,
        LOG_STREAM_NAME,
        f"Read {df.count()} rows from {source_path}.",
    )

    df_transformed = df.withColumn("id", col("id")).withColumn(
        "title", upper(col("title"))
    )

    log(
        logs,
        LOG_GROUP_NAME,
        LOG_STREAM_NAME,
        f"Transformed {df_transformed.count()} rows.",
    )

    df_transformed = DynamicFrame.fromDF(df_transformed, glue_context, "df_transformed")

    target_bucket = "transformed-data"
    # TODO: Change to something more dynamic.
    target_key = "year=2023/month=05/day=29/"
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

    job.commit()


main()
