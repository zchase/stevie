package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zchase/stevie/pkg/application"
	"github.com/zchase/stevie/pkg/utils"
)

func AddUI(cmd *cobra.Command, args []string) {
	// Create a spinner to show progress.
	createReactAppSpinner := utils.CreateNewTerminalSpinner(
		"Creating new React project",
		"Successfully created a new React project",
		"Failed to create a React project",
	)

	// Check that npx is install.
	npxExecPath, err := utils.RunCommand("which", []string{"npx"})
	if err != nil {
		createReactAppSpinner.Fail()
		utils.HandleError("Error checking if npx is installed", err)
	}
	if npxExecPath == "" {
		createReactAppSpinner.Fail()
		utils.HandleError("Error finding npx. Please ensure you have version >= 5.2 of npm installed.", nil)
	}

	// Create the react app.
	err = utils.RunCommandWithOutput("npx", []string{"create-react-app", "ui", "--template", "redux-typescript"})
	if err != nil {
		createReactAppSpinner.Fail()
		utils.HandleError("Error creating React application", err)
	}

	// Add the UI utils.
	err = application.CreateUIUtilsPackage("ui")
	if err != nil {
		createReactAppSpinner.Fail()
		utils.HandleError("Error creating UI utils module", err)
	}

	createReactAppSpinner.Stop()
}

var addUICommand = cobra.Command{
	Use:   "add-ui",
	Short: "Adds a UI using create-react-app",
	Long:  "Adds a UI using create-react-app",
	Run:   AddUI,
}

func init() {
	RootCmd.AddCommand(&addUICommand)
}
