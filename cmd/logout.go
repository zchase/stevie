package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zchase/stevie/pkg/auto_pulumi"
	"github.com/zchase/stevie/pkg/utils"
)

func runPulumiLogout(cmd *cobra.Command, args []string) {
	err := auto_pulumi.LogoutCurrentUser(args)
	if err != nil {
		utils.HandleError("Error logging out of the Pulumi CLI", err)
	}

	utils.Print("Successfully logged out of the Pulumi CLI.")
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout of the Pulumi CLI",
	Run:   runPulumiLogout,
}

func init() {
	RootCmd.AddCommand(logoutCmd)
}
