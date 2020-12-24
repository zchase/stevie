package utils

import (
	"encoding/json"
	"strings"
)

type JSONSchemaItemProperties struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	HashKey     bool   `json:"hashKey"`
	RangeKey    bool   `json:"rangeKey"`
}

type JSONSchemaItem struct {
	Properties map[string]JSONSchemaItemProperties `json:"properties"`
	Type       string
}

type JSONSchemaDefinition map[string]JSONSchemaItem

type JSONSchema struct {
	Schema      string               `json:"$schema"`
	Definitions JSONSchemaDefinition `json:"definitions"`
}

func GenerateModelSchemaFromFile(filePath string) (JSONSchema, error) {
	// language, err := DetectFileLanguageFromExtension(filePath)
	// if err != nil {
	// 	return nil, NewErrorMessage("Error detecting language from file", err)
	// }

	language := "typescript"

	switch language {
	case "typescript":
		schema, err := generateJSONSchemaFromTypeScriptFile(filePath)
		if err != nil {
			return JSONSchema{}, err
		}
		return schema, nil
	default:
		return JSONSchema{}, NewErrorMessage("Unsupported language detected", nil)
	}
}

func generateJSONSchemaFromTypeScriptFile(filePath string) (JSONSchema, error) {
	// Generate the schema.
	schemaString, err := RunCommand("typescript-json-schema", []string{filePath, "*", "--validationKeywords", "hashKey", "rangeKey"})
	if err != nil {
		return JSONSchema{}, NewErrorMessage("Error generating JSON schema from model file", err)
	}

	// Unmarshal the JSON
	var schema JSONSchema
	err = json.Unmarshal([]byte(strings.TrimSpace(schemaString)), &schema)
	if err != nil {
		return JSONSchema{}, NewErrorMessage("Error unmarshaling JSON schema", err)
	}

	return schema, nil
}
