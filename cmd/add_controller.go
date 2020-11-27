package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zchase/stevie/pkg/application"
	"github.com/zchase/stevie/pkg/utils"
)

var validControllerMethods = map[string]bool{
	"GET":    true,
	"POST":   true,
	"PUT":    true,
	"DELETE": true,
}

// Flags
var name string
var methods []string

var addControllerCmd = &cobra.Command{
	Use:   "add-controller",
	Short: "Create a new file.",
	Long:  "Create a new file.",
	Run:   addNewFile,
}

func addNewFile(cmd *cobra.Command, args []string) {
	var err error

	// To start let's make sure we have all the config values we
	// need to create the route.
	configSpinner := utils.TerminalSpinner{
		SpinnerText:   "Configuring controller",
		CompletedText: "✅ Successfully configured controller.",
		FailureText:   "❌ Failed configuring controller.",
	}
	configSpinner.Create()

	// Read the base config
	baseConfig, err := ReadBaseConfig(application.ApplicationConfigPath)
	if err != nil {
		utils.HandleError("Error reading application config: %v", err)
	}

	// Set the backend langauge for creating the controller.
	backendLanguage := baseConfig.BackendLanguage

	// Check if the name is defined and if it is not prompt
	// the user for the controller name.
	if name == "" {
		name, err = utils.PromptRequiredString("What is the name of the controller?")
		if err != nil {
			configSpinner.Fail()
			utils.HandleError("Error prompting for controller name: ", err)
		}
	}

	// Check the provided methods are valid.
	var controllerMethods []string
	for _, inputMethod := range methods {
		method := strings.ToUpper(inputMethod)
		if !validControllerMethods[method] {
			configSpinner.Fail()
			errMsg := fmt.Sprintf("Unkown method provided: %s", inputMethod)
			utils.HandleError(errMsg, nil)
		}

		controllerMethods = append(controllerMethods, method)
	}
	configSpinner.Stop()

	// Create the controller file.
	controllerSpinner := utils.TerminalSpinner{
		SpinnerText:   "Creating controller files.",
		CompletedText: "✅ Successfully created controller.",
		FailureText:   "❌ Failed to create controller.",
	}
	controllerSpinner.Create()

	// Create the controller file.
	switch backendLanguage {
	case "typescript":
		controllerFile, err := application.CreateNewTypeScriptController(name, controllerMethods)
		if err != nil {
			utils.HandleError("Error creating TypeScript controller: ", err)
		}

		err = AddAPIRouteToConfig(application.ApplicationConfigPath, name, fmt.Sprintf("/%s", name), controllerFile, controllerMethods)
		if err != nil {
			utils.HandleError("Error adding route to config: ", err)
		}
	default:
		errMessage := fmt.Sprintf("Unsupported langauge in config: %x", backendLanguage)
		utils.HandleError(errMessage, nil)
	}

	// Output the controller has been created.
	controllerSpinner.Stop()
	fmt.Println("Controller command finished.")
}

func init() {
	RootCmd.AddCommand(addControllerCmd)

	addControllerCmd.Flags().StringVar(&name, "name", "", "The name for the controller in camelCase.")
	addControllerCmd.Flags().StringSliceVar(&methods, "methods", []string{"GET"}, "The methods for your route. GET, POST, PUT, & DELETE.")
}
