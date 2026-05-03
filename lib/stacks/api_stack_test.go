package stacks

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	_jsii_ "github.com/aws/jsii-runtime-go"
	"github.com/umekikazuya/me.infra/lib/config"
)

func TestApiStack(t *testing.T) {
	defer _jsii_.Close()

	ensurePlaceholderAPIArtifact(t)

	app := awscdk.NewApp(nil)
	cfg := config.Load(app)
	dataStack := NewDataStack(app, cfg.StackName("data"), &StackProps{
		StackProps: cfg.StackProps(),
		Config:     cfg,
	})

	stack := NewApiStack(app, cfg.StackName("api"), &StackProps{
		StackProps: cfg.StackProps(),
		Config:     cfg,
		Data:       dataStack,
	})

	template := assertions.Template_FromStack(stack.Stack, nil)

	template.ResourceCountIs(_jsii_.String("AWS::Lambda::Function"), _jsii_.Number(1))
	template.ResourceCountIs(_jsii_.String("AWS::ApiGatewayV2::Api"), _jsii_.Number(1))
	template.ResourceCountIs(_jsii_.String("AWS::ApiGatewayV2::Stage"), _jsii_.Number(1))
	template.ResourceCountIs(_jsii_.String("AWS::ApiGatewayV2::Integration"), _jsii_.Number(1))
	template.ResourceCountIs(_jsii_.String("AWS::ApiGatewayV2::Route"), _jsii_.Number(1))
	template.ResourceCountIs(_jsii_.String("AWS::Lambda::Permission"), _jsii_.Number(1))
	template.ResourceCountIs(_jsii_.String("AWS::Logs::LogGroup"), _jsii_.Number(1))

	template.HasResourceProperties(_jsii_.String("AWS::Lambda::Function"), map[string]interface{}{
		"FunctionName": "me-dev-api",
		"Runtime":      "provided.al2023",
		"Architectures": []interface{}{
			"arm64",
		},
		"Handler": "bootstrap",
		"Environment": map[string]interface{}{
			"Variables": assertions.Match_ObjectLike(&map[string]interface{}{
				"DYNAMODB_TABLE_NAME": assertions.Match_AnyValue(),
				"JWT_SECRET":          assertions.Match_AnyValue(),
				"QIITA_TOKEN":         "replace-me",
				"ME_ID":               "replace-me",
				"ZENN_USERNAME":       "replace-me",
				"LOG_LEVEL":           "info",
			}),
		},
	})

	template.HasResourceProperties(_jsii_.String("AWS::ApiGatewayV2::Api"), map[string]interface{}{
		"Name":         "me-dev-http-api",
		"ProtocolType": "HTTP",
	})

	template.HasResourceProperties(_jsii_.String("AWS::Logs::LogGroup"), map[string]interface{}{
		"LogGroupName":    "/aws/lambda/me-dev-api",
		"RetentionInDays": 30,
	})

	template.HasResource(_jsii_.String("AWS::Logs::LogGroup"), map[string]interface{}{
		"DeletionPolicy":      "Delete",
		"UpdateReplacePolicy": "Delete",
	})

	template.HasOutput(_jsii_.String("ApiFunctionNameOutput"), map[string]interface{}{
		"Export": map[string]interface{}{
			"Name": "me-dev-api:function-name",
		},
	})

	template.HasOutput(_jsii_.String("ApiEndpointOutput"), map[string]interface{}{
		"Export": map[string]interface{}{
			"Name": "me-dev-api:endpoint",
		},
	})
}

func ensurePlaceholderAPIArtifact(t *testing.T) {
	t.Helper()

	root := repositoryRoot()
	cmd := exec.Command("sh", filepath.Join(root, "scripts", "package-placeholder-api.sh"))
	cmd.Dir = root
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("package placeholder api: %v\n%s", err, output)
	}
}
