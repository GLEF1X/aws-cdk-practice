import * as cdk from 'aws-cdk-lib';
import {aws_ec2} from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {R53ResolverVPC} from "./vpc";
import {EC2InstanceConstruct} from "./ec2-instance";
import {FailOverS3Bucket} from "./s3bucket";

export class Route53FailoverStack extends cdk.Stack {
  protected targetVpc: aws_ec2.Vpc;

  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    this.targetVpc = new R53ResolverVPC(this, "R53ResolverTestVPC").vpc;
    new EC2InstanceConstruct(this, "EC2FailOverInstanceConstruct", {
      vpc: this.targetVpc
    })
    new FailOverS3Bucket(this, "FailOverS3Bucket")
  }
}
