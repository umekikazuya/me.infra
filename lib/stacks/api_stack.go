package stacks

import (
	"path/filepath"
	"runtime"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/constructs-go/constructs/v10"
	_jsii_ "github.com/aws/jsii-runtime-go"
)

type ApiStack struct {
	awscdk.Stack
	Function    awslambda.Function
	FunctionURL awslambda.FunctionUrl
	LogGroup    awslogs.LogGroup
}

func NewApiStack(scope constructs.Construct, id string, props *StackProps) *ApiStack {
	stack := newStack(scope, id, props)
	apiStack := &ApiStack{Stack: stack}

	if props == nil || props.Config == nil || props.Data == nil {
		return apiStack
	}

	cfg := props.Config
	data := props.Data

	environment := map[string]*string{
		"DYNAMODB_TABLE_NAME": data.Table.TableName(),
		"JWT_SECRET":          data.JwtSecret.SecretValue().UnsafeUnwrap(),
		"QIITA_TOKEN":         data.QiitaTokenSecret.SecretValue().UnsafeUnwrap(),
		"ME_ID":               data.MeIDParameter.StringValue(),
		"ZENN_USERNAME":       data.ZennUsernameParameter.StringValue(),
		"LOG_LEVEL":           data.LogLevelParameter.StringValue(),
	}

	logGroup := awslogs.NewLogGroup(stack, _jsii_.String("ApiLogGroup"), &awslogs.LogGroupProps{
		LogGroupName:  _jsii_.String(cfg.APILogGroupName()),
		Retention:     awslogs.RetentionDays_ONE_MONTH,
		RemovalPolicy: awscdk.RemovalPolicy_RETAIN,
	})

	function := awslambda.NewFunction(stack, _jsii_.String("ApiFunction"), &awslambda.FunctionProps{
		FunctionName: _jsii_.String(cfg.APIFunctionName()),
		Description:  _jsii_.String("Placeholder me API Lambda. The app repository replaces the code with update-function-code after initial deployment."),
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Architecture: awslambda.Architecture_ARM_64(),
		Handler:      _jsii_.String("bootstrap"),
		Code: awslambda.Code_FromAsset(_jsii_.String(placeholderAPIAssetPath()), &awss3assets.AssetOptions{
			DeployTime: _jsii_.Bool(true),
		}),
		Environment: &environment,
		MemorySize:  _jsii_.Number(512),
		Timeout:     awscdk.Duration_Seconds(_jsii_.Number(30)),
		LogGroup:    logGroup,
	})

	data.Table.GrantReadWriteData(function)
	data.JwtSecret.GrantRead(function, nil)
	data.QiitaTokenSecret.GrantRead(function, nil)
	data.MeIDParameter.GrantRead(function)
	data.ZennUsernameParameter.GrantRead(function)
	data.LogLevelParameter.GrantRead(function)

	functionURL := function.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_AWS_IAM,
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("ApiFunctionNameOutput"), &awscdk.CfnOutputProps{
		Value:      function.FunctionName(),
		ExportName: _jsii_.String(cfg.ExportName("api", "function-name")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("ApiFunctionArnOutput"), &awscdk.CfnOutputProps{
		Value:      function.FunctionArn(),
		ExportName: _jsii_.String(cfg.ExportName("api", "function-arn")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("ApiFunctionUrlOutput"), &awscdk.CfnOutputProps{
		Value:      functionURL.Url(),
		ExportName: _jsii_.String(cfg.ExportName("api", "function-url")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("ApiLogGroupNameOutput"), &awscdk.CfnOutputProps{
		Value:      logGroup.LogGroupName(),
		ExportName: _jsii_.String(cfg.ExportName("api", "log-group-name")),
	})

	apiStack.Function = function
	apiStack.FunctionURL = functionURL
	apiStack.LogGroup = logGroup

	return apiStack
}

func repositoryRoot() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".."))
}

func placeholderAPIAssetPath() string {
	return filepath.Join(repositoryRoot(), ".artifacts", "api", "bootstrap.zip")
}
