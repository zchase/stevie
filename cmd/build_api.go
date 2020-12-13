package cmd

import (
	"context"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
	"github.com/zchase/stevie/pkg/application"
	"github.com/zchase/stevie/pkg/auto_pulumi"
)

func CreateAPIDeployment(environment string) (auto.Stack, error) {
	context := context.Background()

	// Check that user is authenticated.
	username, err := auto_pulumi.GetCurrentPulumiUser()
	if err != nil {
		return auto.Stack{}, nil
	}

	// Read in the base config file.
	appConfig, err := ReadBaseConfig(application.ApplicationConfigPath)
	if err != nil {
		return auto.Stack{}, err
	}

	// Create the deploy function.
	deployFunc := func(ctx *pulumi.Context) error {
		return application.BuildAPIRoutes(ctx, appConfig.DashCaseName, environment, appConfig.Routes)
	}

	// Create the stack.
	stackName := auto.FullyQualifiedStackName(username, appConfig.DashCaseName, environment)
	stack, err := auto.UpsertStackInlineSource(context, stackName, appConfig.DashCaseName, deployFunc)
	if err != nil {
		return auto.Stack{}, err
	}

	return stack, nil
}
