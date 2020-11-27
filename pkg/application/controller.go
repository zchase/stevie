package application

import (
	"fmt"
	"path"
	"strings"

	"github.com/zchase/stevie/pkg/utils"
)

type TypeScriptControllerFileArgs struct {
	FunctionName string
	Methods      []TypeScriptRouteMethodArgs
}

type TypeScriptRouteMethodArgs struct {
	Name         string
	FunctionName string
	HandlerName  string
}

// createTypeScriptRouteMethodArgs creates the args for creating the routes in the controller
// file for TypeScript.
func createTypeScriptRouteMethodArgs(controllerName string, methods []string) []TypeScriptRouteMethodArgs {
	var result []TypeScriptRouteMethodArgs
	for _, method := range methods {
		// Create the function name.
		functionNameParts := fmt.Sprintf("%s %s", controllerName, method)
		functionName := utils.SentenceToCamelCase(functionNameParts)

		// Create the handler name.
		handlerName := fmt.Sprintf("%sHandler", strings.ToLower(method))

		result = append(result, TypeScriptRouteMethodArgs{
			Name:         method,
			FunctionName: functionName,
			HandlerName:  handlerName,
		})
	}

	return result
}

// CreateNewTypeScriptController creates a new TypeScript controller.
func CreateNewTypeScriptController(name string, methods []string) (string, error) {
	// Create the file paths.
	controllerTemplatePath := path.Join(FileTemplatePath, TypeScriptFileTemplatesDirectoryName, "controller.tmpl")
	controllerFileName := fmt.Sprintf("%s.ts", name)
	controllerFilePath := path.Join(TypeScriptAppDirectoryName, TypeScriptControllersDirectoryName, controllerFileName)

	// Create the method args.
	methodArgs := createTypeScriptRouteMethodArgs(name, methods)

	err := utils.WriteOutTemplateToFile(controllerTemplatePath, controllerFilePath, TypeScriptControllerFileArgs{
		FunctionName: name,
		Methods:      methodArgs,
	})
	if err != nil {
		return "", err
	}

	return controllerFilePath, nil
}
