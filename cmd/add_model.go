package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zchase/stevie/pkg/application"
	"github.com/zchase/stevie/pkg/utils"
)

var modelLanguage string
var hashKeyName string
var hashKeyType string
var rangeKeyName string
var rangeKeyType string
var modelName string

// promptForControllerLanguage prompts the user to select a language for
// their new controller.
func promptForModelLanguage(message string) string {
	language, err := utils.PromptSelection(message, []string{"typescript"})
	utils.CheckForNilAndHandleError(err, "Error selecting controller language")

	return language
}

func createNewModel(cmd *cobra.Command, args []string) {
	var err error

	// Create a spinner while configuring the model.
	modelSpinner := utils.CreateNewTerminalSpinner(
		"Creating your new model",
		"Successfully creating your new model",
		"Failed to create your new model",
	)

	// Check if the name is defined and if it is not prompt
	// the user for the controller name.
	if modelName == "" {
		modelName, err = utils.PromptRequiredString("What is the name of the model?")
		if err != nil {
			modelSpinner.Fail()
			utils.HandleError("Error prompting for model name", err)
		}
	}

	switch modelLanguage {
	case application.TypeScriptLanguage:
		break
	case "":
		modelLanguage = promptForModelLanguage("Please pick the language to write your model with")
		break
	default:
		modelLanguage = promptForModelLanguage("Unknown model langauge provided. Please selected a valid language")
		break
	}

	if hashKeyName == "" {
		hashKeyName, err = utils.PromptRequiredString("What is the name of your hash key?")
		if err != nil {
			modelSpinner.Fail()
			utils.HandleError("Error prompting for hash key name", err)
		}
	}

	if hashKeyType == "" {
		hashKeyType, err = utils.PromptSelection("What is the type of your hash key?", []string{"string", "number"})
	}

	modelArgs := application.ModelArgs{
		Name:         utils.DashCaseToSentenceCase(modelName),
		HashKeyName:  hashKeyName,
		HashKeyType:  hashKeyType,
		RangeKeyName: rangeKeyName,
		RangeKeyType: rangeKeyType,
	}
	err = application.CreateNewModel(modelName, modelLanguage, modelArgs)
	utils.CheckForNilAndHandleError(err, "Error creating new model")
	modelSpinner.Stop()
}

var addModelCmd = &cobra.Command{
	Use:   "add-model",
	Short: "Create a new model.",
	Long:  "Create a new model.",
	Run:   createNewModel,
}

func init() {
	RootCmd.AddCommand(addModelCmd)

	addModelCmd.Flags().StringVar(&modelName, "name", "", "The name of the model.")
	addModelCmd.Flags().StringVar(&modelLanguage, "language", "", "The language to define your model in.")
	addModelCmd.Flags().StringVar(&hashKeyName, "hashKeyName", "", "The name of the hash key for your table.")
	addModelCmd.Flags().StringVar(&hashKeyType, "hashKeyType", "", "The type for your hash key.")
	addModelCmd.Flags().StringVar(&rangeKeyName, "rangeKeyName", "", "The name of the range key for your table.")
	addModelCmd.Flags().StringVar(&rangeKeyType, "rangeKeyType", "", "The type of your range key.")
}
