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

	dataStack := stacks.NewDataStack(app, cfg.StackName("data"), &stacks.StackProps{
		StackProps: cfg.StackProps(),
		Config:     cfg,
	})

	apiStack := stacks.NewApiStack(app, cfg.StackName("api"), &stacks.StackProps{
		StackProps: cfg.StackProps(),
		Config:     cfg,
		Data:       dataStack,
	})

	webStack := stacks.NewWebStack(app, cfg.StackName("web"), &stacks.StackProps{
		StackProps: cfg.StackProps(),
		Config:     cfg,
	})

	stacks.NewDeployStack(app, cfg.StackName("deploy"), &stacks.StackProps{
		StackProps: cfg.StackProps(),
		Config:     cfg,
		Api:        apiStack,
		Web:        webStack,
	})

	app.Synth(nil)
}
