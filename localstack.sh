#!/bin/bash
source .env

export LOCALSTACK_API_KEY=${LOCALSTACK_API_KEY}
export ENFORCE_IAM=1

# Check if `LOCALSTACK_API_KEY`` is set.
if [ -z "$LOCALSTACK_API_KEY" ]; then
  echo "LOCALSTACK_API_KEY is not set. Please set it in .env file."
  exit 1
fi

# Checks if container exists and removes it.
if docker ps -a | grep localstack_main; then
  docker rm -f localstack_main
fi

# Check if `localstack` is locally installed.
if ! command -v localstack &> /dev/null
then
  # Start localstack with `docker`.
  docker compose up --build
else
  # Start localstack with `localstack` command.
  localstack start
fi
