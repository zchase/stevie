package utils

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type DefaultConfigFile struct {
	Name        string
	Environment string
}

// createConfigFileContents creates a config file contents.
func createDefaultConfigFileContents(name string, env string) (string, error) {
	// Create the config file data.
	configFile := DefaultConfigFile{
		Name:        name,
		Environment: env,
	}

	// Marshal the struct into YAML.
	contents, err := yaml.Marshal(&configFile)
	if err != nil {
		return "", err
	}

	// Return the YAML string.
	return string(contents), nil
}

func CreateConfigFile(configDirectoryName string, name string, env string) error {
	// Create the file contents.
	configContents, err := createDefaultConfigFileContents(name, env)
	if err != nil {
		return err
	}

	// Write out the new file.
	fileName := fmt.Sprintf("%s.yaml", env)
	err = WriteNewFile(configDirectoryName, fileName, configContents)
	if err != nil {
		return err
	}

	return nil
}

func ReadConfigFile(configDirectoryName string, env string) (DefaultConfigFile, error) {
	configFileName := fmt.Sprintf("%s/%s.yaml", configDirectoryName, env)
	configFileContents, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return DefaultConfigFile{}, fmt.Errorf("Error reading config file %s: %v", configFileName, err)
	}

	var configFile DefaultConfigFile
	err = yaml.Unmarshal(configFileContents, &configFile)
	if err != nil {
		return DefaultConfigFile{}, fmt.Errorf("Error reading config file %s: %v", configFileName, err)
	}

	return configFile, nil
}
