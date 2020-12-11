package application

import (
	"fmt"
	"path"

	"github.com/zchase/stevie/pkg/utils"
)

var (
	// Config
	ApplicationConfigPath = "config"

	// File template paths
	FileTemplatePath = "/pkg/application/file_templates"
	LocalPackagePath = "/lib"

	// TypeScript
	TypeScriptFileTemplatesDirectoryName = "typescript"
	TypeScriptTSConfigFileName           = "tsconfig.json"
	TypeScriptTSConfigTemplateName       = "tsconfig_json.tmpl"
	TypeScriptPackageJSONFileName        = "package.json"
	TypeScriptPackageJSONTemplateName    = "package_json.tmpl"
	TypeScriptTestDirectoryName          = "test"
	TypeScriptAppDirectoryName           = "app"
	TypeScriptModelsDirectoryName        = "models"
	TypeScriptControllersDirectoryName   = "controllers"
	TypesScriptUtilsDirectory            = "stevie-utils"
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

	// Create the controllers directory in the app directory.
	controllersDirectoryPath := path.Join(TypeScriptAppDirectoryName, TypeScriptControllersDirectoryName)
	err = utils.CreateNewDirectory(controllersDirectoryPath)
	if err != nil {
		return err
	}

	// Copy in the utils package.
	utilsPackagePath := path.Join(LocalPackagePath, TypeScriptFileTemplatesDirectoryName)
	err = utils.CreateNewDirectory(TypesScriptUtilsDirectory)
	if err != nil {
		return fmt.Errorf("Error creating utils directory: %v", err)
	}

	err = utils.CopyPackagedDirectory(utilsPackagePath, TypesScriptUtilsDirectory, []string{"node_modules"})
	if err != nil {
		return fmt.Errorf("Error copy in TypeScript utilities: %v", err)
	}

	// Install the application dependencies.
	err = utils.RunCommandWithOutput("yarn", []string{})
	if err != nil {
		return fmt.Errorf("Error install app dependencies: %v", err)
	}

	return nil
}
