# Kinesis Data Forwarder Service

The Kinesis Data Forwarder service is a service that allows you to forward data from the
API Gateway to a Kinesis Data Stream. For that to work, the payload from the websocket over
the API gateway, needs to contain an action which is set to `kinesis-data-forwarder`.

## Building the service

To build the service for deployment, you have to run the following command:

```sh
$ pnpm build
```

This will create a `dist` folder with the compiled code and a zip file that will be
uploaded to AWS Lambda.
