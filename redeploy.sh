#!/bin/bash

ECHO Preparing Resources...
ECHO Compiling the AWS Lambda...
pushd lendingclub-auto-withdrawal > /dev/null
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/main main.go
popd > /dev/null

ECHO Create a unique AWS S3 bucket name.
BUCKET=my-auto-withdrawal-lambda-$(uuidgen | tr [:upper:] [:lower:])

ECHO Creating the bucket...
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
aws cloudformation update-stack \
    --stack-name ${STACK_NAME:-LendingClub-Auto-Withdrawal} \
    --capabilities CAPABILITY_AUTO_EXPAND CAPABILITY_NAMED_IAM CAPABILITY_IAM \
    --template-body file://stack-template.transformed.json \
    --parameters ParameterKey=LendingClubInvestorID,UsePreviousValue=true ParameterKey=LendingClubAPIKey,UsePreviousValue=true ParameterKey=RuleSchedule,UsePreviousValue=true

if [[ $? != 0 ]]; then
    ECHO Deployment error...so deleting the S3 Bucket now...
    ./remove-bucket.sh $BUCKET
else
    ECHO "You can log into the AWS Console, and go to CloudFormation/Stacks to see if the update process is complete."
    ECHO "Or monitor it from the command line using:"
    ECHO
    ECHO "    aws cloudformation describe-stack-events --stack-name ${STACK_NAME:-LendingClub-Auto-Withdrawal} --max-items 1"
    ECHO
    ECHO "Once it's complete, run the command below to remove the S3 bucket created at the start, as it is no longer needed."
    ECHO
    ECHO "    ./remove-bucket.sh $BUCKET"
    ECHO
fi
