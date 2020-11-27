package cmd

import (
	"context"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zchase/stevie/pkg/application"
	"github.com/zchase/stevie/pkg/auto_pulumi"
	"github.com/zchase/stevie/pkg/utils"
)

var tmpDirName = "tmp"

func previewAPIChanges(cmd *cobra.Command, args []string) {
	utils.Print(utils.TextColor("Previewing Changes to API infrastructure\n", color.Bold))
	utils.Print("Preview steps:")
	ctx := context.Background()

	// Check the environment is set and error out if it is not.
	//
	// TODO: add pick list to choose env if it is not set.
	if Environment == "" {
		utils.HandleError("Please provide an environment via the enviroment flag", nil)
	}

	// Create the preview action.
	previewAction := auto_pulumi.PulumiAction{
		Environment:      Environment,
		CreateDeployment: CreateAPIDeployment,
	}

	// Set up the preview action
	err := previewAction.SetUp(ctx, application.ApplicationConfigPath)
	if err != nil {
		utils.HandleError("Error setting up preview: ", err)
	}

	// Run the preview
	err = previewAction.Preview(ctx)
	if err != nil {
		utils.HandleError("Error running preview: ", err)
	}

	utils.ClearLine()
	utils.Print("Preview Completed!")
}

var previewCmd = &cobra.Command{
	Use:   "preview [environment]",
	Short: "Preview the changes for the API.",
	Long:  `Preview the changes for the API.`,
	Run:   previewAPIChanges,
}

func init() {
	RootCmd.AddCommand(previewCmd)
}
