package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsssm"
	"github.com/aws/constructs-go/constructs/v10"
	_jsii_ "github.com/aws/jsii-runtime-go"
)

type DataStack struct {
	awscdk.Stack
	Table                 awsdynamodb.Table
	JwtSecret             awssecretsmanager.Secret
	QiitaTokenSecret      awssecretsmanager.Secret
	MeIDParameter         awsssm.StringParameter
	ZennUsernameParameter awsssm.StringParameter
	LogLevelParameter     awsssm.StringParameter
}

func NewDataStack(scope constructs.Construct, id string, props *StackProps) *DataStack {
	stack := newStack(scope, id, props)
	dataStack := &DataStack{Stack: stack}

	if props == nil || props.Config == nil {
		return dataStack
	}

	cfg := props.Config

	table := awsdynamodb.NewTable(stack, _jsii_.String("Table"), &awsdynamodb.TableProps{
		TableName:   _jsii_.String(cfg.Data.TableName),
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
		PartitionKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("PK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TimeToLiveAttribute: _jsii_.String("ttl"),
		RemovalPolicy:       awscdk.RemovalPolicy_RETAIN,
	})

	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: _jsii_.String("GSI1"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("GSI1PK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("GSI1SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: _jsii_.String("GSI2"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("GSI2PK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("GSI2SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: _jsii_.String("GSI3"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("GSI3PK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("GSI3SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: _jsii_.String("GSI_EMAIL"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("GSI_EMAIL_PK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	jwtSecret := awssecretsmanager.NewSecret(stack, _jsii_.String("JwtSecret"), &awssecretsmanager.SecretProps{
		SecretName:  _jsii_.String(cfg.SecretName("jwtSecret")),
		Description: _jsii_.String("JWT signing secret for the me API."),
		GenerateSecretString: &awssecretsmanager.SecretStringGenerator{
			ExcludePunctuation: _jsii_.Bool(true),
		},
		RemovalPolicy: awscdk.RemovalPolicy_RETAIN,
	})

	qiitaToken := awssecretsmanager.NewSecret(stack, _jsii_.String("QiitaTokenSecret"), &awssecretsmanager.SecretProps{
		SecretName:  _jsii_.String(cfg.SecretName("qiitaToken")),
		Description: _jsii_.String("Qiita API token placeholder. Replace the generated value after deployment."),
		GenerateSecretString: &awssecretsmanager.SecretStringGenerator{
			ExcludePunctuation: _jsii_.Bool(true),
		},
		RemovalPolicy: awscdk.RemovalPolicy_RETAIN,
	})

	meIDParameter := awsssm.NewStringParameter(stack, _jsii_.String("MeIdParameter"), &awsssm.StringParameterProps{
		ParameterName: _jsii_.String(cfg.ParameterName("meId")),
		Description:   _jsii_.String("Profile ID used by the me application."),
		StringValue:   _jsii_.String(cfg.Data.MeID),
		Tier:          awsssm.ParameterTier_STANDARD,
	})

	zennUsernameParameter := awsssm.NewStringParameter(stack, _jsii_.String("ZennUsernameParameter"), &awsssm.StringParameterProps{
		ParameterName: _jsii_.String(cfg.ParameterName("zennUsername")),
		Description:   _jsii_.String("Zenn username used by the me application."),
		StringValue:   _jsii_.String(cfg.Data.ZennUsername),
		Tier:          awsssm.ParameterTier_STANDARD,
	})

	logLevelParameter := awsssm.NewStringParameter(stack, _jsii_.String("LogLevelParameter"), &awsssm.StringParameterProps{
		ParameterName: _jsii_.String(cfg.ParameterName("logLevel")),
		Description:   _jsii_.String("Application log level."),
		StringValue:   _jsii_.String(cfg.Data.LogLevel),
		Tier:          awsssm.ParameterTier_STANDARD,
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("TableNameOutput"), &awscdk.CfnOutputProps{
		Value:      table.TableName(),
		ExportName: _jsii_.String(cfg.ExportName("data", "table-name")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("TableArnOutput"), &awscdk.CfnOutputProps{
		Value:      table.TableArn(),
		ExportName: _jsii_.String(cfg.ExportName("data", "table-arn")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("JwtSecretArnOutput"), &awscdk.CfnOutputProps{
		Value:      jwtSecret.SecretArn(),
		ExportName: _jsii_.String(cfg.ExportName("data", "jwt-secret-arn")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("QiitaTokenSecretArnOutput"), &awscdk.CfnOutputProps{
		Value:      qiitaToken.SecretArn(),
		ExportName: _jsii_.String(cfg.ExportName("data", "qiita-token-secret-arn")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("MeIdParameterNameOutput"), &awscdk.CfnOutputProps{
		Value:      meIDParameter.ParameterName(),
		ExportName: _jsii_.String(cfg.ExportName("data", "me-id-parameter-name")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("ZennUsernameParameterNameOutput"), &awscdk.CfnOutputProps{
		Value:      zennUsernameParameter.ParameterName(),
		ExportName: _jsii_.String(cfg.ExportName("data", "zenn-username-parameter-name")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("LogLevelParameterNameOutput"), &awscdk.CfnOutputProps{
		Value:      logLevelParameter.ParameterName(),
		ExportName: _jsii_.String(cfg.ExportName("data", "log-level-parameter-name")),
	})

	dataStack.Table = table
	dataStack.JwtSecret = jwtSecret
	dataStack.QiitaTokenSecret = qiitaToken
	dataStack.MeIDParameter = meIDParameter
	dataStack.ZennUsernameParameter = zennUsernameParameter
	dataStack.LogLevelParameter = logLevelParameter

	return dataStack
}
