package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zchase/stevie/pkg/application"
	"github.com/zchase/stevie/pkg/auto_pulumi"
	"github.com/zchase/stevie/pkg/utils"
)

var (
	// Default environments.
	DefaultEnvironments = []string{"production", "development", "testing"}
)

// Shared flags
var Environment string
var BackendLanguage string

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
	if err != nil {
		utils.HandleError("Error checking contents of current working directory: ", err)
	}

	// Throw an error if the current working directory is not empty.
	if isEmptyDir == false {
		utils.HandleError("Current working directory is not empty", nil)
	}

	// Check the user is logged in.
	username, err := auto_pulumi.GetCurrentPulumiUser()
	if err != nil {
		utils.HandleError("Error checking for authenticated user: %v", err)
	}

	// Set the config path
	configPath := application.ApplicationConfigPath

	// Check that the backend langauge is supported.
	backendLanguage := strings.ToLower(BackendLanguage)
	switch backendLanguage {
	case "typescript":
	default:
		errMessage := fmt.Sprintf("%s is not a supported language.", BackendLanguage)
		utils.HandleError(errMessage, nil)
	}

	// Prompt the user for the project name and description.
	appConfig, err := CreateApplicationConfig(configPath, "", "", backendLanguage, DefaultEnvironments)
	if err != nil {
		utils.HandleError("Error setting up application config: ", err)
	}

	// Create the spinner for the new project.
	setupSpinner := utils.TerminalSpinner{
		SpinnerText:   "Setting up your new project",
		CompletedText: "✅ Successfully set up project.",
		FailureText:   "❌ Failed to set up project",
	}
	setupSpinner.Create()

	// Create the Pulumi Project.
	fmt.Println("Creating Pulumi project and stacks...")
	for _, env := range DefaultEnvironments {
		err = auto_pulumi.CreatePulumiProject(ctx, username, appConfig.DashCaseName, env, appConfig.Description)
		if err != nil {
			setupSpinner.Fail()
			utils.HandleError("Error creating Pulumi project: ", err)
		}
	}
	setupSpinner.Stop()

	// Create the intitial project structure. First we will create
	// the application directories.
	createProjectSpinner := utils.TerminalSpinner{
		SpinnerText:   "Creating your new project",
		CompletedText: "✅ Successfully created your project.",
		FailureText:   "❌ Failed to create your project",
	}
	createProjectSpinner.Create()

	// Create the project structure based on the backend-language chosen.
	switch backendLanguage {
	case "typescript":
		err = application.CreateTypeScriptProject(appConfig.DashCaseName, appConfig.Description)
		if err != nil {
			utils.HandleError("Error creating TypeScript project: ", err)
		}
	}

	createProjectSpinner.Stop()
	fmt.Println("Project has been successfully created!")
}

func init() {
	RootCmd.AddCommand(initCmd)

	RootCmd.PersistentFlags().StringVarP(&Environment, "environment", "e", "", "The environemnt you are deploying to.")
	RootCmd.Flags().StringVar(&BackendLanguage, "backend-language", "typescript", "The langauge for writitng your serverless functions. Defaults to TypeScript.")
}
