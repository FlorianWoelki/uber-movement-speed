# Preprocessing Service

The Preprocessing service is a service that allows you to preprocess the raw data from the
Kinesis Data Stream and upload it to a new S3 bucket in CSV format. It also adds the data
to a dynamodb table for further possible processing by users.

## Building the service

To build the service for deployment, you have to run the following command:

```sh
$ ./build.sh
```

This simple bash script will create an executable file that can be run on linux systems.
This executable will be zipped and uploaded to AWS Lambda.
