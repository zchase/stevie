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
var corsEnabled bool
var controllerLanguage string

var addControllerCmd = &cobra.Command{
	Use:   "add-controller",
	Short: "Create a new file.",
	Long:  "Create a new file.",
	Run:   addNewController,
}

// promptForControllerLanguage prompts the user to select a language for
// their new controller.
func promptForControllerLanguage(message string) string {
	language, err := utils.PromptSelection(message, SupportedLanguages)
	utils.CheckForNilAndHandleError(err, "Error selecting controller language")

	return language
}

func addNewController(cmd *cobra.Command, args []string) {
	var err error

	// To start let's make sure we have all the config values we
	// need to create the route.
	configSpinner := utils.CreateNewTerminalSpinner(
		"Configuring controller",
		"Successfully configured controller.",
		"Failed configuring controller.",
	)

	// Check the backend language is valid otherwise we need to prompt the user to
	// choose their language.
	switch controllerLanguage {
	case application.TypeScriptControllerLanguage, application.GoControllerLanguage, application.DotNetControllerLanguage:
		break
	case "":
		controllerLanguage = promptForControllerLanguage("Please pick the language to write your controller with")
		break
	default:
		controllerLanguage = promptForControllerLanguage("Unknown controller langauge provided. Please selected a valid language")
		break
	}

	// Check if the name is defined and if it is not prompt
	// the user for the controller name.
	if name == "" {
		name, err = utils.PromptRequiredString("What is the name of the controller?")
		if err != nil {
			configSpinner.Fail()
			utils.HandleError("Error prompting for controller name", err)
		}
	}

	// Check the provided methods are valid.
	var controllerMethods []string
	for _, inputMethod := range methods {
		method := strings.ToUpper(inputMethod)
		if !validControllerMethods[method] {
			configSpinner.Fail()
			errMsg := fmt.Sprintf("Unknown method provided: %s", inputMethod)
			utils.HandleError(errMsg, nil)
		}

		controllerMethods = append(controllerMethods, strings.ToLower(method))
	}

	// Create the controller file.
	configSpinner.Stop()
	controllerSpinner := utils.CreateNewTerminalSpinner(
		"Creating controller files.",
		"Successfully created controller.",
		"Failed to create controller.",
	)

	// Create the controller file(s).
	controllerPath, err := application.CreateNewController(name, methods, controllerLanguage)
	utils.CheckForNilAndHandleError(err, "Error creating controller files")

	// Add the controller to the config.
	err = AddAPIRouteToConfig(application.ApplicationConfigPath, name, fmt.Sprintf("/%s", name), controllerPath, corsEnabled)
	utils.CheckForNilAndHandleError(err, "Error writing new controller to config")

	// Output the controller has been created.
	controllerSpinner.Stop()
	fmt.Println("Controller command finished.")
}

func init() {
	RootCmd.AddCommand(addControllerCmd)

	addControllerCmd.Flags().StringVar(&name, "name", "", "The name for the controller in camelCase.")
	addControllerCmd.Flags().StringSliceVar(&methods, "methods", []string{"GET"}, "The methods for your route. GET, POST, PUT, & DELETE.")
	addControllerCmd.Flags().StringVar(&controllerLanguage, "language", "", "The language you are using to write your controller.")
	addControllerCmd.Flags().BoolVar(&corsEnabled, "cors", false, "Enable CORS on your path.")
}
