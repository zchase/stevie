package auto_pulumi

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lambda"
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
		lowerCaseName := strings.ToLower(contentName)
		switch lowerCaseName {
		case "get", "post", "put", "delete":
			methods = append(methods, lowerCaseName)
		}
	}

	// If there are no valid methods we should return an error.
	if len(methods) == 0 {
		return nil, fmt.Errorf("No valid methods found in the controller folder.")
	}
	return methods, nil
}

func CreateAPI(
	ctx *pulumi.Context, projectName, environment string, routes []APIRoute,
	tableNames []DynamoDBTable,
) ([]APIEndpoint, error) {
	apiName := fmt.Sprintf("%s-api", projectName)

	// Create the API Gateway
	var functions []APIEndpointFunction
	var apiResources []*apigateway.Resource
	var permissions []*lambda.Permission
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

		endpointData, err := CreateAPIEndpoint(ctx, gateway, environment, route, routeMethods, tableNames)
		if err != nil {
			return nil, err
		}

		endpoint := APIEndpoint{
			Name: route.Name,
			URL:  endpointData.EndpointUrl,
		}
		endpoints = append(endpoints, endpoint)
		functions = append(functions, endpointData.Functions...)
		apiResources = append(apiResources, endpointData.ApiResource)
		permissions = append(permissions, endpointData.Permissions...)
	}

	err = CreateApiGatewayDeployment(ctx, functions, apiResources, gateway, permissions, environment)
	if err != nil {
		return nil, fmt.Errorf("Error creating API deployment: %v", err)
	}

	return endpoints, nil
}
