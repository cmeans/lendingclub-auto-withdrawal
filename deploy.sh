#!/bin/bash

read -p 'LendingClub Investor ID (*required): ' INVESTOR_ID
if [[ $INVESTOR_ID == "" ]]; then
    echo You must supply your LendingClub Investor ID.
    exit
fi

read -p 'LendingClub API Key (*required): ' API_KEY
if [[ $API_KEY == "" ]]; then
    echo You must supply your LendingClub Developer API KEY.
    exit
fi

read -p 'AWS CloudFormation Stack Name (LendingClub-Auto-Withdrawal): ' STACK_NAME
read -p 'Rule schedule (cron(0 23 * * ? *)): ' RULE_SCHEDULE

if [[ $STACK_NAME == "" ]]; then
    STACK_NAME="LendingClub-Auto-Withdrawal"
fi

if [[ $RULE_SCHEDULE == "" ]]; then
    RULE_SCHEDULE="cron(0 23 * * ? *)"
fi

ECHO Preparing Resources...
ECHO Compiling the AWS Lambda...
pushd lambda > /dev/null
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/main main.go
popd > /dev/null

ECHO Create a unique AWS S3 bucket name.
BUCKET=my-auto-withdrawal-lambda-$(uuidgen | tr [:upper:] [:lower:])

ECHO Creating the bucket.
# Create a unique bucket.
aws s3 mb s3://$BUCKET

ECHO Packaging and uploading the Lambda to the S3 bucket...
# Upload the code to our new bucket.
aws cloudformation package \
    --template-file stack-template.json \
    --s3-bucket $BUCKET \
    --output-template-file stack-template.transformed.json \
    --use-json > /dev/null

ECHO Uploading the AWS CloudFormation Stack template...to kickoff the Stack creation...
# Create the new AWS CloudFormation Stack.
aws cloudformation create-stack \
    --stack-name $STACK_NAME \
    --capabilities CAPABILITY_AUTO_EXPAND CAPABILITY_NAMED_IAM CAPABILITY_IAM \
    --template-body file://stack-template.transformed.json \
    --parameters ParameterKey=LendingClubInvestorID,ParameterValue="$INVESTOR_ID" ParameterKey=LendingClubAPIKey,ParameterValue="$API_KEY" ParameterKey=RuleSchedule,ParameterValue="$RULE_SCHEDULE"

ECHO "You can log into the AWS Console, CloudFormation Stacks to see if the creation process is complete."
ECHO "Or monitor it from the command line using:"
ECHO
ECHO "aws cloudformation describe-stack-events --stack-name LendingClub-Auto-Withdrawal --max-items 1"
ECHO
ECHO "Once it's complete, run the command below to remove the S3 bucket created at the start, as it is no longer needed."
ECHO
ECHO "./remove-bucket.sh $BUCKET"
ECHO
