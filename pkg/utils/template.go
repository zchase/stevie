package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/markbates/pkger"
)

func WriteOutTemplateToFile(templatePath string, filePath string, args interface{}) error {
	tpl := template.New("")

	// Parse the template.
	templateContents, err := pkger.Open(templatePath)
	if err != nil {
		return fmt.Errorf("[WriteOutTemplateToFile]: Error opening template: %v", err)
	}
	defer templateContents.Close()

	// Read template to bytes
	templateBytes, err := ioutil.ReadAll(templateContents)
	if err != nil {
		return fmt.Errorf("[WriteOutTemplateToFile]: Error reading template: %v", err)
	}

	// Parse the template.
	parsedTemplate, err := tpl.Parse(string(templateBytes))
	if err != nil {
		return fmt.Errorf("[WriteOutTemplateToFile]: Error parsing template: %v", err)
	}

	// Create the file.
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("[WriteOutTemplateToFile]: Error creating template file: %v", err)
	}

	// Write out the template to the new file.
	err = parsedTemplate.Execute(file, args)
	if err != nil {
		return fmt.Errorf("[WriteOutTemplateToFile]: Error writing out file from template: %v", err)
	}

	return nil
}
