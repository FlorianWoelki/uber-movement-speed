import time
import sys
import boto3

endpoint_url = "http://localhost.localstack.cloud:4566"


def start_glue_job(job_name: str):
    """
    Starts a glue job with the given name.

    Args:
        job_name (str): The name of the glue job.

    Returns:
        str: The ID of the glue job run.
    """
    glue = boto3.client("glue", endpoint_url=endpoint_url)
    response = glue.start_job_run(
        JobName=job_name,
    )
    return response["JobRunId"]


def stop_glue_job(job_name: str, run_id: str):
    """
    Stops a glue job with the given name and run ID.

    Args:
        job_name (str): The name of the glue job.
        run_id (str): The ID of the glue job run.
    """
    glue = boto3.client("glue", endpoint_url=endpoint_url)
    glue.batch_stop_job_run(
        JobName=job_name,
        JobRunIds=[run_id],
    )


def get_run_logs(job_name: str, run_id: str) -> list[str]:
    """
    Gets the logs of a glue job with the given name and run ID.

    Args:
        job_name (str): The name of the glue job.
        run_id (str): The ID of the glue job run.

    Returns:
        list[str]: The logs of the glue job run.
    """
    glue = boto3.client("glue", endpoint_url=endpoint_url)
    runs = glue.get_job_runs(
        JobName=job_name,
    )["JobRuns"]
    run_id = next(run["Id"] for run in runs if run["Id"] == run_id)

    logs = boto3.client("logs", endpoint_url=endpoint_url)
    response = logs.get_log_events(
        logGroupName=f"/aws/glue/jobs",
        logStreamName=run_id,
    )

    result: list[str] = []
    for event in response["events"]:
        result.append(event["message"])

    return result


def main():
    job_name = "raw-data-etl"

    # Starts a glue job.
    if len(sys.argv) == 2 and sys.argv[1] == "start":
        job_run_id = start_glue_job(job_name)
        print(f"Started Glue job run with ID: {job_run_id}")

        # Wait for the job to finish.
        glue = boto3.client("glue", endpoint_url=endpoint_url)
        job_run = glue.get_job_run(JobName=job_name, RunId=job_run_id)
        while job_run["JobRun"]["JobRunState"] == "RUNNING":
            job_run = glue.get_job_run(JobName=job_name, RunId=job_run_id)
            print(f"Glue job run with ID: {job_run_id} is still running...")
            time.sleep(4)

        print(
            f"Glue job run with ID: {job_run_id} finished with status: {job_run['JobRun']['JobRunState']}"
        )
        return

    # Gets the logs of a glue job.
    if len(sys.argv) == 3 and sys.argv[1] == "logs":
        job_run_id = sys.argv[2]
        logs = get_run_logs(job_name, job_run_id)
        for log in logs:
            print(log)
        return

    # Stops a glue job.
    if len(sys.argv) == 3 and sys.argv[1] == "stop":
        job_run_id = sys.argv[2]
        stop_glue_job(job_name, job_run_id)
        print(f"Stopped Glue job run with ID: {job_run_id}")
        return

    print(f"Usage: {sys.argv[0]} [start|stop|logs]")


if __name__ == "__main__":
    main()
