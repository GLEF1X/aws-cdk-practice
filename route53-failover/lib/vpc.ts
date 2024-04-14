import {Construct} from "constructs";
import * as ec2 from "aws-cdk-lib/aws-ec2";

export class R53ResolverVPC extends Construct {
  public vpc: ec2.Vpc;

  constructor(scope: Construct, id: string) {
    super(scope, id);
    // @see https://docs.aws.amazon.com/cdk/api/v2/docs/aws-cdk-lib.aws_ec2.Vpc.html
    this.vpc = new ec2.Vpc(this, "R53FailoverTestVPC", {
      ipAddresses: ec2.IpAddresses.cidr("10.24.0.0/16"),
      subnetConfiguration: [
        {
          cidrMask: 24,
          name: "DB",
          subnetType: ec2.SubnetType.PRIVATE_ISOLATED,
        },
        {
          cidrMask: 24,
          name: "Application",
          subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS,
        },
        {
          cidrMask: 24,
          name: "Web",
          subnetType: ec2.SubnetType.PUBLIC,
        },
      ],
      natGateways: 0,
      enableDnsSupport: true,
      enableDnsHostnames: true,
    });
  }
}