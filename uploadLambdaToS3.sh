#!/bin/bash

BUCKET=my-lambda-$(uuidgen | tr [:upper:] [:lower:])

ECHO S3 Bucket name will be: $BUCKET

# Create a unique bucket.
aws s3 mb s3://$BUCKET

aws cloudformation package \
    --template uploadLambdaToS3.json \
    --s3-bucket $BUCKET \
    --output json > uploadLambdaToS3.output.txt
