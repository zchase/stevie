package application

import (
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/zchase/stevie/pkg/auto_pulumi"
)

// CreateAPIRoute creates an API Route.
func CreateAPIRoute(name, route, pathToFiles string, corsEnabled bool) auto_pulumi.APIRoute {
	return auto_pulumi.APIRoute{
		Name:        name,
		Route:       route,
		CorsEnabled: corsEnabled,
		PathToFiles: pathToFiles,
	}
}

// BuildAPIRoutes builds the API Routes.
func BuildAPIRoutes(
	ctx *pulumi.Context, projectName, environment string, routes []auto_pulumi.APIRoute,
	tableNames []auto_pulumi.DynamoDBTable,
) error {
	endpointURLS, err := auto_pulumi.CreateAPI(ctx, projectName, environment, routes, tableNames)
	if err != nil {
		return err
	}

	for _, url := range endpointURLS {
		ctx.Export(url.Name, url.URL)
	}

	return nil
}
