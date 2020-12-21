package application

import (
	"fmt"
	"path"
	"strings"

	"github.com/zchase/stevie/pkg/utils"
)

const (
	TypeScriptControllerLanguage = "typescript"
	GoControllerLanguage         = "go"
	DotNetControllerLanguage     = "dotnet"
)

// CreateNewController creates a new controller.
func CreateNewController(name string, methods []string, language string) (string, error) {
	// Create the main controller directory.
	controllerDirectoryPath := path.Join(ApplicationFolder, ControllersFolder, name)
	err := utils.CreateNewDirectory(controllerDirectoryPath)
	if err != nil {
		return "", fmt.Errorf("Error creating main controller directory: %v", err)
	}

	// Create any top level files needed for the controller.
	switch language {
	case GoControllerLanguage:
		err = createGoTopLevelFiles(controllerDirectoryPath)
		if err != nil {
			return "", err
		}
		break
	case TypeScriptControllerLanguage:
		err = createTypeScriptTopLevelFiles(controllerDirectoryPath, name, "")
		if err != nil {
			return "", err
		}
		break
	case DotNetControllerLanguage:
		break
	default:
		return "", fmt.Errorf("Language not supported: %s", language)
	}

	// Loop through the methods.
	for _, method := range methods {
		// Create the handler directory.
		controllerHandlerDirectoryPath := path.Join(controllerDirectoryPath, method)
		err = utils.CreateNewDirectory(controllerHandlerDirectoryPath)
		if err != nil {
			return "", fmt.Errorf("Error creating handler directory for %s method on route %s: %v", name, method, err)
		}

		// Create the handler's files.
		switch language {
		case GoControllerLanguage:
			err = createNewGoController(controllerHandlerDirectoryPath, name, method)
			if err != nil {
				return "", err
			}
			break
		case TypeScriptControllerLanguage:
			err = createNewTypeScriptController(controllerHandlerDirectoryPath, name, method)
			if err != nil {
				return "", err
			}
			break
		case DotNetControllerLanguage:
			err = createNewDotNetController(controllerHandlerDirectoryPath, name, method)
			if err != nil {
				return "", err
			}
			break
		}
	}

	return controllerDirectoryPath, nil
}

type DotNetControllerFileArgs struct {
	Route        string
	Method       string
	FunctionName string
}

// createNewDotNetController creates a new dotnet controller.
func createNewDotNetController(dirPath, name, method string) error {
	dotNetControllerTemplatePath := path.Join(FileTemplatePath, DotNetFileTemplateDirectoryName, "controller.tmpl")
	dotNetControllerFileName := fmt.Sprintf("%s.cs", utils.DashCaseToSentenceCase(method))
	dotNetControllerFilePath := path.Join(dirPath, dotNetControllerFileName)

	err := utils.WriteOutTemplateToFile(dotNetControllerTemplatePath, dotNetControllerFilePath, DotNetControllerFileArgs{
		FunctionName: utils.DashCaseToSentenceCase(method),
		Method:       method,
		Route:        fmt.Sprintf("/%s", name),
	})
	if err != nil {
		return fmt.Errorf("Error creating controller file: %v", err)
	}

	dotNetControllerProjectTemplatePath := path.Join(FileTemplatePath, DotNetFileTemplateDirectoryName, "app_csproj.tmpl")
	dotNetControllerProjectFilePath := path.Join(dirPath, "app.csproj")
	err = utils.WriteOutTemplateToFile(dotNetControllerProjectTemplatePath, dotNetControllerProjectFilePath, nil)
	if err != nil {
		return fmt.Errorf("Error creating dotnet project file: %v", err)
	}

	return nil
}

// createGoTopLevelFiles creates the top level files for a go controller.
func createGoTopLevelFiles(dirPath string) error {
	// Create the go.mod and go.sum files.
	err := utils.WriteNewFile("", GoGoModName, "")
	if err != nil {
		return fmt.Errorf("Error creating go.mod: %v", err)
	}

	err = utils.WriteNewFile("", GoGoSumName, "")
	if err != nil {
		return fmt.Errorf("Error creating go.sum: %v", err)
	}

	return nil
}

type GoControllerFileArgs struct {
	Method string
	Route  string
}

// CreateNewGoController creates a new Go controller.
func createNewGoController(dirPath, name, method string) error {
	goControllerTemplatePath := path.Join(FileTemplatePath, GoFileTemplatesDirectoryName, "controller.tmpl")
	goControllerFileName := fmt.Sprintf("%s.go", strings.ToLower(method))
	goControllerFilePath := path.Join(dirPath, goControllerFileName)

	err := utils.WriteOutTemplateToFile(goControllerTemplatePath, goControllerFilePath, GoControllerFileArgs{
		Method: method,
		Route:  fmt.Sprintf("/%s", name),
	})
	if err != nil {
		return fmt.Errorf("Error creating controller file: %v", err)
	}

	return nil
}

// createTypeScriptTopLevelFiles creates the top level files for a TypeScript handler.
func createTypeScriptTopLevelFiles(dirPath, name, description string) error {
	// Create the package.json file.
	packageJSONArgs := PackageJsonArgs{
		Name:        name,
		Description: description,
	}
	err := writeOutTemplateFile(dirPath, TypeScriptPackageJSONTemplateName, TypeScriptPackageJSONFileName, packageJSONArgs)
	if err != nil {
		return fmt.Errorf("Error creating top level package.json: %v", err)
	}

	// Create tsconfig.json file.
	err = writeOutTemplateFile(dirPath, TypeScriptTSConfigTemplateName, TypeScriptTSConfigFileName, nil)
	if err != nil {
		return fmt.Errorf("Error creating top level tsconfig.json: %v", err)
	}

	// Copy in the utils package.
	utilsPackagePath := path.Join(LocalPackagePath, TypeScriptFileTemplatesDirectoryName)
	controllerUtilsDirectoryPath := path.Join(dirPath, TypesScriptUtilsDirectory)
	err = utils.CreateNewDirectory(controllerUtilsDirectoryPath)
	if err != nil {
		return fmt.Errorf("Error creating utils directory: %v", err)
	}

	err = utils.CopyPackagedDirectory(utilsPackagePath, controllerUtilsDirectoryPath, []string{"node_modules", "lib"})
	if err != nil {
		return fmt.Errorf("Error copy in TypeScript utilities: %v", err)
	}

	return nil
}

type TypeScriptControllerFileArgs struct {
	FunctionName string
	HandlerName  string
}

// createNewTypeScriptController creates a new TypeScript controller.
func createNewTypeScriptController(dirPath, name, method string) error {
	// Create the file paths.
	controllerTemplatePath := path.Join(FileTemplatePath, TypeScriptFileTemplatesDirectoryName, "controller.tmpl")
	controllerFileName := fmt.Sprintf("%s.ts", strings.ToLower(method))
	controllerFilePath := path.Join(dirPath, controllerFileName)

	// Create the function and handler names.
	functionNameParts := fmt.Sprintf("%s %s", utils.DashCaseToCamelCase(name), method)
	functionName := utils.SentenceToCamelCase(functionNameParts)
	handlerName := fmt.Sprintf("%sHandler", strings.ToLower(method))

	err := utils.WriteOutTemplateToFile(controllerTemplatePath, controllerFilePath, TypeScriptControllerFileArgs{
		FunctionName: functionName,
		HandlerName:  handlerName,
	})
	if err != nil {
		return fmt.Errorf("Error writing out controller file for %s method on route %s: %v", method, name, err)
	}

	return nil
}
