package auto_pulumi

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/zchase/stevie/pkg/utils"
)

// APIRoute represents and API Route in the config.
type APIRoute struct {
	Name        string
	Route       string
	PathToFiles string
	CorsEnabled bool
}

type APIEndpoint struct {
	Name string
	URL  pulumi.StringOutput
}

// readRoutesFromControllerDirectory reads the different controller methods in a given
// controller directory.
func readRoutesFromControllerDirectory(controllerDirectoryPath string) ([]string, error) {
	var methods []string

	controllerDirectoryContents, err := utils.ReadDirectoryContents(controllerDirectoryPath)
	if err != nil {
		return nil, err
	}

	for _, contentName := range controllerDirectoryContents {
		switch contentName {
		case "get", "post", "put", "delete":
			methods = append(methods, contentName)
		}
	}

	// If there are no valid methods we should return an error.
	if len(methods) == 0 {
		return nil, fmt.Errorf("No valid methods found in the controller folder.")
	}
	return methods, nil
}

func CreateAPI(ctx *pulumi.Context, projectName, environment string, routes []APIRoute) ([]APIEndpoint, error) {
	apiName := fmt.Sprintf("%s-api", projectName)

	// Create the API Gateway
	gateway, err := CreateNewAPIGateway(ctx, apiName)
	if err != nil {
		return nil, err
	}

	// Create the endpoints for the controllers.
	var endpoints []APIEndpoint
	for _, route := range routes {
		routeMethods, err := readRoutesFromControllerDirectory(route.PathToFiles)
		if err != nil {
			return nil, err
		}

		endpointURL, err := CreateAPIEndpoint(ctx, gateway, environment, route, routeMethods)
		if err != nil {
			return nil, err
		}

		endpoint := APIEndpoint{
			Name: route.Name,
			URL:  endpointURL,
		}
		endpoints = append(endpoints, endpoint)
	}

	return endpoints, nil
}
