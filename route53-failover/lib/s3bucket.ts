import {Construct} from "constructs";
import * as s3 from "aws-cdk-lib/aws-s3";
import * as s3Deployment from "aws-cdk-lib/aws-s3-deployment";
import * as path from "node:path";
import * as cdk from "aws-cdk-lib";
import {RemovalPolicy} from "aws-cdk-lib";

export class FailOverS3Bucket extends Construct {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    const bucket = new s3.Bucket(scope, "FailOverS3Bucket", {
      websiteIndexDocument: 'index.html',
      publicReadAccess: true,
      blockPublicAccess: new s3.BlockPublicAccess({
        blockPublicAcls: false,
        blockPublicPolicy: false,
        restrictPublicBuckets: false,
        ignorePublicAcls: false
      }),
      versioned: true,
      bucketName: 'www.worldwideapex.com',
      autoDeleteObjects: true,
      removalPolicy: RemovalPolicy.DESTROY
    })
    
    new s3Deployment.BucketDeployment(scope, "FailOverS3BucketDeployment", {
      sources: [s3Deployment.Source.asset(path.resolve(__dirname, 'data'))],
      destinationBucket: bucket
    });


    new cdk.CfnOutput(this, 'BucketWebsiteUrl', {
      value: bucket.bucketWebsiteUrl
    });
  }
}