package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "stevie",
	Short: "Stevie helps you build serverless applications",
	Long: `Stevie helps build serverless applications by letting you
focus on your business logic instead of the cloud services needed to run it.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
