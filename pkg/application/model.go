package application

import (
	"fmt"
	"path"

	"github.com/zchase/stevie/pkg/utils"
)

type ModelArgs struct {
	Name         string
	HashKeyName  string
	HashKeyType  string
	RangeKeyName string
	RangeKeyType string
}

// CreateNewModel creates a new model.
func CreateNewModel(name string, language string, args ModelArgs) error {
	var fileExtension string
	switch language {
	case TypeScriptLanguage:
		fileExtension = "ts"
		break
	default:
		return fmt.Errorf("Unsupported language provided")
	}

	modelTemplateFilePath := path.Join(FileTemplatePath, TypeScriptLanguage, "model.tmpl")
	modelFilePath := fmt.Sprintf("%s.%s", path.Join(ApplicationFolder, ModelsFolder, name), fileExtension)
	err := utils.WriteOutTemplateToFile(modelTemplateFilePath, modelFilePath, args)
	if err != nil {
		return utils.NewErrorMessage("Error writing new model file", err)
	}

	return nil
}
