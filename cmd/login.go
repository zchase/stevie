package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zchase/stevie/pkg/auto_pulumi"
	"github.com/zchase/stevie/pkg/utils"
)

func runPulumiLogin(cmd *cobra.Command, args []string) {
	err := auto_pulumi.PromptForPulumiAccessToken()
	if err != nil {
		utils.HandleError("Error logging into the Pulumi CLI", err)
	}

	utils.Print("Successfully logged into the Pulumi CLI.")
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the Pulumi CLI",
	Run:   runPulumiLogin,
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
