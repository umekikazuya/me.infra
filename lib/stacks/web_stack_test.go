package stacks

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	_jsii_ "github.com/aws/jsii-runtime-go"
	"github.com/umekikazuya/me.infra/lib/config"
)

func TestWebStack(t *testing.T) {
	defer _jsii_.Close()

	ensurePlaceholderAPIArtifact(t)

	context := map[string]any{
		"appDomain":         "www.example.com",
		"appCertificateArn": "arn:aws:acm:us-east-1:123456789012:certificate/frontend",
	}
	app := awscdk.NewApp(&awscdk.AppProps{Context: &context})
	cfg := config.Load(app)

	stack := NewWebStack(app, cfg.StackName("web"), &StackProps{
		StackProps: cfg.StackProps(),
		Config:     cfg,
	})

	template := assertions.Template_FromStack(stack.Stack, nil)

	template.ResourceCountIs(_jsii_.String("AWS::S3::Bucket"), _jsii_.Number(1))
	template.ResourceCountIs(_jsii_.String("AWS::CloudFront::Distribution"), _jsii_.Number(1))
	template.ResourceCountIs(_jsii_.String("AWS::CloudFront::Function"), _jsii_.Number(1))
	template.ResourceCountIs(_jsii_.String("AWS::CloudFront::OriginAccessControl"), _jsii_.Number(1))

	template.HasResourceProperties(_jsii_.String("AWS::S3::Bucket"), map[string]any{
		"PublicAccessBlockConfiguration": map[string]any{
			"BlockPublicAcls":       true,
			"BlockPublicPolicy":     true,
			"IgnorePublicAcls":      true,
			"RestrictPublicBuckets": true,
		},
	})

	template.HasResource(_jsii_.String("AWS::S3::Bucket"), map[string]any{
		"DeletionPolicy":      "Delete",
		"UpdateReplacePolicy": "Delete",
	})

	template.HasResourceProperties(_jsii_.String("AWS::CloudFront::Function"), map[string]any{
		"AutoPublish": true,
		"FunctionConfig": map[string]any{
			"Comment": "Rewrite SPA routes to /index.html while keeping static asset requests intact.",
			"Runtime": "cloudfront-js-2.0",
		},
	})

	template.HasResourceProperties(_jsii_.String("AWS::CloudFront::Distribution"), map[string]any{
		"DistributionConfig": assertions.Match_ObjectLike(&map[string]any{
			"Aliases": []any{
				"www.example.com",
			},
			"DefaultRootObject": "index.html",
			"PriceClass":        "PriceClass_200",
			"DefaultCacheBehavior": assertions.Match_ObjectLike(&map[string]any{
				"ViewerProtocolPolicy": "redirect-to-https",
				"FunctionAssociations": assertions.Match_ArrayWith(&[]any{
					assertions.Match_ObjectLike(&map[string]any{
						"EventType":   "viewer-request",
						"FunctionARN": assertions.Match_AnyValue(),
					}),
				}),
			}),
		}),
	})

	template.HasOutput(_jsii_.String("FrontendBucketNameOutput"), map[string]any{
		"Export": map[string]any{
			"Name": "me-dev-web:bucket-name",
		},
	})

	template.HasOutput(_jsii_.String("FrontendDistributionIdOutput"), map[string]any{
		"Export": map[string]any{
			"Name": "me-dev-web:distribution-id",
		},
	})

	template.HasOutput(_jsii_.String("FrontendDistributionDomainNameOutput"), map[string]any{
		"Export": map[string]any{
			"Name": "me-dev-web:distribution-domain-name",
		},
	})

	template.HasOutput(_jsii_.String("FrontendEndpointOutput"), map[string]any{
		"Value": "https://www.example.com",
		"Export": map[string]any{
			"Name": "me-dev-web:endpoint",
		},
	})
}
