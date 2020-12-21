package auto_pulumi

import (
	"fmt"
	"os"
	"path"

	"github.com/zchase/stevie/pkg/utils"
)

// PackageDotNetLambda packages a DotNet Lambda for distribution.
func PackageDotNetLambda(tmpDirectoryName, routeName, method string) (string, error) {
	name := fmt.Sprintf("%s-%s-handler", routeName, method)
	lambdaDirectoryPath := path.Join(tmpDirectoryName, name)

	// Create the temp directory for packaging the lambda if it doesn't exist.
	err := utils.CreateNewDirectory(lambdaDirectoryPath)
	if err != nil {
		return "", err
	}

	dotNetControllerDirectory := path.Join("app", "controllers", routeName, method)
	err = utils.CopyDirectory(dotNetControllerDirectory, lambdaDirectoryPath, nil)
	if err != nil {
		return "", utils.NewErrorMessage("Error copy dotent lambda files", err)
	}

	_, err = utils.RunCommand("dotnet", []string{"publish", lambdaDirectoryPath})
	if err != nil {
		return "", utils.NewErrorMessage("Error building dotnet lambda", err)
	}

	dotNetLambdaEntryPoint := path.Join(lambdaDirectoryPath, "bin", "Debug", "netcoreapp3.1", "publish")
	return dotNetLambdaEntryPoint, nil
}

// PackageGoLambda packages a Go Lambda for distribution.
func PackageGoLambda(tmpDirectoryName, routeName, method string) (string, error) {
	name := fmt.Sprintf("%s-%s-handler", routeName, method)
	lambdaDirectoryPath := path.Join(tmpDirectoryName, name)

	// Create the temp directory for packaging the lambda if it doesn't exist.
	err := utils.CreateNewDirectory(lambdaDirectoryPath)
	if err != nil {
		return "", err
	}

	goFilePath := path.Join("app", "controllers", routeName, method, fmt.Sprintf("%s.go", method))
	goFileOutputPath := path.Join(lambdaDirectoryPath, name)
	err = os.Setenv("GOOS", "linux")
	if err != nil {
		return "", err
	}
	err = os.Setenv("GOARCH", "amd64")
	if err != nil {
		return "", err
	}
	_, err = utils.RunCommand("go", []string{"build", "-o", goFileOutputPath, goFilePath})
	if err != nil {
		return "", utils.NewErrorMessage("Error building lambda", err)
	}
	err = os.Setenv("GOOS", "")
	if err != nil {
		return "", err
	}
	err = os.Setenv("GOARCH", "")
	if err != nil {
		return "", err
	}

	zipFilePath := path.Join(lambdaDirectoryPath, fmt.Sprintf("%s.zip", name))
	_, err = utils.RunCommand("zip", []string{"-j", zipFilePath, goFileOutputPath})
	if err != nil {
		return "", utils.NewErrorMessage("Error zipping lambda", err)
	}

	return path.Join(lambdaDirectoryPath, fmt.Sprintf("%s.zip", name)), nil
}

// PackageTypeScriptLambda packages a TypeScript Lambda for distribution.
func PackageTypeScriptLambda(tmpDirectoryName string, routeName, method string) (string, error) {
	name := fmt.Sprintf("%s-%s-handler", routeName, method)
	zipOutputPath := path.Join(tmpDirectoryName, fmt.Sprintf("%s.zip", name))
	lambdaDirectoryPath := path.Join(tmpDirectoryName, name)

	// Create the temp directory for packaging the lambda if it doesn't exist.
	err := utils.CreateNewDirectory(lambdaDirectoryPath)
	if err != nil {
		return "", err
	}

	// Build the lambda. Technically this is repeated for method but eventually it would
	// be good to support a package.json per method so we only package up what is need.
	routeCodeDir := fmt.Sprintf("app/controllers/%s", routeName)
	_, err = utils.RunCommand("yarn", []string{"--cwd", routeCodeDir})
	_, err = utils.RunCommand("yarn", []string{"--cwd", routeCodeDir, "build"})

	// Copy the files to package into the lambda directory.
	handlerFileSource := fmt.Sprintf("%s/bin/%s/%s.js", routeCodeDir, method, method)
	handlerFileDestination := path.Join(lambdaDirectoryPath, fmt.Sprintf("%s.js", name))
	err = utils.CopyFile(handlerFileSource, handlerFileDestination)
	if err != nil {
		return "", err
	}

	// Add the utils package.
	utilsName := path.Join("app", "controllers", routeName, "stevie-utils")
	utilsPath := path.Join(lambdaDirectoryPath, "stevie-utils")
	err = utils.CreateNewDirectory(utilsPath)
	if err != nil {
		return "", fmt.Errorf("Error creating utils directory for Lambda: %v", err)
	}

	err = utils.CopyDirectory(utilsName, utilsPath, []string{"node_modules"})
	if err != nil {
		return "", fmt.Errorf("Error copy in TypeScript utilities: %v", err)
	}

	// Add the package.json for the project and install the dependencies.
	packageJSONDestination := path.Join(lambdaDirectoryPath, "package.json")
	packageJSONOrigin := path.Join("app", "controllers", routeName, "package.json")
	err = utils.CopyFile(packageJSONOrigin, packageJSONDestination)
	if err != nil {
		return "", err
	}

	// Install the lambda dependencies.
	_, err = utils.RunCommand("yarn", []string{"--cwd", lambdaDirectoryPath, "install"})
	if err != nil {
		return "", err
	}

	// Zip the files and write out the zip file.
	err = utils.ZipDirectory(lambdaDirectoryPath, zipOutputPath)
	if err != nil {
		return "", err
	}

	// Return the path to the zip file.
	return zipOutputPath, nil
}
