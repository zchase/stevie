package cmd

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/zchase/stevie/pkg/application"
	"github.com/zchase/stevie/pkg/auto_pulumi"
	"github.com/zchase/stevie/pkg/utils"
	"gopkg.in/yaml.v2"
)

var baseConfigName = "base"

// Represents the values for the base config file.
type ApplicationConfig struct {
	Name         string
	DashCaseName string
	Description  string
	Routes       []auto_pulumi.APIRoute
}

// Represents the values for an environment config file.
type EnvironmentConfigFile struct {
	Name        string
	Environment string
}

// WriteOutBaseEnvironmentConfigFile writes out the base config file for a given
// environment.
func (c *ApplicationConfig) WriteOutBaseEnvironmentConfigFile(configPath, env string) error {
	// Assign the config values.
	configValues := EnvironmentConfigFile{
		Name:        c.DashCaseName,
		Environment: env,
	}

	// Convert the struct to a YAML bytes
	contents, err := yaml.Marshal(&configValues)
	if err != nil {
		return err
	}

	// Create the config file.
	envFileName := fmt.Sprintf("%s.yaml", env)
	return utils.WriteNewFile(configPath, envFileName, string(contents))
}

// WriteOutBaseConfigFile writes out the base config file.
func (c *ApplicationConfig) WriteOutBaseConfigFile(configPath string) error {
	baseConfigFileName := fmt.Sprintf("%s.yaml", baseConfigName)

	// Create the base config contents.
	contents, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	// Write out the base config file.
	return utils.WriteNewFile(configPath, baseConfigFileName, string(contents))
}

// AddAPIRouteToConfig adds an API Route to the base config.
func AddAPIRouteToConfig(configPath, name, route, pathToHandlerFiles string, corsEnabled bool) error {
	// Read in the base config file.
	baseConfig, err := ReadBaseConfig(configPath)
	if err != nil {
		return err
	}

	// If the routes doesn't exist let's add this new route as the first object
	// otherwise we append the new route.
	if baseConfig.Routes == nil {
		baseConfig.Routes = []auto_pulumi.APIRoute{
			application.CreateAPIRoute(name, route, pathToHandlerFiles, corsEnabled),
		}
	} else {
		baseConfig.Routes = append(baseConfig.Routes, application.CreateAPIRoute(name, route, pathToHandlerFiles, corsEnabled))
	}

	// Write out the new config file.
	return baseConfig.WriteOutBaseConfigFile(configPath)
}

// ReadBaseConfig reads out the base config.
func ReadBaseConfig(configPath string) (ApplicationConfig, error) {
	// Read in the config file.
	baseConfigPath := path.Join(configPath, fmt.Sprintf("%s.yaml", baseConfigName))
	configBytes, err := ioutil.ReadFile(baseConfigPath)
	if err != nil {
		return ApplicationConfig{}, err
	}

	// Unmarshal the config file.
	var appConfig ApplicationConfig
	err = yaml.Unmarshal(configBytes, &appConfig)
	return appConfig, err
}

// CreateApplicationConfig creates the application config.
func CreateApplicationConfig(configPath, name, description string, environments []string) (ApplicationConfig, error) {
	var err error

	// If either the name or description is empty prompt the user
	// for the values.
	if name == "" {
		name, err = utils.PromptRequiredString("Project Name")
		if err != nil {
			return ApplicationConfig{}, err
		}
	}

	if description == "" {
		description, err = utils.PromptRequiredString("Project Description")
		if err != nil {
			return ApplicationConfig{}, err
		}
	}

	// Create the config struct.
	config := ApplicationConfig{
		Name:         name,
		DashCaseName: utils.SentenceToDashCase(name),
		Description:  description,
	}

	// Create the config directory.
	err = utils.CreateNewDirectory(configPath)
	if err != nil {
		return ApplicationConfig{}, err
	}

	// Create the base config file.
	err = config.WriteOutBaseConfigFile(configPath)
	if err != nil {
		return ApplicationConfig{}, err
	}

	// Create the environment config files.
	for _, env := range environments {
		config.WriteOutBaseEnvironmentConfigFile(configPath, env)
	}

	return config, nil
}
