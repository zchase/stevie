package application

import (
	"fmt"
	"path"

	"github.com/zchase/stevie/pkg/utils"
)

var (
	// Config
	ApplicationConfigPath = "config"

	// Directories
	ApplicationFolder = "app"
	ControllersFolder = "controllers"

	// File template paths
	FileTemplatePath = "/pkg/application/file_templates"
	LocalPackagePath = "/lib"

	// TypeScript
	TypeScriptFileTemplatesDirectoryName = "typescript"
	TypeScriptTSConfigFileName           = "tsconfig.json"
	TypeScriptTSConfigTemplateName       = "tsconfig_json.tmpl"
	TypeScriptPackageJSONFileName        = "package.json"
	TypeScriptPackageJSONTemplateName    = "package_json.tmpl"
	TypesScriptUtilsDirectory            = "stevie-utils"

	// Go
	GoFileTemplatesDirectoryName = "go"
	GoGoModName                  = "go.mod"
	GoGoSumName                  = "go.sum"
)

// PackageJsonArgs define the args need to generate a package.json file.
type PackageJsonArgs struct {
	Name        string
	Description string
}

func CreateProjectStructure(name, description string) error {
	// Create the application directory.
	err := utils.CreateNewDirectory(ApplicationFolder)
	if err != nil {
		return fmt.Errorf("Error creating application directory: %v", err)
	}

	// Create the controllers directory.
	controllersDirectoryPath := path.Join(ApplicationFolder, ControllersFolder)
	err = utils.CreateNewDirectory(controllersDirectoryPath)
	if err != nil {
		return err
	}

	return nil
}

// WriteOutTemplateFile writes out a template file.
func writeOutTemplateFile(destinationPath, templateName, name string, args interface{}) error {
	templatePath := path.Join(FileTemplatePath, TypeScriptFileTemplatesDirectoryName, templateName)
	destination := path.Join(destinationPath, name)

	err := utils.WriteOutTemplateToFile(templatePath, destination, args)
	return err
}
