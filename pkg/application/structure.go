package application

import (
	"fmt"
	"path"

	"github.com/zchase/stevie/pkg/utils"
)

var (
	// Config
	ApplicationConfigPath = "config"

	// File template path
	FileTemplatePath = "/pkg/application/file_templates"

	// TypeScript
	TypeScriptFileTemplatesDirectoryName = "typescript"
	TypeScriptTSConfigFileName           = "tsconfig.json"
	TypeScriptTSConfigTemplateName       = "tsconfig_json.tmpl"
	TypeScriptPackageJSONFileName        = "package.json"
	TypeScriptPackageJSONTemplateName    = "package_json.tmpl"
	TypeScriptControllerBuilderFileName  = "controller_builder"
	TypeScriptServerArgsFileName         = "server_args"
	TypeScriptTestDirectoryName          = "test"
	TypeScriptAppDirectoryName           = "app"
	TypeScriptModelsDirectoryName        = "models"
	TypeScriptControllersDirectoryName   = "controllers"
)

// WriteOutTemplateFile writes out a template file.
func writeOutTemplateFile(destinationPath, templateName, name string, args interface{}) error {
	templatePath := path.Join(FileTemplatePath, TypeScriptFileTemplatesDirectoryName, templateName)
	destination := path.Join(destinationPath, name)

	err := utils.WriteOutTemplateToFile(templatePath, destination, args)
	return err
}

// PackageJsonArgs define the args need to generate a package.json file.
type PackageJsonArgs struct {
	Name        string
	Description string
}

// CreateTypeScriptProjectStructure creates the base structure for a TypeScript project.
func CreateTypeScriptProject(name, description string) error {
	// Create the package.json file.
	packageJSONArgs := PackageJsonArgs{
		Name:        name,
		Description: description,
	}
	err := writeOutTemplateFile("", TypeScriptPackageJSONTemplateName, TypeScriptPackageJSONFileName, packageJSONArgs)
	if err != nil {
		return err
	}

	// Create tsconfig.json file.
	err = writeOutTemplateFile("", TypeScriptTSConfigTemplateName, TypeScriptTSConfigFileName, nil)
	if err != nil {
		return err
	}

	// Create the test directory.
	err = utils.CreateNewDirectory(TypeScriptTestDirectoryName)
	if err != nil {
		return err
	}

	// Create the app directory.
	err = utils.CreateNewDirectory(TypeScriptAppDirectoryName)
	if err != nil {
		return err
	}

	// Create the models directory in the app directory.
	modelsDirectoryPath := path.Join(TypeScriptAppDirectoryName, TypeScriptModelsDirectoryName)
	err = utils.CreateNewDirectory(modelsDirectoryPath)
	if err != nil {
		return err
	}

	// Create the controllers directory in the app directory.
	controllersDirectoryPath := path.Join(TypeScriptAppDirectoryName, TypeScriptControllersDirectoryName)
	err = utils.CreateNewDirectory(controllersDirectoryPath)
	if err != nil {
		return err
	}

	// Create the controller builder file in the controller directory.
	controllerBuilderOutputFileName := fmt.Sprintf("%s.ts", TypeScriptControllerBuilderFileName)
	controllerBuilderTemplateName := fmt.Sprintf("%s.tmpl", TypeScriptControllerBuilderFileName)
	err = writeOutTemplateFile(controllersDirectoryPath, controllerBuilderTemplateName, controllerBuilderOutputFileName, nil)
	if err != nil {
		return err
	}

	// Create the server args file in the controller directory.
	serverArgsOutputFileName := fmt.Sprintf("%s.ts", TypeScriptServerArgsFileName)
	serverArgsTemplateName := fmt.Sprintf("%s.tmpl", TypeScriptServerArgsFileName)
	err = writeOutTemplateFile(controllersDirectoryPath, serverArgsTemplateName, serverArgsOutputFileName, nil)
	if err != nil {
		return err
	}

	// Install the application dependencies.
	err = utils.RunCommandWithOutput("yarn", []string{})
	if err != nil {
		utils.HandleError("Erroring install app dependencies: ", err)
	}

	return nil
}
