package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awselasticbeanstalk"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AwsCdkPracticeStackProps struct {
	awscdk.StackProps
}

func NewAwsCdkPracticeStack(scope constructs.Construct, id string, props *AwsCdkPracticeStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	sourceDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting source path: ", err)
		return nil
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	vpc := awsec2.NewVpc(stack, jsii.String("PracticeVPC"), &awsec2.VpcProps{
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				CidrMask:   jsii.Number(24),
				Name:       jsii.String("practice-private"),
				SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
			},
			{
				CidrMask:   jsii.Number(24),
				Name:       jsii.String("practice-public"),
				SubnetType: awsec2.SubnetType_PUBLIC,
			},
		},
	})

	elasticBeanStalkEC2Role := awsiam.NewRole(stack, jsii.String("ElasticBeanstalkEC2Role"), &awsiam.RoleProps{
		RoleName:  jsii.String("practice-aws-elasticbeanstalk-ec2-role"),
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AWSElasticBeanstalkWebTier")),
		},
	})

	elbZipArchive := awss3assets.NewAsset(stack, jsii.String("PracticeGoAppAsset"), &awss3assets.AssetProps{
		Path: jsii.String(filepath.Join(sourceDir, "go-app.zip")),
	})

	app := awselasticbeanstalk.NewCfnApplication(stack, jsii.String("AwsCdkPracticeGoApplication"), &awselasticbeanstalk.CfnApplicationProps{
		ApplicationName: jsii.String("AwsCdkPracticeGoApplication"),
	})

	appVersion := awselasticbeanstalk.NewCfnApplicationVersion(stack, jsii.String("AwsCdkPracticeGoApplicationVersion"), &awselasticbeanstalk.CfnApplicationVersionProps{
		ApplicationName: jsii.String("AwsCdkPracticeGoApplication"),
		SourceBundle: &awselasticbeanstalk.CfnApplicationVersion_SourceBundleProperty{
			S3Bucket: elbZipArchive.S3BucketName(),
			S3Key:    elbZipArchive.S3ObjectKey(),
		},
	})

	elasticBeanstalkApp := awselasticbeanstalk.NewCfnEnvironment(stack, jsii.String("AwsCdkPracticeGoEnvironment"), &awselasticbeanstalk.CfnEnvironmentProps{
		ApplicationName:   app.ApplicationName(),
		EnvironmentName:   jsii.String("AwsCdkPracticeGoEnvironment"),
		SolutionStackName: jsii.String("64bit Amazon Linux 2023 v4.0.5 running Go 1"),
		VersionLabel:      appVersion.Ref(),
		OptionSettings: []awselasticbeanstalk.CfnEnvironment_OptionSettingProperty{
			{
				Namespace:  jsii.String("aws:ec2:vpc"),
				OptionName: jsii.String("VPCId"),
				Value:      vpc.VpcId(),
			},
			{
				Namespace:  jsii.String("aws:ec2:vpc"),
				OptionName: jsii.String("Subnets"),
				Value:      awscdk.Fn_Join(jsii.String(","), &[]*string{(*vpc.PrivateSubnets())[0].SubnetId()}),
			},
			{
				Namespace:  jsii.String("aws:ec2:vpc"),
				OptionName: jsii.String("ELBSubnets"),
				Value:      awscdk.Fn_Join(jsii.String(","), &[]*string{(*vpc.PublicSubnets())[0].SubnetId()}),
			},
			{
				Namespace:  jsii.String("aws:autoscaling:launchconfiguration"),
				OptionName: jsii.String("InstanceType"),
				Value:      jsii.String("t2.micro"),
			},
			{
				Namespace:  jsii.String("aws:autoscaling:launchconfiguration"),
				OptionName: jsii.String("IamInstanceProfile"),
				Value:      elasticBeanStalkEC2Role.RoleName(),
			},
		},
	})

	awscdk.NewCfnOutput(stack, jsii.String("ElasticBeanstalkURL"), &awscdk.CfnOutputProps{
		Value:       elasticBeanstalkApp.AttrEndpointUrl(),
		Description: jsii.String("The URL of the Elastic Beanstalk Application"),
	})

	appVersion.AddDependency(app)
	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewAwsCdkPracticeStack(app, "AwsCdkPracticeStack", &AwsCdkPracticeStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
