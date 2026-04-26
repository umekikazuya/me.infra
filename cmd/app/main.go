package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
	"github.com/umekikazuya/me.infra/lib/config"
	"github.com/umekikazuya/me.infra/lib/stacks"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	cfg := config.Load(app)

	commonProps := &stacks.StackProps{
		StackProps: cfg.StackProps(),
		Config:     cfg,
	}

	stacks.NewDataStack(app, cfg.StackName("data"), commonProps)
	stacks.NewApiStack(app, cfg.StackName("api"), commonProps)
	stacks.NewWebStack(app, cfg.StackName("web"), commonProps)

	app.Synth(nil)
}
