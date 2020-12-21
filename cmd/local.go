package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zchase/stevie/pkg/application"
	"github.com/zchase/stevie/pkg/utils"
)

var localServerPort int

// runLocal runs the API routes locally.
func runLocal(cmd *cobra.Command, args []string) {
	utils.Print("Running API Routes locally.")

	// Read the config.
	config, err := ReadBaseConfig(application.ApplicationConfigPath)
	utils.CheckForNilAndHandleError(err, "Error reading base config")

	// Run the routes.
	err = application.RunTypeScriptRoutesLocally(config.Routes, localServerPort)
	utils.CheckForNilAndHandleError(err, "Error running local server")

	utils.Print("Finished running local server")
}

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Run the API routes locally",
	Run:   runLocal,
}

func init() {
	RootCmd.AddCommand(localCmd)

	localCmd.Flags().IntVarP(&localServerPort, "port", "p", 3000, "The port you'd like to run your server on. Defaults to 3000.")
}
