# Dynamo Getter Service

This service is responsible for getting data from DynamoDB by the query parameter string
`id`. It will be deployed as an AWS Lambda function to a HTTP endpoint.

This service is accessible under the following GET request:

```text
GET http://localhost:4566/restapis/<api-id>/dev/_user_request_/dynamo-getter?id=<id>
```

## Building the service

To build the service for deployment, you have to run the following command:

```sh
$ ./build.sh
```

This simple bash script will create an executable file that can be run on linux systems.
This executable will be zipped and uploaded to AWS Lambda.
