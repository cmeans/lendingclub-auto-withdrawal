# lendingclub-auto-withdrawal

Simple routine to transfer all available cash out of LendingClub.

LendingClub has stopped trading Notes, see their [site](https://www.lendingclub.com/account/summary.action) for details.

However, they have not provided a mechanism to automatically return our money, and, as much fun as it is to log into the site every few days and manually transfer the available cash, I decided I wanted an easier way to get my money back.

Fortunately, LendingClub does provide free access to a [developer API](https://www.lendingclub.com/developers), and that is what is being leveraged here.

You have choices in how to use/deploy this routine:

1. Compile and run the routine locally as a standalone utility.  You can run it manually, or set up a scheduler to trigger it for you automatically on whatever schedule you like.

2. Deploy it as an AWS CloudFormation Stack (you will need to have/create an AWS account) and have AWS run the routine for you on whatever schedule you like.  This routine should have no problems keeping you in the Free Tier, but that's under your control, so be careful.

3. You come up with some other approach, deploying it as a Docker container or whatever.  You're free to use this code however you would like.

## Approach #1: Standalone utility

### Download & Install

* [GoLang](https://golang.org/dl/)

1. Get a [Lending Club Developer API Key](https://www.lendingclub.com/account/profile.action).

2. Clone this repo:
```bash
git clone https://github.com/cmeans/lendingclub-auto-withdrawal
```

3. Move into the project and then lambda code folder:
```bash
cd lendingclub-auto-withdrawal/lendingclub-auto-withdrawal
```

4. Compile with:
```bash
go build
```

5. Run the utility:
```bash
./lendingclub-auto-withdrawal <investorID> <lendingclubAPIKey>
```

As it is compiled code, you can copy/move the utility wherever it'll be used.

## Approach #2: AWS CloudFormation Stack

Though my original tack was a combination of [IFTTT](https://ifttt.com) and [AWS Lambda](https://aws.amazon.com/lambda/), I later learned about [Amazon EventBridge](https://aws.amazon.com/eventbridge/), and now have a scheduled rule take care of triggering the Lambda, and everything is wrapped into a [AWS CloudFormation](https://docs.aws.amazon.com/cloudformation/index.html) [stack](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/stacks.html).

### Installation

The AWS Lambda is written in [Go](https://golang.org/), which must be compiled for the Amazon target platform, so you will need to have/install a recent version of Go...but don't be concerned, you will not need to learn Go or much of anything else to use this stack.

The deployment routines handle everything for you, but they rely on the AWS CLI.  If you do not want to install it you will need to create the necessary resources yourself.  Inspect the `stack-template.json` for all the goodies.

#### Download & Install

* [GoLang](https://golang.org/dl/)
* [AWS CLI](https://aws.amazon.com/cli/)

### Setup (for Mac users)

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

You will need to setup your AWS account locally or at least set some environment variables so `aws` will run without issue.

You will be prompted for your LendingClub Investor ID, Developer API key.  If you provide dummy/fake values, the Lambda will not work, but you can log into your AWS account and fix the values manually at any time.

```bash
./deploy.sh
```

By default, the routine will create a CloudFormation Stack called *lendingclub-auto-withdrawal*, and will schedule the Lambda to trigger at 6pm Central.  Both of these defaults can be changed by editing the `deploy.sh` script before yuo run it.

The script will compile the lambda, and then upload the ZIPped code to a new S3 Bucket (it creates a new bucket each time it's run) you can delete the bucket once the stack installation/update is done.

6. Delete the S3 bucket, by running this script, with the bucket name as it's only parameter:

```bash
./remove-bucket.sh <bucket-name>
```

7. If you mistyped your InvestorID or API Key, you can fix them in the Lambda settings via the AWS console.

8. By default the routine will not initiate a transfer if the availableCash is $10.00 or less.  You can override this default value by creating a new Environment Variable in the Lambda Configuration section, with the name MINIMUM_AMOUNT set to an integer value.

9. If you encounter issues (sometimes the deployment has some timing issues) it's perfectly fine to delete the stack from AWS and redeploy it as many times as you like.

10. If you get adventurous and update the code or make a tweak to the `stack-template.json`, you can simply redeploy your changes using the `redeploy.sh` script.
