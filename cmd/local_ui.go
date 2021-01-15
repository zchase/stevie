package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
	"github.com/spf13/cobra"
	"github.com/zchase/stevie/pkg/application"
	"github.com/zchase/stevie/pkg/auto_pulumi"
	"github.com/zchase/stevie/pkg/utils"
)

var backendEnvironment string

func RunLocalUI(cmd *cobra.Command, args []string) {
	utils.Print("Running React App Locally")

	if backendEnvironment == "" {
		utils.HandleError("Please provide an environment via the --environment flag", nil)
	}

	// Check that user is authenticated.
	username, err := auto_pulumi.GetCurrentPulumiUser()
	utils.CheckForNilAndHandleError(err, "Error getting current Pulumi user")

	// Read in the base config file.
	appConfig, err := ReadBaseConfig(application.ApplicationConfigPath)
	utils.CheckForNilAndHandleError(err, "Error reading application config")

	ctx := context.Background()
	stackName := auto.FullyQualifiedStackName(username, appConfig.DashCaseName, backendEnvironment)
	stack, err := auto.UpsertStackInlineSource(ctx, stackName, appConfig.DashCaseName, nil)
	utils.CheckForNilAndHandleError(err, "Error selecting stack")

	endpoints, err := stack.Outputs(ctx)
	utils.CheckForNilAndHandleError(err, "Error getting stack outputs")

	for name, value := range endpoints {
		envVarName := fmt.Sprintf("REACT_APP_%s_ENDPOINT", strings.ToUpper(name))
		os.Setenv(envVarName, value.Value.(string))
	}

	// Run the app.
	err = utils.RunCommandWithOutput("yarn", []string{"--cwd", "ui", "start"})
}

var localUICmd = &cobra.Command{
	Use:   "local-ui",
	Short: "Run the react app locally against an environment",
	Run:   RunLocalUI,
}

func init() {
	RootCmd.AddCommand(localUICmd)
	localUICmd.Flags().StringVarP(&backendEnvironment, "environment", "e", "", "The backend environment to run your app against.")
}
