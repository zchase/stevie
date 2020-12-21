package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zchase/stevie/pkg/application"
	"github.com/zchase/stevie/pkg/auto_pulumi"
	"github.com/zchase/stevie/pkg/utils"
)

var (
	// Default environments.
	DefaultEnvironments = []string{"production", "development", "testing"}

	// Supported languages
	SupportedLanguages = []string{
		application.TypeScriptControllerLanguage,
		application.GoControllerLanguage,
		application.DotNetControllerLanguage,
	}
)

// Shared flags
var Environment string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new project.",
	Long:  `Create a new Stevie project in an empty directory.`,
	Run:   createNewProject,
}

func createNewProject(cmd *cobra.Command, args []string) {
	fmt.Println("Creating a new Stevie project")
	ctx := context.Background()

	// Check if the working directory is empty.
	isEmptyDir, err := utils.IsCurrentDirectoryEmpty()
	utils.CheckForNilAndHandleError(err, "Error checking contents of current working directory")

	// Throw an error if the current working directory is not empty.
	if isEmptyDir == false {
		utils.HandleError("Current working directory is not empty", nil)
	}

	// Check the user is logged in.
	username, err := auto_pulumi.GetCurrentPulumiUser()
	utils.CheckForNilAndHandleError(err, "Error checking for authenticated user")

	// Set the config path
	configPath := application.ApplicationConfigPath

	// Prompt the user for the project name and description.
	appConfig, err := CreateApplicationConfig(configPath, "", "", DefaultEnvironments)
	utils.CheckForNilAndHandleError(err, "Error setting up application config")

	// Create the spinner for the new project.
	setupSpinner := utils.CreateNewTerminalSpinner(
		"Setting up your new project",
		"Successfully set up project.",
		"Failed to set up project",
	)

	// Create the Pulumi Project.
	for _, env := range DefaultEnvironments {
		projectName, err := auto_pulumi.CreatePulumiProject(ctx, username, appConfig.DashCaseName, env, appConfig.Description)
		if err != nil {
			setupSpinner.Fail()
			utils.HandleError("Error creating Pulumi project", err)
		}

		//utils.ClearLine()
		utils.Printf("Created project: %s", projectName)
	}
	setupSpinner.Stop()

	// Create the initial project structure. First we will create
	// the application directories.
	createProjectSpinner := utils.CreateNewTerminalSpinner(
		"Creating your new project",
		"Successfully created your project.",
		"Failed to create your project.",
	)

	// Create the project structure based on the backend-language chosen.
	err = application.CreateProjectStructure(appConfig.DashCaseName, appConfig.Description)
	utils.CheckForNilAndHandleError(err, "Error creating the project structure")

	createProjectSpinner.Stop()
	utils.ClearLine()
	fmt.Println("Project has been successfully created!")
}

func init() {
	RootCmd.AddCommand(initCmd)

	RootCmd.PersistentFlags().StringVarP(&Environment, "environment", "e", "", "The environment you are deploying to.")
}
