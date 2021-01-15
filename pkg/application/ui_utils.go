package application

import (
	"fmt"
	"path"
	"strings"

	"github.com/zchase/stevie/pkg/utils"
)

type APIModel struct {
	Name   string
	Schema utils.JSONSchema
}

type ActionModelStoreImportItems struct {
	ImportName string
}
type ActionModelStoreIndexArgs struct {
	Items []ActionModelStoreImportItems
}

// readAndCopyAPIModels reads the API models and copies them into the utils package.
func readAndCopyAPIModels(utilsModelFilePath string) ([]APIModel, error) {
	var models []APIModel
	modelsDirPath := path.Join(ApplicationFolder, ModelsFolder)
	modelContents, err := utils.ReadDirectoryContents(modelsDirPath)
	if err != nil {
		return nil, err
	}

	for _, modelFile := range modelContents {
		oldFilePath := path.Join(modelsDirPath, modelFile)
		newFilePath := path.Join(utilsModelFilePath, modelFile)
		err = utils.CopyFile(oldFilePath, newFilePath)
		if err != nil {
			return nil, err
		}

		schema, err := utils.GenerateModelSchemaFromFile(newFilePath)
		if err != nil {
			return nil, err
		}

		models = append(models, APIModel{
			Name:   strings.Split(modelFile, ".")[0],
			Schema: schema,
		})
	}

	return models, nil
}

type ActionFileActionFunctionArgs struct {
	Name         string
	ArgsType     string
	Method       string
	ConstantName string
}

type ActionFileArgs struct {
	UrlEnvVarName   string
	ModelImportName string
	Name            string
	Actions         []ActionFileActionFunctionArgs
}

type ConstantArgs struct {
	Name  string
	Value string
}

type ConstantFileArgs struct {
	Constants []ConstantArgs
}

// createConstant creates a constant.
func createConstant(method, name string) string {
	return fmt.Sprintf("%s_%s", strings.ToUpper(method), strings.ToUpper(name))
}

// writeActionFiles write the action files.
func writeActionFiles(templatePath, utilsPath string, models []APIModel) ([]ConstantArgs, error) {
	var result []ConstantArgs

	for _, model := range models {
		// Add the different methods to support starting with get.
		var actionsFunctions []ActionFileActionFunctionArgs
		methods := [4]string{"get", "put", "post", "delete"}
		for _, method := range methods {
			methodConstant := createConstant(method, model.Name)
			actionsFunctions = append(actionsFunctions, ActionFileActionFunctionArgs{
				Name:         utils.DashCaseToSentenceCase(fmt.Sprintf("%s-%s", method, model.Name)),
				ArgsType:     utils.DashCaseToSentenceCase(model.Name),
				Method:       method,
				ConstantName: methodConstant,
			})
			result = append(result, ConstantArgs{
				Name:  methodConstant,
				Value: methodConstant,
			})
		}

		actionFileArgs := ActionFileArgs{
			Name:            model.Name,
			UrlEnvVarName:   strings.ToUpper(model.Name),
			ModelImportName: utils.DashCaseToSentenceCase(model.Name),
			Actions:         actionsFunctions,
		}

		fileName := fmt.Sprintf("%s.ts", model.Name)
		templatePath := path.Join(templatePath, "action.tmpl")
		filePath := path.Join(utilsPath, fileName)
		err := utils.WriteOutTemplateToFile(templatePath, filePath, actionFileArgs)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// CreateUIUtilsPackage generates the UI utils package.
func CreateUIUtilsPackage(uiDir string) error {
	utilsDirPath := path.Join(uiDir, "stevie-ui-utils")
	uiTemplatesPath := path.Join(FileTemplatePath, "ui")

	// Create the main directory.
	err := utils.CreateNewDirectory(utilsDirPath)
	if err != nil {
		return err
	}

	// Write the root files.
	rootFileTemplatePath := path.Join(uiTemplatesPath, "index.tmpl")
	rootFilePath := path.Join(utilsDirPath, "index.ts")
	err = utils.WriteOutTemplateToFile(rootFileTemplatePath, rootFilePath, nil)
	if err != nil {
		return utils.NewErrorMessage("Error writing root index file", err)
	}

	packageJsonFileTemplatePath := path.Join(uiTemplatesPath, "package_json.tmpl")
	packageJsonFilePath := path.Join(utilsDirPath, "package.json")
	err = utils.WriteOutTemplateToFile(packageJsonFileTemplatePath, packageJsonFilePath, nil)
	if err != nil {
		return utils.NewErrorMessage("Error writing root package.json file", err)
	}

	// Create the models directory.
	modelsDirPath := path.Join(utilsDirPath, "models")
	err = utils.CreateNewDirectory(modelsDirPath)
	if err != nil {
		return utils.NewErrorMessage("Error creating models directory", err)
	}

	models, err := readAndCopyAPIModels(modelsDirPath)
	if err != nil {
		return utils.NewErrorMessage("Error copying API models", err)
	}

	var modelsIndexFileArgs []ActionModelStoreImportItems
	for _, model := range models {
		modelsIndexFileArgs = append(modelsIndexFileArgs, ActionModelStoreImportItems{
			ImportName: model.Name,
		})
	}

	modelsActionsStoresFileTemplatePath := path.Join(uiTemplatesPath, "action_store_index_file.tmpl")
	modelsFilePath := path.Join(modelsDirPath, "index.ts")
	err = utils.WriteOutTemplateToFile(modelsActionsStoresFileTemplatePath, modelsFilePath, ActionModelStoreIndexArgs{
		Items: modelsIndexFileArgs,
	})
	if err != nil {
		return utils.NewErrorMessage("Error writing model index file", err)
	}

	// Create the actions directory.
	actionsDirPath := path.Join(utilsDirPath, "actions")
	err = utils.CreateNewDirectory(actionsDirPath)
	if err != nil {
		return utils.NewErrorMessage("Error creating actions directory", err)
	}

	// Write the actions files.
	constants, err := writeActionFiles(uiTemplatesPath, actionsDirPath, models)
	if err != nil {
		return err
	}

	createActionFileTemplatePath := path.Join(uiTemplatesPath, "create_action.tmpl")
	createActionFileName := path.Join(actionsDirPath, "create_action.ts")
	err = utils.WriteOutTemplateToFile(createActionFileTemplatePath, createActionFileName, nil)
	if err != nil {
		return err
	}

	restClientFileTemplatePath := path.Join(uiTemplatesPath, "rest_client.tmpl")
	restClientFileName := path.Join(actionsDirPath, "rest_client.ts")
	err = utils.WriteOutTemplateToFile(restClientFileTemplatePath, restClientFileName, nil)
	if err != nil {
		return err
	}

	actionsFilePath := path.Join(actionsDirPath, "index.ts")
	err = utils.WriteOutTemplateToFile(modelsActionsStoresFileTemplatePath, actionsFilePath, ActionModelStoreIndexArgs{
		Items: modelsIndexFileArgs,
	})
	if err != nil {
		return utils.NewErrorMessage("Error writing actions index file", err)
	}

	// Create the stores directory.
	storesDirPath := path.Join(utilsDirPath, "stores")
	err = utils.CreateNewDirectory(storesDirPath)
	if err != nil {
		return utils.NewErrorMessage("Error creating stores directory", err)
	}

	// Write the stores file.
	storesFilePath := path.Join(storesDirPath, "index.ts")
	err = utils.WriteOutTemplateToFile(modelsActionsStoresFileTemplatePath, storesFilePath, nil)
	if err != nil {
		return utils.NewErrorMessage("Error writing stores index file", err)
	}

	// Create the constants directory and write the constants file.
	constantsDirPath := path.Join(utilsDirPath, "constants")
	err = utils.CreateNewDirectory(constantsDirPath)
	if err != nil {
		return utils.NewErrorMessage("Error constants directory", err)
	}

	constantsFileTemplatePath := path.Join(uiTemplatesPath, "constants.tmpl")
	constantsFilePath := path.Join(constantsDirPath, "index.ts")
	err = utils.WriteOutTemplateToFile(constantsFileTemplatePath, constantsFilePath, ConstantFileArgs{
		Constants: constants,
	})
	if err != nil {
		return utils.NewErrorMessage("Error writing constants file", err)
	}

	// Install the module.
	err = utils.RunCommandWithOutput("yarn", []string{"--cwd", "ui", "add", "file:./stevie-ui-utils"})
	if err != nil {
		return utils.NewErrorMessage("Error installing utils modules", err)
	}

	return nil
}
