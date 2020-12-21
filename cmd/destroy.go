package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zchase/stevie/pkg/application"
	"github.com/zchase/stevie/pkg/auto_pulumi"
	"github.com/zchase/stevie/pkg/utils"
)

func destroyAPI(cmd *cobra.Command, args []string) {
	fmt.Println("Destroying API infrastructure:")
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
	utils.CheckForNilAndHandleError(err, "Error setting up API resource destroy operation")

	// Run the preview
	err = updateAction.Destroy(ctx)
	utils.CheckForNilAndHandleError(err, "Error destroy API resources")

	fmt.Println("Destroy Completed!")
}

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the API.",
	Long:  `Destroy the API.`,
	Run:   destroyAPI,
}

func init() {
	RootCmd.AddCommand(destroyCmd)
}
