package cmd

import (
	"fmt"
	"path"
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
var fromModelName string

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

func addNewControllerFromModel(modelName, controllerLanguage string) error {
	var err error

	// Create a spinner.
	controllerFromModelSpinner := utils.CreateNewTerminalSpinner(
		fmt.Sprintf("Creating CRUD APIs from %s model", modelName),
		fmt.Sprintf("Successfully created CRUD APIs from %s model", modelName),
		fmt.Sprintf("Failed to create CRUD APIs from %s model", modelName),
	)

	// Verify the language is valid.
	switch controllerLanguage {
	case application.TypeScriptControllerLanguage:
		break
	case "":
		controllerLanguage, err = utils.PromptSelection("Unknown controller langauge provided. Please selected a valid language", []string{"typescript"})
		utils.CheckForNilAndHandleError(err, "Error selecting controller language")
		break
	default:
		controllerLanguage, err = utils.PromptSelection("No controller langauge provided. Please selected a valid language", []string{"typescript"})
		utils.CheckForNilAndHandleError(err, "Error selecting controller language")
		break
	}

	// Lookup the model.
	modelDirPath := path.Join(application.ApplicationFolder, application.ModelsFolder)
	modelDirContents, err := utils.ReadDirectoryContents(modelDirPath)
	utils.CheckForNilAndHandleError(err, "Error reading contents of model directory")

	var modelFileName string
	for _, content := range modelDirContents {
		fileName := strings.Split(content, ".")[0]
		if fileName == modelName {
			modelFileName = content
		}
	}

	if modelFileName == "" {
		utils.HandleError("No matching model found for the provided model name", nil)
	}

	modelFilePath := path.Join(modelDirPath, modelFileName)
	fmt.Println(modelFilePath)
	schema, err := utils.GenerateModelSchemaFromFile(modelFilePath)
	utils.CheckForNilAndHandleError(err, "Error generating schema from model file")

	// Create the controllers from the model definition.
	controllerPath, err := application.GenerateControllerFromModelSchema(modelName, controllerLanguage, schema)
	utils.CheckForNilAndHandleError(err, "Error creating controller files from model")

	err = AddAPIRouteToConfig(application.ApplicationConfigPath, modelName, fmt.Sprintf("/%s", modelName), controllerPath, corsEnabled)
	utils.CheckForNilAndHandleError(err, "Error writing new controller to config")

	controllerFromModelSpinner.Stop()
	return nil
}

func addNewController(cmd *cobra.Command, args []string) {
	var err error

	// If the controller is being generated from a model we need to handle
	// that case and return once we've generated the controllers.
	if fromModelName != "" {
		err = addNewControllerFromModel(fromModelName, controllerLanguage)
		utils.CheckForNilAndHandleError(err, "Error generating controllers from model definition")
		return
	}

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
	addControllerCmd.Flags().StringSliceVar(&methods, "methods", []string{"get"}, "The methods for your route. GET, POST, PUT, & DELETE.")
	addControllerCmd.Flags().StringVar(&controllerLanguage, "language", "", "The language you are using to write your controller.")
	addControllerCmd.Flags().BoolVar(&corsEnabled, "cors", false, "Enable CORS on your path.")
	addControllerCmd.Flags().StringVar(&fromModelName, "from-model-name", "", "Generate CRUD APIs from a model.")
}
