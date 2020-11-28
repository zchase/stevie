package application

import (
	"fmt"
	"path"
	"strings"

	"github.com/zchase/stevie/pkg/auto_pulumi"
	"github.com/zchase/stevie/pkg/utils"
)

var (
	// Directory name.
	tmpLocalDirName = "tmp-local"

	// TypeScript
	typeScriptLocalFileTemplateName = "local.tmpl"
	typeScriptLocalFileOutputName   = "index.ts"
)

type TypeScriptLocalRoute struct {
	ImportName   string
	FunctionName string
	Path         string
	Method       string
}

type TypeScriptLocalImport struct {
	Name string
	File string
}

type TypeScriptLocalArguments struct {
	Imports []TypeScriptLocalImport
	Routes  []TypeScriptLocalRoute
	Port    int
}

// createTypeScriptLocalRoutes creates the args for generating the routes locally.
func createTypeScriptLocalArgs(routes []auto_pulumi.APIRoute, port int) TypeScriptLocalArguments {
	var localImports []TypeScriptLocalImport
	var localRoutes []TypeScriptLocalRoute

	for _, route := range routes {
		localImports = append(localImports, TypeScriptLocalImport{
			Name: route.Name,
			File: strings.Split(route.HandlerFile, ".")[0],
		})

		for _, method := range route.Methods {
			functionName := utils.SentenceToCamelCase(strings.Join([]string{route.Name, method}, " "))

			localRoutes = append(localRoutes, TypeScriptLocalRoute{
				ImportName:   route.Name,
				FunctionName: functionName,
				Method:       strings.ToLower(method),
				Path:         route.Route,
			})
		}
	}

	return TypeScriptLocalArguments{
		Port:    port,
		Routes:  localRoutes,
		Imports: localImports,
	}
}

// RunTypeScriptRoutesLocally runs the API routes locally.
func RunTypeScriptRoutesLocally(routes []auto_pulumi.APIRoute, port int) error {
	// Create the temporary directory for the file to run lcally.
	tmp := utils.TemporaryDirectory{Name: tmpLocalDirName}
	err := tmp.Create()
	if err != nil {
		return fmt.Errorf("Error creating temporary directory %s: %v", tmpLocalDirName, err)
	}

	// Listen for the program to be terminated and if it is clean up the temp directory.
	utils.ListenForProgramClose(func() {
		utils.Print("Server is exiting.")
		tmp.Clean()
		utils.Print("Server has exited successfully.")
	})

	// Generate the template contents for the local file.
	templatePath := path.Join(FileTemplatePath, TypeScriptFileTemplatesDirectoryName, typeScriptLocalFileTemplateName)
	localFilePath := path.Join(tmpLocalDirName, typeScriptLocalFileOutputName)
	templateArgs := createTypeScriptLocalArgs(routes, port)
	err = utils.WriteOutTemplateToFile(templatePath, localFilePath, templateArgs)
	if err != nil {
		return fmt.Errorf("Error generating template file contents: %v", err)
	}

	// Run the file.
	err = utils.RunCommandWithOutput("yarn", []string{"run-local"})
	if err != nil {
		return fmt.Errorf("Error running local server: %v", err)
	}

	// Clean up the temp directory and return no error.
	tmp.Clean()
	return nil
}
