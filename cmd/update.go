package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zchase/stevie/pkg/application"
	"github.com/zchase/stevie/pkg/auto_pulumi"
	"github.com/zchase/stevie/pkg/utils"
)

func updateAPIChanges(cmd *cobra.Command, args []string) {
	fmt.Println("Updating Changes to API infrastructure:")
	ctx := context.Background()

	// Check the environment is set and error out if it is not.
	//
	// TODO: add pick list to choose env if it is not set.
	if Environment == "" {
		utils.HandleError("Please provide an environment via the environment flag", nil)
	}

	// Create the preview action.
	updateAction := auto_pulumi.PulumiAction{
		Environment:      Environment,
		CreateDeployment: CreateAPIDeployment,
	}

	// Set up the preview action
	err := updateAction.SetUp(ctx, application.ApplicationConfigPath)
	utils.CheckForNilAndHandleError(err, "Error setting up update")

	// Run the preview
	err = updateAction.Update(ctx)
	utils.CheckForNilAndHandleError(err, "Error running update")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the API.",
	Long:  `Update the API.`,
	Run:   updateAPIChanges,
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
