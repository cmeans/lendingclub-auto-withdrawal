# lendingclub-auto-withdrawal

Simple routine to transfer all available cash out of LendingClub.

LendingClub has stopped trading Notes, see their [site](https://www.lendingclub.com/account/summary.action) for details.

However, they have not provided a mechanism to automatically return our money, and, as much fun as it is to log into the site every few days and manually transfer the available cash, I decided I wanted an easier way to get my money back.

Fortunately, LendingClub does provide free access to a [developer API](https://www.lendingclub.com/developers), and that is what is being leveraged here.

Though my original approach was a combination of [IFTTT](https://ifttt.com) and [AWS Lambda](https://aws.amazon.com/lambda/), I later learned about [Amazon EventBridge](https://aws.amazon.com/eventbridge/), and now have a scheduled rule take care of triggering the Lambda, and everything is wrapped into a [AWS CloudFormation](https://docs.aws.amazon.com/cloudformation/index.html) [stack](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/stacks.html).

## Installation

The AWS Lambda is written in [Go](https://golang.org/), which must be compiled for the Amazon target platform, so you will need to have/install a recent version of Go...but don't be concerned, you will not need to learn Go or much of anything else to use this stack.

I also use the AWS CLI utility to upload the Lambda, you can avoid installing it and upload it yourself manually.

### Download & Install

* [GoLang](https://golang.org/dl/)
* [AWS CLI](https://aws.amazon.com/cli/)

## Setup

1. Get a [Lending Club Developer API Key](https://www.lendingclub.com/account/profile.action).

2. Setup an AWS Account.  Assuming you only trigger the lambda to run once a day at most, you should be able to operate within the free tier, but that is dependent upon your usage, and is your responsibility to keep an eye on.

3. Run this script:
Note:  You'll need to setup your AWS account locally or at least set some environment vars so the ```aws cli``` will run without issue.
```bash
./uploadLambdaToS3.sh
```
The script will compile, and then upload the lambda code to a new S3 Bucket (it creates a new bucket each time it's run) you can delete the bucket once the stack installation is complete.  Check the output of the script by viewing the file ```uploadLambdaToS3.output.txt```.  The very last line will include the bucket name, and the bucket key of the lambda...you'll need them soon.

4. [Sign into the AWS console](https://aws.amazon.com/)

5. Go to CloudFormation.

6. Create a Stack.

7. Select "Template is ready".

8. Select "Upload a template".

9. Click "Choose file", and navigate to where you cloned this repo, and select the file named: ```stack-template.json```.

...more to come...
