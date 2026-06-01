package config

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	_jsii_ "github.com/aws/jsii-runtime-go"
)

func TestLoadDefaults(t *testing.T) {
	defer _jsii_.Close()

	app := awscdk.NewApp(nil)
	cfg := Load(app)

	if cfg.AppName != AppName {
		t.Fatalf("AppName = %q, want %q", cfg.AppName, AppName)
	}

	if cfg.Environment.Name != "dev" {
		t.Fatalf("Environment.Name = %q, want %q", cfg.Environment.Name, "dev")
	}

	if cfg.Environment.Prefix != "me-dev" {
		t.Fatalf("Environment.Prefix = %q, want %q", cfg.Environment.Prefix, "me-dev")
	}

	if cfg.Data.MeID != "replace-me" {
		t.Fatalf("Data.MeID = %q, want %q", cfg.Data.MeID, "replace-me")
	}

	if cfg.Data.ZennUsername != "replace-me" {
		t.Fatalf("Data.ZennUsername = %q, want %q", cfg.Data.ZennUsername, "replace-me")
	}

	if cfg.Data.LogLevel != "info" {
		t.Fatalf("Data.LogLevel = %q, want %q", cfg.Data.LogLevel, "info")
	}

	if cfg.Data.JWTSecret == "" {
		t.Fatal("Data.JWTSecret should be auto-generated when unset")
	}

	if cfg.Data.QiitaToken != "replace-me" {
		t.Fatalf("Data.QiitaToken = %q, want %q", cfg.Data.QiitaToken, "replace-me")
	}

	if cfg.Domain.AppDomain != "" {
		t.Fatalf("Domain.AppDomain = %q, want empty", cfg.Domain.AppDomain)
	}

	if cfg.Domain.APIDomain != "" {
		t.Fatalf("Domain.APIDomain = %q, want empty", cfg.Domain.APIDomain)
	}

	if cfg.HasFrontendCustomDomain() {
		t.Fatal("HasFrontendCustomDomain() = true, want false")
	}

	if cfg.HasAPICustomDomain() {
		t.Fatal("HasAPICustomDomain() = true, want false")
	}

	if cfg.StackName("data") != "me-dev-data" {
		t.Fatalf("StackName(data) = %q, want %q", cfg.StackName("data"), "me-dev-data")
	}

	if cfg.APIFunctionName() != "me-dev-api" {
		t.Fatalf("APIFunctionName() = %q, want %q", cfg.APIFunctionName(), "me-dev-api")
	}

	if cfg.APILogGroupName() != "/aws/lambda/me-dev-api" {
		t.Fatalf("APILogGroupName() = %q, want %q", cfg.APILogGroupName(), "/aws/lambda/me-dev-api")
	}

	if cfg.ExportName("data", "table-name") != "me-dev-data:table-name" {
		t.Fatalf("ExportName(data, table-name) = %q, want %q", cfg.ExportName("data", "table-name"), "me-dev-data:table-name")
	}

	if cfg.EnvironmentRef() != nil {
		t.Fatal("EnvironmentRef() should be nil when account or region is unset")
	}
}

func TestLoadUsesContextOverrides(t *testing.T) {
	defer _jsii_.Close()

	context := map[string]any{
		"envName":           "prod",
		"account":           "123456789012",
		"region":            "ap-northeast-1",
		"prefix":            "me-prod",
		"meId":              "umeki",
		"zennUsername":      "umekikazuya",
		"logLevel":          "debug",
		"jwtSecret":         "test-jwt-secret",
		"qiitaToken":        "test-qiita-token",
		"appDomain":         "www.example.com",
		"appCertificateArn": "arn:aws:acm:us-east-1:123456789012:certificate/frontend",
		"apiDomain":         "api.example.com",
		"apiCertificateArn": "arn:aws:acm:ap-northeast-1:123456789012:certificate/api",
	}

	app := awscdk.NewApp(&awscdk.AppProps{Context: &context})
	cfg := Load(app)

	if cfg.Environment.Name != "prod" {
		t.Fatalf("Environment.Name = %q, want %q", cfg.Environment.Name, "prod")
	}

	if cfg.Environment.Account != "123456789012" {
		t.Fatalf("Environment.Account = %q, want %q", cfg.Environment.Account, "123456789012")
	}

	if cfg.Environment.Region != "ap-northeast-1" {
		t.Fatalf("Environment.Region = %q, want %q", cfg.Environment.Region, "ap-northeast-1")
	}

	if cfg.Environment.Prefix != "me-prod" {
		t.Fatalf("Environment.Prefix = %q, want %q", cfg.Environment.Prefix, "me-prod")
	}

	if cfg.Data.MeID != "umeki" {
		t.Fatalf("Data.MeID = %q, want %q", cfg.Data.MeID, "umeki")
	}

	if cfg.Data.ZennUsername != "umekikazuya" {
		t.Fatalf("Data.ZennUsername = %q, want %q", cfg.Data.ZennUsername, "umekikazuya")
	}

	if cfg.Data.LogLevel != "debug" {
		t.Fatalf("Data.LogLevel = %q, want %q", cfg.Data.LogLevel, "debug")
	}

	if cfg.Data.JWTSecret != "test-jwt-secret" {
		t.Fatalf("Data.JWTSecret = %q, want %q", cfg.Data.JWTSecret, "test-jwt-secret")
	}

	if cfg.Data.QiitaToken != "test-qiita-token" {
		t.Fatalf("Data.QiitaToken = %q, want %q", cfg.Data.QiitaToken, "test-qiita-token")
	}

	if cfg.Domain.AppDomain != "www.example.com" {
		t.Fatalf("Domain.AppDomain = %q, want %q", cfg.Domain.AppDomain, "www.example.com")
	}

	if cfg.Domain.APIDomain != "api.example.com" {
		t.Fatalf("Domain.APIDomain = %q, want %q", cfg.Domain.APIDomain, "api.example.com")
	}

	if !cfg.HasFrontendCustomDomain() {
		t.Fatal("HasFrontendCustomDomain() = false, want true")
	}

	if !cfg.HasAPICustomDomain() {
		t.Fatal("HasAPICustomDomain() = false, want true")
	}

	if cfg.FrontendURL() != "https://www.example.com" {
		t.Fatalf("FrontendURL() = %q, want %q", cfg.FrontendURL(), "https://www.example.com")
	}

	if cfg.APIURL() != "https://api.example.com" {
		t.Fatalf("APIURL() = %q, want %q", cfg.APIURL(), "https://api.example.com")
	}

	if cfg.APIFunctionName() != "me-prod-api" {
		t.Fatalf("APIFunctionName() = %q, want %q", cfg.APIFunctionName(), "me-prod-api")
	}

	env := cfg.EnvironmentRef()
	if env == nil {
		t.Fatal("EnvironmentRef() returned nil")
	}

	if got := *env.Account; got != "123456789012" {
		t.Fatalf("EnvironmentRef().Account = %q, want %q", got, "123456789012")
	}

	if got := *env.Region; got != "ap-northeast-1" {
		t.Fatalf("EnvironmentRef().Region = %q, want %q", got, "ap-northeast-1")
	}
}

func TestLoadPanicsOnIncompleteDomainPair(t *testing.T) {
	defer _jsii_.Close()

	context := map[string]any{
		"apiDomain": "api.example.com",
	}

	app := awscdk.NewApp(&awscdk.AppProps{Context: &context})

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Load() should panic when apiDomain is set without apiCertificateArn")
		}
	}()

	Load(app)
}
