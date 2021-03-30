#!/bin/bash

# Empty bucket.
aws s3 rm s3://$1 --recursive

# Remove bucket.
aws s3 rb s3://$1
