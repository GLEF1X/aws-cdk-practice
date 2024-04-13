import {Duration, Stack} from "aws-cdk-lib";
import {Construct} from "constructs";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as iam from "aws-cdk-lib/aws-iam";
import {CloudFormationInit, InstanceClass, InstanceSize, InstanceType, MachineImage} from "aws-cdk-lib/aws-ec2";

export interface EC2InstanceConstructProps {
  vpc: ec2.Vpc
}

export class EC2InstanceConstruct extends Construct {
  constructor(scope: Construct, id: string, props: EC2InstanceConstructProps) {
    super(scope, id);
    const {vpc} = props

    const ec2InstanceSecurityGroup = new ec2.SecurityGroup(
      this,
      'ec2InstanceSecurityGroup',
      {vpc: vpc, allowAllOutbound: true},
    );
    ec2InstanceSecurityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.allTraffic())

    const serverRole = new iam.Role(this, 'serverEc2Role', {
      assumedBy: new iam.ServicePrincipal('ec2.amazonaws.com'),
      inlinePolicies: {
        ['RetentionPolicy']: new iam.PolicyDocument({
          statements: [
            new iam.PolicyStatement({
              resources: ['*'],
              actions: ['logs:PutRetentionPolicy'],
            }),
          ],
        }),
      },
      managedPolicies: [
        iam.ManagedPolicy.fromAwsManagedPolicyName('AmazonSSMManagedInstanceCore'),
        iam.ManagedPolicy.fromAwsManagedPolicyName('CloudWatchAgentServerPolicy'),
      ],
    });

    new ec2.Instance(this, "MainInstance", {
      vpc: vpc,
      allowAllOutbound: true,
      instanceType: InstanceType.of(InstanceClass.T2, InstanceSize.MICRO),
      machineImage: MachineImage.latestAmazonLinux2023(),
      securityGroup: ec2InstanceSecurityGroup,
      init: CloudFormationInit.fromConfigSets({
        configSets: {
          default: ['initial', 'install_httpd', 'start']
        },
        configs: {
          initial: new ec2.InitConfig([
            ec2.InitCommand.shellCommand('sudo yum update -y'),
            ec2.InitFile.fromObject('/etc/stack.json', {
              stackId: Stack.of(this).stackId,
              stackName: Stack.of(this).stackName,
              region: Stack.of(this).region,
            }),
          ]),
          install_httpd: new ec2.InitConfig([
            ec2.InitPackage.yum('httpd'),
          ]),
          start: new ec2.InitConfig([
            ec2.InitCommand.shellCommand('sudo httpd'),
          ])
        },
      }),
      initOptions: {
        timeout: Duration.minutes(10),
        includeUrl: true,
        includeRole: true,
        printLog: true,
      },
      role: serverRole,
    })
  }
}