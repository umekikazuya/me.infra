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
	template.ResourceCountIs(_jsii_.String("AWS::SSM::Parameter"), _jsii_.Number(0))

	template.HasResourceProperties(_jsii_.String("AWS::DynamoDB::Table"), map[string]interface{}{
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

	template.ResourcePropertiesCountIs(_jsii_.String("AWS::DynamoDB::Table"), map[string]interface{}{
		"TableName": assertions.Match_AnyValue(),
	}, _jsii_.Number(0))

	template.HasResource(_jsii_.String("AWS::DynamoDB::Table"), map[string]interface{}{
		"DeletionPolicy":      "Delete",
		"UpdateReplacePolicy": "Delete",
	})

	template.HasOutput(_jsii_.String("TableNameOutput"), map[string]interface{}{
		"Export": map[string]interface{}{
			"Name": "me-dev-data:table-name",
		},
	})
}
