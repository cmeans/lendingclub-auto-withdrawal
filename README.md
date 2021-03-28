# lendingclub-auto-withdrawl

Simple routine to transfer all available cash out of LendingClub.

LendingClub has stopped trading Notes, see their [site](https://www.lendingclub.com/account/summary.action) for details.

However, they have not provided a mechanism to automatically retrieve our money, as much fun as it is to log into the site every few days, I decided it would be easier to leverage their [API](https://www.lendingclub.com/developers), and do it automatically myself.

My approach:

* [IFTTT](https://ifttt.com) routine to trigger a webhook daily.
* [AWS API Gateway](https://aws.amazon.com/api-gateway/) is the webhook.
* The [AWS Lambda](https://aws.amazon.com/lambda/) is fired by the API Gateway.

## Install

* GoLang
* AWS CLI

## Setup

* Get a Lending Club Developer Key.
* Setup an AWS Account.
* Create a Role to execute the Lambda.
* Update the ```deploy``` script with the Role, your LendingClub investor ID, and API key.
* ...
