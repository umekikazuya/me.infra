package stacks

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	_jsii_ "github.com/aws/jsii-runtime-go"
)

type DeployStack struct {
	awscdk.Stack
	Provider awsiam.OpenIdConnectProvider
	Role     awsiam.Role
}

func NewDeployStack(scope constructs.Construct, id string, props *StackProps) *DeployStack {
	if props == nil {
		panic("deploy stack requires stack props")
	}
	if props.Config == nil {
		panic("deploy stack requires app config")
	}
	if props.Api == nil || props.Api.Function == nil {
		panic("deploy stack requires api stack with function")
	}
	if props.Web == nil || props.Web.Bucket == nil || props.Web.Distribution == nil {
		panic("deploy stack requires web stack with bucket and distribution")
	}

	stack := newStack(scope, id, props)
	deployStack := &DeployStack{Stack: stack}
	cfg := props.Config

	provider := awsiam.NewOpenIdConnectProvider(stack, _jsii_.String("GitHubOidcProvider"), &awsiam.OpenIdConnectProviderProps{
		Url: _jsii_.String("https://token.actions.githubusercontent.com"),
		ClientIds: &[]*string{
			_jsii_.String("sts.amazonaws.com"),
		},
	})

	principal := awsiam.NewOpenIdConnectPrincipal(provider, &map[string]interface{}{
		"StringEquals": map[string]string{
			"token.actions.githubusercontent.com:aud": "sts.amazonaws.com",
		},
		"StringLike": map[string]string{
			"token.actions.githubusercontent.com:sub": cfg.GitHubOIDCSubjectPattern(),
		},
	})

	role := awsiam.NewRole(stack, _jsii_.String("AppDeployRole"), &awsiam.RoleProps{
		RoleName:    _jsii_.String(cfg.Deploy.RoleName),
		Description: _jsii_.String(fmt.Sprintf("Role assumed by GitHub Actions OIDC to deploy app artifacts for %s.", cfg.AppName)),
		AssumedBy:   principal,
	})

	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			_jsii_.String("lambda:UpdateFunctionCode"),
		},
		Resources: &[]*string{
			props.Api.Function.FunctionArn(),
		},
	}))

	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			_jsii_.String("s3:ListBucket"),
			_jsii_.String("s3:GetBucketLocation"),
		},
		Resources: &[]*string{
			props.Web.Bucket.BucketArn(),
		},
	}))

	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			_jsii_.String("s3:PutObject"),
			_jsii_.String("s3:DeleteObject"),
		},
		Resources: &[]*string{
			props.Web.Bucket.ArnForObjects(_jsii_.String("*")),
		},
	}))

	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			_jsii_.String("cloudfront:CreateInvalidation"),
		},
		Resources: &[]*string{
			cloudFrontDistributionArn(props.Web.Distribution.DistributionId()),
		},
	}))

	awscdk.NewCfnOutput(stack, _jsii_.String("DeployRoleArnOutput"), &awscdk.CfnOutputProps{
		Value:      role.RoleArn(),
		ExportName: _jsii_.String(cfg.ExportName("deploy", "role-arn")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("DeployRoleNameOutput"), &awscdk.CfnOutputProps{
		Value:      role.RoleName(),
		ExportName: _jsii_.String(cfg.ExportName("deploy", "role-name")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("DeployGitHubSubjectPatternOutput"), &awscdk.CfnOutputProps{
		Value:      _jsii_.String(cfg.GitHubOIDCSubjectPattern()),
		ExportName: _jsii_.String(cfg.ExportName("deploy", "github-subject-pattern")),
	})

	deployStack.Provider = provider
	deployStack.Role = role

	return deployStack
}

func cloudFrontDistributionArn(distributionID *string) *string {
	return awscdk.Fn_Join(_jsii_.String(""), &[]*string{
		_jsii_.String("arn:"),
		awscdk.Aws_PARTITION(),
		_jsii_.String(":cloudfront::"),
		awscdk.Aws_ACCOUNT_ID(),
		_jsii_.String(":distribution/"),
		distributionID,
	})
}
