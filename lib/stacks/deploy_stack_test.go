package stacks

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	_jsii_ "github.com/aws/jsii-runtime-go"
	"github.com/umekikazuya/me.infra/lib/config"
)

func TestDeployStack(t *testing.T) {
	defer _jsii_.Close()

	ensurePlaceholderAPIArtifact(t)

	context := map[string]any{
		"githubRepo":       "umekikazuya/me",
		"githubRefPattern": "refs/heads/*",
		"deployRoleName":   "me-app-deploy",
	}
	app := awscdk.NewApp(&awscdk.AppProps{Context: &context})
	cfg := config.Load(app)

	dataStack := NewDataStack(app, cfg.StackName("data"), &StackProps{
		StackProps: cfg.StackProps(),
		Config:     cfg,
	})

	apiStack := NewApiStack(app, cfg.StackName("api"), &StackProps{
		StackProps: cfg.StackProps(),
		Config:     cfg,
		Data:       dataStack,
	})

	webStack := NewWebStack(app, cfg.StackName("web"), &StackProps{
		StackProps: cfg.StackProps(),
		Config:     cfg,
	})

	stack := NewDeployStack(app, cfg.StackName("deploy"), &StackProps{
		StackProps: cfg.StackProps(),
		Config:     cfg,
		Api:        apiStack,
		Web:        webStack,
	})

	template := assertions.Template_FromStack(stack.Stack, nil)

	template.HasResourceProperties(_jsii_.String("AWS::IAM::Role"), map[string]any{
		"RoleName": "me-app-deploy",
		"AssumeRolePolicyDocument": assertions.Match_ObjectLike(&map[string]any{
			"Statement": assertions.Match_ArrayWith(&[]any{
				assertions.Match_ObjectLike(&map[string]any{
					"Action": "sts:AssumeRoleWithWebIdentity",
					"Condition": map[string]any{
						"StringEquals": map[string]any{
							"token.actions.githubusercontent.com:aud": "sts.amazonaws.com",
						},
						"StringLike": map[string]any{
							"token.actions.githubusercontent.com:sub": "repo:umekikazuya/me:ref:refs/heads/*",
						},
					},
				}),
			}),
		}),
	})

	template.HasResourceProperties(_jsii_.String("AWS::IAM::Policy"), map[string]any{
		"PolicyDocument": assertions.Match_ObjectLike(&map[string]any{
			"Statement": assertions.Match_ArrayWith(&[]any{
				assertions.Match_ObjectLike(&map[string]any{
					"Action": "lambda:UpdateFunctionCode",
				}),
				assertions.Match_ObjectLike(&map[string]any{
					"Action": []any{
						"s3:ListBucket",
						"s3:GetBucketLocation",
					},
				}),
				assertions.Match_ObjectLike(&map[string]any{
					"Action": []any{
						"s3:PutObject",
						"s3:DeleteObject",
					},
				}),
				assertions.Match_ObjectLike(&map[string]any{
					"Action": "cloudfront:CreateInvalidation",
				}),
			}),
		}),
	})

	template.HasOutput(_jsii_.String("DeployRoleArnOutput"), map[string]any{
		"Export": map[string]any{
			"Name": "me-dev-deploy:role-arn",
		},
	})

	template.HasOutput(_jsii_.String("DeployGitHubSubjectPatternOutput"), map[string]any{
		"Value": "repo:umekikazuya/me:ref:refs/heads/*",
		"Export": map[string]any{
			"Name": "me-dev-deploy:github-subject-pattern",
		},
	})
}
