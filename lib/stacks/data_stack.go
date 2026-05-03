package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v10"
	_jsii_ "github.com/aws/jsii-runtime-go"
)

type DataStack struct {
	awscdk.Stack
	Table awsdynamodb.Table
}

func NewDataStack(scope constructs.Construct, id string, props *StackProps) *DataStack {
	stack := newStack(scope, id, props)
	dataStack := &DataStack{Stack: stack}

	if props == nil || props.Config == nil {
		return dataStack
	}

	cfg := props.Config

	table := awsdynamodb.NewTable(stack, _jsii_.String("Table"), &awsdynamodb.TableProps{
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
		PartitionKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("PK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TimeToLiveAttribute: _jsii_.String("ttl"),
		RemovalPolicy:       awscdk.RemovalPolicy_DESTROY,
	})

	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: _jsii_.String("GSI1"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("GSI1PK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("GSI1SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: _jsii_.String("GSI2"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("GSI2PK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("GSI2SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: _jsii_.String("GSI3"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("GSI3PK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("GSI3SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: _jsii_.String("GSI_EMAIL"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("GSI_EMAIL_PK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: _jsii_.String("SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("TableNameOutput"), &awscdk.CfnOutputProps{
		Value:      table.TableName(),
		ExportName: _jsii_.String(cfg.ExportName("data", "table-name")),
	})

	awscdk.NewCfnOutput(stack, _jsii_.String("TableArnOutput"), &awscdk.CfnOutputProps{
		Value:      table.TableArn(),
		ExportName: _jsii_.String(cfg.ExportName("data", "table-arn")),
	})

	dataStack.Table = table

	return dataStack
}
