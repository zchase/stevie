package auto_pulumi

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// APIRoute represents and API Route in the config.
type APIRoute struct {
	Name        string
	Route       string
	HandlerFile string
	Methods     []string
}

type APIEndpoint struct {
	Name string
	URL  pulumi.StringOutput
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
		endpointURL, err := CreateAPIEndpoint(ctx, gateway, environment, route)
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
