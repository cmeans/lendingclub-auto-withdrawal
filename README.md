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

## Setup (for Mac users)

1. Get a [Lending Club Developer API Key](https://www.lendingclub.com/account/profile.action).

2. Setup an AWS Account.  Assuming you only trigger the lambda to run once a day at most, you should be able to operate within the free tier, but that is dependent upon *your* usage, and is *your responsibility* to keep an eye on.

3. Clone this repo:
```bash
git clone https://github.com/cmeans/lendingclub-auto-withdrawal
```

4. Move into the project folder:
```
cd lendingclub-auto-withdrawal
```

5. Execute the `deploy.sh` script:

You will need to setup your AWS account locally or at least set some environment variables so the `aws cli` will run without issue.

You will be prompted for your LendingClub Investor ID, Developer API key, what to name the application stack, and what [cron](https://docs.aws.amazon.com/lambda/latest/dg/services-cloudwatchevents-expressions.html) schedule to use for regularly triggering the lambda.
You *must* have your Investor ID and API Key, the last two options provide default values, which I suggest you use, at least initially.

```bash
./deploy.sh
```

The script will compile, and then upload the lambda code to a new S3 Bucket (it creates a new bucket each time it's run) you can delete the bucket once the stack installation is complete.

6. Once the stack creation is complete, or there was an issue, you should delete the S3 bucket as it will not be used again.

You can run this script, with the bucket name as it's only parameter:

```bash
./remove-bucket.sh <bucket-name>
```

7. If you mistyped your InvestorID or API Key, you can fix them in the Lambda settings via the AWS console.

8. By default the routine will not initiate a transfer if the availableCash is $10.00 or less.  You can override this default value by creating a new Environment Variable in the Lambda Configuration section, with the name MINIMUM_AMOUNT set to an integer value.
