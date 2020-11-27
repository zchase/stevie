package auto_pulumi

import (
	"fmt"
	"path"

	"github.com/zchase/stevie/pkg/utils"
)

// PackageTypeScriptLambda packages a TypeScript Lambda for distribution.
func PackageTypeScriptLambda(tmpDirectoryName string, name string) (string, error) {
	// If the zip file exists we can just move on.
	zipOutputPath := path.Join(tmpDirectoryName, fmt.Sprintf("%s.zip", name))
	lambdaDirectoryPath := path.Join(tmpDirectoryName, name)
	zipExists, err := utils.DoesFileExist(lambdaDirectoryPath)
	if err != nil {
		return "", nil
	}
	if zipExists {
		return zipOutputPath, nil
	}

	// Create the temp directory for packaging the lambda if it doesn't exist.
	lambdaDirectoryExists, err := utils.DoesFileExist(lambdaDirectoryPath)
	if err != nil {
		return "", err
	}
	if !lambdaDirectoryExists {
		err := utils.CreateNewDirectory(lambdaDirectoryPath)
		if err != nil {
			return "", err
		}
	}

	// Copy the files to package into the lambda directory.
	handlerFileSource := fmt.Sprintf("bin/controllers/%s.js", name)
	handlerFileDestination := path.Join(lambdaDirectoryPath, fmt.Sprintf("%s.js", name))
	err = utils.CopyFile(handlerFileSource, handlerFileDestination)
	if err != nil {
		return "", err
	}

	// Add the helper files.
	helperFiles := [2]string{"controller_builder.js", "server_args.js"}
	for _, file := range helperFiles {
		source := fmt.Sprintf("bin/controllers/%s", file)
		destination := path.Join(lambdaDirectoryPath, file)
		err = utils.CopyFile(source, destination)
		if err != nil {
			return "", err
		}
	}

	// Add the package.json for the project and install the dependencies.
	packageJSONDestination := path.Join(lambdaDirectoryPath, "package.json")
	err = utils.CopyFile("package.json", packageJSONDestination)
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
