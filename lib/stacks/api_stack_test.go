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
	template.ResourceCountIs(_jsii_.String("AWS::Lambda::Url"), _jsii_.Number(1))
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
				"QIITA_TOKEN":         assertions.Match_AnyValue(),
				"ME_ID":               assertions.Match_AnyValue(),
				"ZENN_USERNAME":       assertions.Match_AnyValue(),
				"LOG_LEVEL":           assertions.Match_AnyValue(),
			}),
		},
	})

	template.HasResourceProperties(_jsii_.String("AWS::Lambda::Url"), map[string]interface{}{
		"AuthType": "AWS_IAM",
	})

	template.HasResourceProperties(_jsii_.String("AWS::Logs::LogGroup"), map[string]interface{}{
		"LogGroupName":    "/aws/lambda/me-dev-api",
		"RetentionInDays": 30,
	})

	template.HasOutput(_jsii_.String("ApiFunctionNameOutput"), map[string]interface{}{
		"Export": map[string]interface{}{
			"Name": "me-dev-api:function-name",
		},
	})

	template.HasOutput(_jsii_.String("ApiFunctionUrlOutput"), map[string]interface{}{
		"Export": map[string]interface{}{
			"Name": "me-dev-api:function-url",
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
