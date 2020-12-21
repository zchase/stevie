package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zchase/stevie/pkg/auto_pulumi"
	"github.com/zchase/stevie/pkg/utils"
)

// whoAmI outputs the user name of the currently logged in Pulumi
// user.
func whoAmI(cmd *cobra.Command, args []string) {
	// Get the Pulumi username.
	username, err := auto_pulumi.GetCurrentPulumiUser()
	utils.CheckForNilAndHandleError(err, "Error getting currently logged in user")

	utils.Printf("The currently logged in user is: %s", utils.TextColor(username, color.FgGreen))
}

var whoAmICommand = &cobra.Command{
	Use:   "whoami",
	Short: "Output the currently logged in user.",
	Long:  "Output the currently logged in Pulumi user.",
	Run:   whoAmI,
}

func init() {
	RootCmd.AddCommand(whoAmICommand)
}
