package application

import (
	"path"

	"github.com/zchase/stevie/pkg/utils"
)

var (
	// Config
	ApplicationConfigPath = "config"

	// Directories
	ApplicationFolder = "app"
	ControllersFolder = "controllers"
	ModelsFolder      = "models"

	// File template paths
	FileTemplatePath = "/pkg/application/file_templates"
	LocalPackagePath = "/lib"

	// TypeScript
	TypeScriptLanguage                   = "typescript"
	TypeScriptFileTemplatesDirectoryName = "typescript"
	TypeScriptTSConfigFileName           = "tsconfig.json"
	TypeScriptTSConfigTemplateName       = "tsconfig_json.tmpl"
	TypeScriptPackageJSONFileName        = "package.json"
	TypeScriptPackageJSONTemplateName    = "package_json.tmpl"
	TypesScriptUtilsDirectory            = "stevie-utils"

	// Go
	GoLanguage                   = "go"
	GoFileTemplatesDirectoryName = "go"
	GoGoModName                  = "go.mod"
	GoGoSumName                  = "go.sum"

	// Dotnet
	DotNetLangauge                  = "dotnet"
	DotNetFileTemplateDirectoryName = "dotnet"
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
		return utils.NewErrorMessage("Error creating application directory", err)
	}

	// Create the controllers directory.
	controllersDirectoryPath := path.Join(ApplicationFolder, ControllersFolder)
	err = utils.CreateNewDirectory(controllersDirectoryPath)
	if err != nil {
		return utils.NewErrorMessage("Error creating controllers directory", err)
	}

	// Create the models directory.
	modelsDirectoryPath := path.Join(ApplicationFolder, ModelsFolder)
	err = utils.CreateNewDirectory(modelsDirectoryPath)
	if err != nil {
		return utils.NewErrorMessage("Error creating models directory", err)
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
