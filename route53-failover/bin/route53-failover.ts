#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import {Route53FailoverStack} from '../lib/route53-failover-stack';

const app = new cdk.App();
new Route53FailoverStack(app, 'Route53FailoverStack', {
  env: {account: process.env.CDK_DEFAULT_ACCOUNT, region: process.env.CDK_DEFAULT_REGION},
});