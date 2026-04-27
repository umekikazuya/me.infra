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

	app := awscdk.NewApp(nil)
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

	stack := NewWebStack(app, cfg.StackName("web"), &StackProps{
		StackProps: cfg.StackProps(),
		Config:     cfg,
		Api:        apiStack,
		Data:       dataStack,
	})

	template := assertions.Template_FromStack(stack.Stack, nil)

	template.ResourceCountIs(_jsii_.String("AWS::S3::Bucket"), _jsii_.Number(1))
	template.ResourceCountIs(_jsii_.String("AWS::CloudFront::Distribution"), _jsii_.Number(1))
	template.ResourceCountIs(_jsii_.String("AWS::CloudFront::Function"), _jsii_.Number(1))
	template.ResourceCountIs(_jsii_.String("AWS::CloudFront::OriginAccessControl"), _jsii_.Number(2))

	template.HasResourceProperties(_jsii_.String("AWS::S3::Bucket"), map[string]interface{}{
		"PublicAccessBlockConfiguration": map[string]interface{}{
			"BlockPublicAcls":       true,
			"BlockPublicPolicy":     true,
			"IgnorePublicAcls":      true,
			"RestrictPublicBuckets": true,
		},
	})

	template.HasResourceProperties(_jsii_.String("AWS::CloudFront::Function"), map[string]interface{}{
		"AutoPublish": true,
		"FunctionConfig": map[string]interface{}{
			"Comment": "Rewrite SPA routes to /index.html while keeping static asset requests intact.",
			"Runtime": "cloudfront-js-2.0",
		},
	})

	template.HasResourceProperties(_jsii_.String("AWS::CloudFront::Distribution"), map[string]interface{}{
		"DistributionConfig": assertions.Match_ObjectLike(&map[string]interface{}{
			"DefaultRootObject": "index.html",
			"DefaultCacheBehavior": assertions.Match_ObjectLike(&map[string]interface{}{
				"ViewerProtocolPolicy": "redirect-to-https",
				"FunctionAssociations": assertions.Match_ArrayWith(&[]interface{}{
					assertions.Match_ObjectLike(&map[string]interface{}{
						"EventType":   "viewer-request",
						"FunctionARN": assertions.Match_AnyValue(),
					}),
				}),
			}),
			"CacheBehaviors": assertions.Match_ArrayWith(&[]interface{}{
				assertions.Match_ObjectLike(&map[string]interface{}{
					"PathPattern":           "/api/*",
					"ViewerProtocolPolicy":  "redirect-to-https",
					"CachePolicyId":         assertions.Match_AnyValue(),
					"OriginRequestPolicyId": assertions.Match_AnyValue(),
				}),
			}),
		}),
	})

	template.HasOutput(_jsii_.String("FrontendBucketNameOutput"), map[string]interface{}{
		"Export": map[string]interface{}{
			"Name": "me-dev-web:bucket-name",
		},
	})

	template.HasOutput(_jsii_.String("FrontendDistributionIdOutput"), map[string]interface{}{
		"Export": map[string]interface{}{
			"Name": "me-dev-web:distribution-id",
		},
	})
}
