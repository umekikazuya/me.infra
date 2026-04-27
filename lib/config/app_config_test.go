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

	if cfg.StackName("data") != "me-dev-data" {
		t.Fatalf("StackName(data) = %q, want %q", cfg.StackName("data"), "me-dev-data")
	}

	if cfg.EnvironmentRef() != nil {
		t.Fatal("EnvironmentRef() should be nil when account or region is unset")
	}
}

func TestLoadUsesContextOverrides(t *testing.T) {
	defer _jsii_.Close()

	context := map[string]interface{}{
		"envName": "prod",
		"account": "123456789012",
		"region":  "ap-northeast-1",
		"prefix":  "me-prod",
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
