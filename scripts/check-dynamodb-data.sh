#!/bin/bash

# This script checks the data in the Aurora MySQL database.
aws --endpoint-url=http://localhost:4566 dynamodb execute-statement \
    --statement "SELECT * FROM street_segment_speeds"
