package cmd

import (
	"context"
	"path"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
	"github.com/zchase/stevie/pkg/application"
	"github.com/zchase/stevie/pkg/auto_pulumi"
	"github.com/zchase/stevie/pkg/utils"
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
		modelDirPath := path.Join(application.ApplicationFolder, application.ModelsFolder)
		tableNames, err := auto_pulumi.BuildDynamoDBTables(ctx, modelDirPath, environment)
		if err != nil {
			return utils.NewErrorMessage("Error creating Dynamo tables", err)
		}

		return application.BuildAPIRoutes(ctx, appConfig.DashCaseName, environment, appConfig.Routes, tableNames)
	}

	// Create the stack.
	stackName := auto.FullyQualifiedStackName(username, appConfig.DashCaseName, environment)
	stack, err := auto.UpsertStackInlineSource(context, stackName, appConfig.DashCaseName, deployFunc)
	if err != nil {
		return auto.Stack{}, err
	}

	return stack, nil
}
