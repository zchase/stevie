package application

import (
	"path"
	"strings"

	"github.com/zchase/stevie/pkg/utils"
)

type APIModel struct {
	Name string
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

		models = append(models, APIModel{
			Name: strings.Split(modelFile, ".")[0],
		})
	}

	return models, nil
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
	actionsFilePath := path.Join(actionsDirPath, "index.ts")
	err = utils.WriteOutTemplateToFile(modelsActionsStoresFileTemplatePath, actionsFilePath, nil)
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
	err = utils.WriteOutTemplateToFile(constantsFileTemplatePath, constantsFilePath, nil)
	if err != nil {
		return utils.NewErrorMessage("Error writing constants file", err)
	}

	// Install the module.
	err = utils.RunCommandWithOutput("yarn", []string{"---cwd", "ui", "add", "file:./ui/stevie-ui-utils"})
	if err != nil {
		return utils.NewErrorMessage("Error installing utils modules", err)
	}

	return nil
}
