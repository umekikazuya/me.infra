package config

import (
	"fmt"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	_jsii_ "github.com/aws/jsii-runtime-go"
)

const AppName = "me"

type EnvironmentConfig struct {
	Name    string
	Account string
	Region  string
	Prefix  string
	Tags    map[string]string
}

type AppConfig struct {
	AppName     string
	Environment EnvironmentConfig
}

func Load(app awscdk.App) *AppConfig {
	envName := contextString(app, "envName", "dev")
	account := contextString(app, "account", os.Getenv("CDK_DEFAULT_ACCOUNT"))
	region := contextString(app, "region", os.Getenv("CDK_DEFAULT_REGION"))
	prefix := contextString(app, "prefix", fmt.Sprintf("%s-%s", AppName, envName))

	return &AppConfig{
		AppName: AppName,
		Environment: EnvironmentConfig{
			Name:    envName,
			Account: account,
			Region:  region,
			Prefix:  prefix,
			Tags: map[string]string{
				"Application": AppName,
				"Environment": envName,
				"ManagedBy":   "aws-cdk",
			},
		},
	}
}

func (c *AppConfig) StackName(name string) string {
	return fmt.Sprintf("%s-%s", c.Environment.Prefix, name)
}

func (c *AppConfig) EnvironmentRef() *awscdk.Environment {
	if c.Environment.Account == "" || c.Environment.Region == "" {
		return nil
	}

	return &awscdk.Environment{
		Account: _jsii_.String(c.Environment.Account),
		Region:  _jsii_.String(c.Environment.Region),
	}
}

func (c *AppConfig) StackProps() awscdk.StackProps {
	props := awscdk.StackProps{}
	if env := c.EnvironmentRef(); env != nil {
		props.Env = env
	}

	return props
}

func (c *AppConfig) ApplyTags(scope constructs.IConstruct) {
	for key, value := range c.Environment.Tags {
		awscdk.Tags_Of(scope).Add(_jsii_.String(key), _jsii_.String(value), nil)
	}
}

func contextString(app awscdk.App, key, defaultValue string) string {
	value := app.Node().TryGetContext(_jsii_.String(key))
	if text, ok := value.(string); ok && text != "" {
		return text
	}

	return defaultValue
}
