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

		endpoints, err := application.BuildAPIRoutes(ctx, appConfig.DashCaseName, environment, appConfig.Routes, tableNames)
		if err != nil {
			return utils.NewErrorMessage("Error creating API Routes", err)
		}

		// If there is a UI directory, we should deploy it.
		dirExists, err := utils.DoesFileExist("ui")
		if err != nil {
			return utils.NewErrorMessage("Error finding UI directory", err)
		}

		if dirExists {
			err = auto_pulumi.CreateWebsiteFromDirectoryContents(ctx, endpoints, "ui", environment)
			if err != nil {
				return utils.NewErrorMessage("Error deploying the website", err)
			}
		}

		return nil
	}

	// Create the stack.
	stackName := auto.FullyQualifiedStackName(username, appConfig.DashCaseName, environment)
	stack, err := auto.UpsertStackInlineSource(context, stackName, appConfig.DashCaseName, deployFunc)
	if err != nil {
		return auto.Stack{}, err
	}

	return stack, nil
}
