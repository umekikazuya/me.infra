package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	_jsii_ "github.com/aws/jsii-runtime-go"
)

type WebStack struct {
	awscdk.Stack
	Bucket          awss3.Bucket
	Distribution    awscloudfront.Distribution
	RewriteFunction awscloudfront.Function
}

func NewWebStack(scope constructs.Construct, id string, props *StackProps) *WebStack {
	if props == nil {
		panic("web stack requires stack props")
	}
	if props.Config == nil {
		panic("web stack requires app config")
	}

	stack := newStack(scope, id, props)
	webStack := &WebStack{Stack: stack}

	cfg := props.Config

	bucket := awss3.NewBucket(stack, _jsii_.String("FrontendBucket"), &awss3.BucketProps{
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		AutoDeleteObjects: _jsii_.Bool(true),
		EnforceSSL:        _jsii_.Bool(true),
		ObjectOwnership:   awss3.ObjectOwnership_BUCKET_OWNER_ENFORCED,
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
	})

	rewriteFunction := awscloudfront.NewFunction(stack, _jsii_.String("SpaRewriteFunction"), &awscloudfront.FunctionProps{
		Comment: _jsii_.String("Rewrite SPA routes to /index.html while keeping static asset requests intact."),
		Code: awscloudfront.FunctionCode_FromInline(_jsii_.String(`function handler(event) {
  var request = event.request;
  var uri = request.uri;

  if (uri.startsWith('/api/')) {
    return request;
  }

  if (uri === '/' || uri.endsWith('/')) {
    request.uri = '/index.html';
    return request;
  }

  var segments = uri.split('/');
  var basename = segments[segments.length - 1];

  if (!basename.includes('.')) {
    request.uri = '/index.html';
  }

  return request;
}`)),
		Runtime: awscloudfront.FunctionRuntime_JS_2_0(),
	})

	distributionProps := &awscloudfront.DistributionProps{
		Comment:           _jsii_.String("me frontend distribution"),
		DefaultRootObject: _jsii_.String("index.html"),
		PriceClass:        awscloudfront.PriceClass_PRICE_CLASS_200,
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			Origin:               awscloudfrontorigins.S3BucketOrigin_WithOriginAccessControl(bucket, nil),
			CachePolicy:          awscloudfront.CachePolicy_CACHING_OPTIMIZED(),
			Compress:             _jsii_.Bool(true),
			ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
			FunctionAssociations: &[]*awscloudfront.FunctionAssociation{
				{
					EventType: awscloudfront.FunctionEventType_VIEWER_REQUEST,
					Function:  rewriteFunction,
				},
			},
		},
	}
	if cfg.HasFrontendCustomDomain() {
		distributionProps.Certificate = awscertificatemanager.Certificate_FromCertificateArn(
			stack,
			_jsii_.String("FrontendCertificate"),
			_jsii_.String(cfg.Domain.AppCertificateARN),
		)
		distributionProps.DomainNames = &[]*string{
			_jsii_.String(cfg.Domain.AppDomain),
		}
	}

	distribution := awscloudfront.NewDistribution(stack, _jsii_.String("FrontendDistribution"), distributionProps)

	awscdk.NewCfnOutput(stack, _jsii_.String("FrontendBucketNameOutput"), &awscdk.CfnOutputProps{
		Value:      bucket.BucketName(),
		ExportName: _jsii_.String(cfg.ExportName("web", "bucket-name")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("FrontendDistributionIdOutput"), &awscdk.CfnOutputProps{
		Value:      distribution.DistributionId(),
		ExportName: _jsii_.String(cfg.ExportName("web", "distribution-id")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("FrontendDistributionDomainNameOutput"), &awscdk.CfnOutputProps{
		Value:      distribution.DistributionDomainName(),
		ExportName: _jsii_.String(cfg.ExportName("web", "distribution-domain-name")),
	})

	frontendEndpoint := awscdk.Fn_Join(_jsii_.String(""), &[]*string{
		_jsii_.String("https://"),
		distribution.DistributionDomainName(),
	})
	if cfg.HasFrontendCustomDomain() {
		frontendEndpoint = _jsii_.String(cfg.FrontendURL())
	}

	awscdk.NewCfnOutput(stack, _jsii_.String("FrontendEndpointOutput"), &awscdk.CfnOutputProps{
		Value:      frontendEndpoint,
		ExportName: _jsii_.String(cfg.ExportName("web", "endpoint")),
	})

	webStack.Bucket = bucket
	webStack.Distribution = distribution
	webStack.RewriteFunction = rewriteFunction

	return webStack
}
