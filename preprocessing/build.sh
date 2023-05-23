#!/bin/bash

# Builds the go binary for the service.
GOOS=linux GOARCH=amd64 go build -o main lambda.go

# Zips the binary and the dependencies.
zip -r my-code.zip main

# Removes the binary.
rm main
