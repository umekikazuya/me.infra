package config

import (
	"crypto/rand"
	"encoding/base64"
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

type DataConfig struct {
	MeID         string
	ZennUsername string
	LogLevel     string
	JWTSecret    string
	QiitaToken   string
}

type DomainConfig struct {
	AppDomain         string
	AppCertificateARN string
	APIDomain         string
	APICertificateARN string
}

type AppConfig struct {
	AppName     string
	Environment EnvironmentConfig
	Data        DataConfig
	Domain      DomainConfig
}

func Load(app awscdk.App) *AppConfig {
	envName := contextString(app, "envName", "dev")
	account := contextString(app, "account", os.Getenv("CDK_DEFAULT_ACCOUNT"))
	region := contextString(app, "region", os.Getenv("CDK_DEFAULT_REGION"))
	prefix := contextString(app, "prefix", fmt.Sprintf("%s-%s", AppName, envName))
	meID := contextString(app, "meId", "replace-me")
	zennUsername := contextString(app, "zennUsername", "replace-me")
	logLevel := contextString(app, "logLevel", "info")
	jwtSecret := contextString(app, "jwtSecret", generateJWTSecret())
	qiitaToken := contextString(app, "qiitaToken", "replace-me")
	appDomain := contextString(app, "appDomain", "")
	appCertificateARN := contextString(app, "appCertificateArn", "")
	apiDomain := contextString(app, "apiDomain", "")
	apiCertificateARN := contextString(app, "apiCertificateArn", "")

	validateDomainPair("appDomain", appDomain, "appCertificateArn", appCertificateARN)
	validateDomainPair("apiDomain", apiDomain, "apiCertificateArn", apiCertificateARN)

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
		Data: DataConfig{
			MeID:         meID,
			ZennUsername: zennUsername,
			LogLevel:     logLevel,
			JWTSecret:    jwtSecret,
			QiitaToken:   qiitaToken,
		},
		Domain: DomainConfig{
			AppDomain:         appDomain,
			AppCertificateARN: appCertificateARN,
			APIDomain:         apiDomain,
			APICertificateARN: apiCertificateARN,
		},
	}
}

func (c *AppConfig) StackName(name string) string {
	return fmt.Sprintf("%s-%s", c.Environment.Prefix, name)
}

func (c *AppConfig) APIFunctionName() string {
	return fmt.Sprintf("%s-%s-api", c.AppName, c.Environment.Name)
}

func (c *AppConfig) APILogGroupName() string {
	return fmt.Sprintf("/aws/lambda/%s", c.APIFunctionName())
}

func (c *AppConfig) ExportName(stackName, exportName string) string {
	return fmt.Sprintf("%s:%s", c.StackName(stackName), exportName)
}

func (c *AppConfig) HasFrontendCustomDomain() bool {
	return c.Domain.AppDomain != "" && c.Domain.AppCertificateARN != ""
}

func (c *AppConfig) HasAPICustomDomain() bool {
	return c.Domain.APIDomain != "" && c.Domain.APICertificateARN != ""
}

func (c *AppConfig) FrontendURL() string {
	if !c.HasFrontendCustomDomain() {
		return ""
	}

	return "https://" + c.Domain.AppDomain
}

func (c *AppConfig) APIURL() string {
	if !c.HasAPICustomDomain() {
		return ""
	}

	return "https://" + c.Domain.APIDomain
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

func generateJWTSecret() string {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Sprintf("generate jwt secret: %v", err))
	}

	return base64.RawURLEncoding.EncodeToString(buf)
}

func validateDomainPair(domainKey, domainValue, certKey, certValue string) {
	if (domainValue == "") == (certValue == "") {
		return
	}

	panic(fmt.Sprintf("%s and %s must be provided together", domainKey, certKey))
}
