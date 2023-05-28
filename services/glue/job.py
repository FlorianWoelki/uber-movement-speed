import time
import sys
import boto3

endpoint_url = "http://localhost.localstack.cloud:4566"


def create_glue_job(job_name: str, script_location: str):
    """
    Creates a glue job with the given name and script location.

    Args:
        job_name (str): The name of the glue job.
        script_location (str): The location of the script to be executed by the glue job.

    Returns:
        str: The name of the glue job.
    """
    glue = boto3.client("glue", endpoint_url=endpoint_url)

    response = glue.create_job(
        Name=job_name,
        Role="arn:aws:iam::000000000000:role/glue-role",
        Command={
            "Name": "pythonshell",
            "ScriptLocation": script_location,
        },
    )

    return response["Name"]


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


def main():
    job_name = "raw-data-etl"
    script_location = "s3://raw-data/scripts/raw_data_etl.py"

    # Starts only a glue job.
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

    # Creats and start the glue job.
    create_glue_job(job_name, script_location)
    job_run_id = start_glue_job(job_name)
    print(f"Started Glue job run with ID: {job_run_id}")


if __name__ == "__main__":
    main()
