package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	_jsii_ "github.com/aws/jsii-runtime-go"
	"github.com/umekikazuya/me.infra/lib/config"
)

type StackProps struct {
	awscdk.StackProps
	Config *config.AppConfig
}

func newStack(scope constructs.Construct, id string, props *StackProps) awscdk.Stack {
	stackProps := awscdk.StackProps{}
	var appConfig *config.AppConfig

	if props != nil {
		stackProps = props.StackProps
		appConfig = props.Config
	}

	stack := awscdk.NewStack(scope, _jsii_.String(id), &stackProps)
	if appConfig != nil {
		appConfig.ApplyTags(stack)
	}

	return stack
}
