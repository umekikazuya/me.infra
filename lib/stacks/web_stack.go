package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
)

func NewWebStack(scope constructs.Construct, id string, props *StackProps) awscdk.Stack {
	return newStack(scope, id, props)
}
