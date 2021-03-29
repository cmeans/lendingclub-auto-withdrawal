#!/bin/bash

ECHO Compile the lambda...
pushd lambda
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/main main.go
popd

BUCKET=my-auto-withdrawal-lambda-$(uuidgen | tr [:upper:] [:lower:])

# Create a unique bucket.
aws s3 mb s3://$BUCKET

# Upload a unique bucket.
aws cloudformation package \
    --template uploadLambdaToS3.json \
    --s3-bucket $BUCKET \
    --output json > uploadLambdaToS3.output.txt
