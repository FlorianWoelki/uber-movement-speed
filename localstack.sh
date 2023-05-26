#!/bin/bash
source .env

# Tries to start localstack for AWS emulation.
export LOCALSTACK_API_KEY=${LOCALSTACK_API_KEY}
localstack start
