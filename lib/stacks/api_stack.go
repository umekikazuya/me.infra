package stacks

import (
	"path/filepath"
	"runtime"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/constructs-go/constructs/v10"
	_jsii_ "github.com/aws/jsii-runtime-go"
)

type ApiStack struct {
	awscdk.Stack
	Function   awslambda.Function
	HTTPAPI    awsapigatewayv2.HttpApi
	LogGroup   awslogs.LogGroup
	DomainName awsapigatewayv2.DomainName
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
		"JWT_SECRET":          _jsii_.String(cfg.Data.JWTSecret),
		"QIITA_TOKEN":         _jsii_.String(cfg.Data.QiitaToken),
		"ME_ID":               _jsii_.String(cfg.Data.MeID),
		"ZENN_USERNAME":       _jsii_.String(cfg.Data.ZennUsername),
		"LOG_LEVEL":           _jsii_.String(cfg.Data.LogLevel),
	}

	if cfg.HasFrontendCustomDomain() {
		environment["CORS_ALLOWED_ORIGINS"] = _jsii_.String(cfg.FrontendURL())
	}

	logGroup := awslogs.NewLogGroup(stack, _jsii_.String("ApiLogGroup"), &awslogs.LogGroupProps{
		LogGroupName:  _jsii_.String(cfg.APILogGroupName()),
		Retention:     awslogs.RetentionDays_ONE_MONTH,
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	function := awslambda.NewFunction(
		stack,
		_jsii_.String("ApiFunction"),
		&awslambda.FunctionProps{
			FunctionName: _jsii_.String(cfg.APIFunctionName()),
			Description:  _jsii_.String("Placeholder me API Lambda. The app repository replaces the code with update-function-code after initial deployment."),
			Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
			Architecture: awslambda.Architecture_ARM_64(),
			Handler:      _jsii_.String("bootstrap"),
			Code:         awslambda.Code_FromAsset(_jsii_.String(placeholderAPIAssetPath()), &awss3assets.AssetOptions{DeployTime: _jsii_.Bool(true)}),
			Environment:  &environment,
			MemorySize:   _jsii_.Number(512),
			Timeout:      awscdk.Duration_Seconds(_jsii_.Number(30)),
			LogGroup:     logGroup,
		},
	)

	data.Table.GrantReadWriteData(function)

	var corsPreflight *awsapigatewayv2.CorsPreflightOptions
	if cfg.HasFrontendCustomDomain() {
		corsPreflight = &awsapigatewayv2.CorsPreflightOptions{
			AllowCredentials: _jsii_.Bool(true),
			AllowHeaders: &[]*string{
				_jsii_.String("Accept"),
				_jsii_.String("Content-Type"),
				_jsii_.String("X-Request-ID"),
				_jsii_.String("X-Requested-With"),
			},
			ExposeHeaders: &[]*string{
				_jsii_.String("X-Request-ID"),
			},
			AllowMethods: &[]awsapigatewayv2.CorsHttpMethod{
				awsapigatewayv2.CorsHttpMethod_GET,
				awsapigatewayv2.CorsHttpMethod_HEAD,
				awsapigatewayv2.CorsHttpMethod_OPTIONS,
				awsapigatewayv2.CorsHttpMethod_POST,
				awsapigatewayv2.CorsHttpMethod_PUT,
				awsapigatewayv2.CorsHttpMethod_PATCH,
				awsapigatewayv2.CorsHttpMethod_DELETE,
			},
			AllowOrigins: &[]*string{
				_jsii_.String(cfg.FrontendURL()),
			},
			MaxAge: awscdk.Duration_Hours(_jsii_.Number(1)),
		}
	}

	var domainName awsapigatewayv2.DomainName
	var defaultDomainMapping *awsapigatewayv2.DomainMappingOptions
	if cfg.HasAPICustomDomain() {
		certificate := awscertificatemanager.Certificate_FromCertificateArn(
			stack,
			_jsii_.String("ApiCertificate"),
			_jsii_.String(cfg.Domain.APICertificateARN),
		)
		domainName = awsapigatewayv2.NewDomainName(stack, _jsii_.String("ApiCustomDomain"), &awsapigatewayv2.DomainNameProps{
			Certificate: certificate,
			DomainName:  _jsii_.String(cfg.Domain.APIDomain),
		})
		defaultDomainMapping = &awsapigatewayv2.DomainMappingOptions{
			DomainName: domainName,
		}
	}

	httpAPI := awsapigatewayv2.NewHttpApi(stack, _jsii_.String("HttpApi"), &awsapigatewayv2.HttpApiProps{
		ApiName:                   _jsii_.String(cfg.StackName("http-api")),
		CorsPreflight:             corsPreflight,
		DefaultDomainMapping:      defaultDomainMapping,
		Description:               _jsii_.String("me API Gateway HTTP API"),
		DisableExecuteApiEndpoint: _jsii_.Bool(cfg.HasAPICustomDomain()),
	})

	integration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		_jsii_.String("DefaultIntegration"),
		function,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{
			PayloadFormatVersion: awsapigatewayv2.PayloadFormatVersion_VERSION_2_0(),
		},
	)
	nonOptionsMethods := &[]awsapigatewayv2.HttpMethod{
		awsapigatewayv2.HttpMethod_GET,
		awsapigatewayv2.HttpMethod_HEAD,
		awsapigatewayv2.HttpMethod_POST,
		awsapigatewayv2.HttpMethod_PUT,
		awsapigatewayv2.HttpMethod_PATCH,
		awsapigatewayv2.HttpMethod_DELETE,
	}
	httpAPI.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:        _jsii_.String("/{proxy+}"),
		Methods:     nonOptionsMethods,
		Integration: integration,
	})

	apiEndpoint := _jsii_.String(cfg.APIURL())
	if !cfg.HasAPICustomDomain() {
		apiEndpoint = httpAPI.ApiEndpoint()
	}

	awscdk.NewCfnOutput(stack, _jsii_.String("ApiFunctionNameOutput"), &awscdk.CfnOutputProps{
		Value:      function.FunctionName(),
		ExportName: _jsii_.String(cfg.ExportName("api", "function-name")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("ApiFunctionArnOutput"), &awscdk.CfnOutputProps{
		Value:      function.FunctionArn(),
		ExportName: _jsii_.String(cfg.ExportName("api", "function-arn")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("ApiEndpointOutput"), &awscdk.CfnOutputProps{
		Value:      apiEndpoint,
		ExportName: _jsii_.String(cfg.ExportName("api", "endpoint")),
	})

	if domainName != nil {
		awscdk.NewCfnOutput(stack, _jsii_.String("ApiCustomDomainNameOutput"), &awscdk.CfnOutputProps{
			Value:      domainName.Name(),
			ExportName: _jsii_.String(cfg.ExportName("api", "custom-domain-name")),
		})

		awscdk.NewCfnOutput(stack, _jsii_.String("ApiCustomDomainRegionalNameOutput"), &awscdk.CfnOutputProps{
			Value:      domainName.RegionalDomainName(),
			ExportName: _jsii_.String(cfg.ExportName("api", "custom-domain-regional-name")),
		})

		awscdk.NewCfnOutput(stack, _jsii_.String("ApiCustomDomainRegionalHostedZoneIdOutput"), &awscdk.CfnOutputProps{
			Value:      domainName.RegionalHostedZoneId(),
			ExportName: _jsii_.String(cfg.ExportName("api", "custom-domain-regional-hosted-zone-id")),
		})
	}

	awscdk.NewCfnOutput(stack, _jsii_.String("ApiLogGroupNameOutput"), &awscdk.CfnOutputProps{
		Value:      logGroup.LogGroupName(),
		ExportName: _jsii_.String(cfg.ExportName("api", "log-group-name")),
	})

	apiStack.Function = function
	apiStack.HTTPAPI = httpAPI
	apiStack.LogGroup = logGroup
	apiStack.DomainName = domainName

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
