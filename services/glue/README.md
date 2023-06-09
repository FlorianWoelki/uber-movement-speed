# Glue Service

This service was developed and tested with the Python version `3.11.3`.

The Glue service is a service that allows you to create a glue job and run it. This service
is split up into two files: `job.py` and `raw_data_etl.py`. The `job.py` file is the simple
CLI that allows you to create and run a glue job or seeing the logs of a glue job.
The `raw_data_etl.py` file is the actual glue job that is run. This file contains the
ETL logic for the raw data and uploads it to a new S3 bucket and a Aurora instance for
further machine learning processes. The `raw_data_etl.py` file is uploaded to S3 and then
the `job.py` file is run to create and run the glue job.

## Running the glue job

To run the glue job, you will need to install the dependencies with
`pip install -r requirements.txt`. After that, you can checkout the CLI for further help
by running `python job.py`. The CLI has the following options:

```sh
$ python job.py start # Starts the glue job.
$ python job.py stop # Stop the glue job.
$ python job.py logs # Print the logs of the glue job.
```
