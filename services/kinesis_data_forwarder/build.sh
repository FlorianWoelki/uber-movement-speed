#!/bin/bash

# Builds the go binary for the service.
GOOS=linux GOARCH=amd64 go build -o main main.go

# Zips the binary and the dependencies.
zip -r kinesis_data_forwarder.zip main

# Removes the binary.
rm main
