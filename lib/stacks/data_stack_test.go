package stacks

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	_jsii_ "github.com/aws/jsii-runtime-go"
	"github.com/umekikazuya/me.infra/lib/config"
)

func TestDataStack(t *testing.T) {
	defer _jsii_.Close()

	app := awscdk.NewApp(nil)
	cfg := config.Load(app)

	stack := NewDataStack(app, cfg.StackName("data"), &StackProps{
		StackProps: cfg.StackProps(),
		Config:     cfg,
	})

	template := assertions.Template_FromStack(stack.Stack, nil)

	template.ResourceCountIs(_jsii_.String("AWS::DynamoDB::Table"), _jsii_.Number(1))
	template.ResourceCountIs(_jsii_.String("AWS::SecretsManager::Secret"), _jsii_.Number(2))
	template.ResourceCountIs(_jsii_.String("AWS::SSM::Parameter"), _jsii_.Number(3))

	template.HasResourceProperties(_jsii_.String("AWS::DynamoDB::Table"), map[string]interface{}{
		"TableName":   "me.",
		"BillingMode": "PAY_PER_REQUEST",
		"TimeToLiveSpecification": map[string]interface{}{
			"AttributeName": "ttl",
			"Enabled":       true,
		},
		"GlobalSecondaryIndexes": assertions.Match_ArrayWith(&[]interface{}{
			map[string]interface{}{
				"IndexName": "GSI1",
				"KeySchema": assertions.Match_ArrayWith(&[]interface{}{
					map[string]interface{}{"AttributeName": "GSI1PK", "KeyType": "HASH"},
					map[string]interface{}{"AttributeName": "GSI1SK", "KeyType": "RANGE"},
				}),
				"Projection": map[string]interface{}{"ProjectionType": "ALL"},
			},
			map[string]interface{}{
				"IndexName": "GSI2",
				"KeySchema": assertions.Match_ArrayWith(&[]interface{}{
					map[string]interface{}{"AttributeName": "GSI2PK", "KeyType": "HASH"},
					map[string]interface{}{"AttributeName": "GSI2SK", "KeyType": "RANGE"},
				}),
				"Projection": map[string]interface{}{"ProjectionType": "ALL"},
			},
			map[string]interface{}{
				"IndexName": "GSI3",
				"KeySchema": assertions.Match_ArrayWith(&[]interface{}{
					map[string]interface{}{"AttributeName": "GSI3PK", "KeyType": "HASH"},
					map[string]interface{}{"AttributeName": "GSI3SK", "KeyType": "RANGE"},
				}),
				"Projection": map[string]interface{}{"ProjectionType": "ALL"},
			},
			map[string]interface{}{
				"IndexName": "GSI_EMAIL",
				"KeySchema": assertions.Match_ArrayWith(&[]interface{}{
					map[string]interface{}{"AttributeName": "GSI_EMAIL_PK", "KeyType": "HASH"},
					map[string]interface{}{"AttributeName": "SK", "KeyType": "RANGE"},
				}),
				"Projection": map[string]interface{}{"ProjectionType": "ALL"},
			},
		}),
	})

	template.HasResourceProperties(_jsii_.String("AWS::SecretsManager::Secret"), map[string]interface{}{
		"Name": "me/dev/jwtSecret",
	})

	template.HasResourceProperties(_jsii_.String("AWS::SecretsManager::Secret"), map[string]interface{}{
		"Name": "me/dev/qiitaToken",
	})

	template.HasResourceProperties(_jsii_.String("AWS::SSM::Parameter"), map[string]interface{}{
		"Name":  "/me/dev/meId",
		"Value": "replace-me",
	})

	template.HasResourceProperties(_jsii_.String("AWS::SSM::Parameter"), map[string]interface{}{
		"Name":  "/me/dev/zennUsername",
		"Value": "replace-me",
	})

	template.HasResourceProperties(_jsii_.String("AWS::SSM::Parameter"), map[string]interface{}{
		"Name":  "/me/dev/logLevel",
		"Value": "info",
	})

	template.HasOutput(_jsii_.String("TableNameOutput"), map[string]interface{}{
		"Export": map[string]interface{}{
			"Name": "me-dev-data:table-name",
		},
	})

	template.HasOutput(_jsii_.String("JwtSecretArnOutput"), map[string]interface{}{
		"Export": map[string]interface{}{
			"Name": "me-dev-data:jwt-secret-arn",
		},
	})
}
