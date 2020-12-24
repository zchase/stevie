package auto_pulumi

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

var tmpDirName = "tmp"

func createAPIGatewayResource(ctx *pulumi.Context, gateway *apigateway.RestApi, name string, route string) (*apigateway.Resource, error) {
	apiName := fmt.Sprintf("%s-api-resource", name)
	pathParts := strings.Split(route, "/")
	pathPart := pathParts[len(pathParts)-1]

	apiresource, err := apigateway.NewResource(ctx, apiName, &apigateway.ResourceArgs{
		RestApi:  gateway.ID(),
		PathPart: pulumi.String(pathPart),
		ParentId: gateway.RootResourceId,
	}, pulumi.DependsOn([]pulumi.Resource{gateway}))
	if err != nil {
		return nil, fmt.Errorf("Error creating API Gateway Resource for %s: %v", name, err)
	}

	return apiresource, nil
}

func createAPIGatewayRouteMethods(
	ctx *pulumi.Context, apiResource *apigateway.Resource, gateway *apigateway.RestApi,
	name string, methods []APIEndpointFunction,
) error {
	gatewayID := gateway.ID()
	resourceID := apiResource.ID()

	for _, method := range methods {
		methodResourceName := fmt.Sprintf("%s-api-%s-method", name, method.Method)

		_, err := apigateway.NewMethod(ctx, methodResourceName, &apigateway.MethodArgs{
			HttpMethod:    pulumi.String(strings.ToUpper(method.Method)),
			Authorization: pulumi.String("NONE"),
			RestApi:       gatewayID,
			ResourceId:    resourceID,
		}, pulumi.DependsOn([]pulumi.Resource{gateway, apiResource}))
		if err != nil {
			return fmt.Errorf("Error creating %s: %v", methodResourceName, err)
		}
	}

	return nil
}

func createAPIGatewayIntegration(
	ctx *pulumi.Context, apiResource *apigateway.Resource, gateway *apigateway.RestApi,
	name string, regionName string, awsAccountID string, methods []APIEndpointFunction,
) ([]*lambda.Permission, error) {
	var lambdaPermissions []*lambda.Permission

	// Add an integration to the API Gateway.
	// This makes communication between the API Gateway and the Lambda function work
	for _, method := range methods {
		function := method.Function
		integrationName := fmt.Sprintf("%s-%s-lambda-integration", name, method.Method)

		_, err := apigateway.NewIntegration(ctx, integrationName, &apigateway.IntegrationArgs{
			HttpMethod:            pulumi.String(strings.ToUpper(method.Method)),
			IntegrationHttpMethod: pulumi.String("POST"),
			ResourceId:            apiResource.ID(),
			RestApi:               gateway.ID(),
			Type:                  pulumi.String("AWS_PROXY"),
			Uri:                   function.InvokeArn,
		}, pulumi.DependsOn([]pulumi.Resource{gateway, apiResource, function}))
		if err != nil {
			return nil, err
		}

		// Add a resource based policy to the Lambda function.
		// This is the final step and allows AWS API Gateway to communicate with the AWS Lambda function
		permissionName := fmt.Sprintf("%s-%s-api-permission", name, method.Method)

		permission, err := lambda.NewPermission(ctx, permissionName, &lambda.PermissionArgs{
			Action:    pulumi.String("lambda:InvokeFunction"),
			Function:  function.Name,
			Principal: pulumi.String("apigateway.amazonaws.com"),
			SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/*/*/*", regionName, awsAccountID, gateway.ID()),
		}, pulumi.DependsOn([]pulumi.Resource{gateway, apiResource, function}))
		if err != nil {
			return nil, err
		}

		lambdaPermissions = append(lambdaPermissions, permission)
	}

	return lambdaPermissions, nil
}

func enableCORSOnResource(ctx *pulumi.Context, apiResource *apigateway.Resource, gateway *apigateway.RestApi, route APIRoute) error {
	apiGatewayId := gateway.ID()
	apiGatewayResourceId := apiResource.ID()

	// Create the CORS Method.
	methodName := fmt.Sprintf("%s-api-options-cors-method", route.Name)
	corsMethod, err := apigateway.NewMethod(ctx, methodName, &apigateway.MethodArgs{
		Authorization: pulumi.String("NONE"),
		HttpMethod:    pulumi.String("OPTIONS"),
		ResourceId:    apiGatewayResourceId,
		RestApi:       apiGatewayId,
	})
	if err != nil {
		return fmt.Errorf("Error creating CORS API Gateway Method: %v", err)
	}

	// Create the integration.
	integrationName := fmt.Sprintf("%s-options-cors-lambda-integration", route.Name)
	integration, err := apigateway.NewIntegration(ctx, integrationName, &apigateway.IntegrationArgs{
		HttpMethod: corsMethod.HttpMethod,
		ResourceId: apiGatewayResourceId,
		RestApi:    apiGatewayId,
		Type:       pulumi.String("MOCK"),
		RequestTemplates: pulumi.StringMap{
			"application/json": pulumi.String("{ statusCode: 200 }"),
		},
	}, pulumi.DependsOn([]pulumi.Resource{corsMethod}))
	if err != nil {
		return fmt.Errorf("Error creating API Integrations for CORS method: %v", err)
	}

	// Create the method response.
	responseName := fmt.Sprintf("%s-options-cors-api-response", route.Name)
	response200, err := apigateway.NewMethodResponse(ctx, responseName, &apigateway.MethodResponseArgs{
		HttpMethod: corsMethod.HttpMethod,
		ResourceId: apiGatewayResourceId,
		RestApi:    apiGatewayId,
		StatusCode: pulumi.String("200"),
		ResponseParameters: pulumi.BoolMap{
			"method.response.header.Access-Control-Allow-Origin":  pulumi.Bool(true),
			"method.response.header.Access-Control-Allow-Methods": pulumi.Bool(true),
			"method.response.header.Access-Control-Allow-Headers": pulumi.Bool(true),
		},
	}, pulumi.DependsOn([]pulumi.Resource{integration}))
	if err != nil {
		return fmt.Errorf("Error creating response for CORS method: %v", err)
	}

	// Create the integration response.
	integrationResponseName := fmt.Sprintf("%s-options-cors-integration-response", route.Name)
	_, err = apigateway.NewIntegrationResponse(ctx, integrationResponseName, &apigateway.IntegrationResponseArgs{
		HttpMethod:        corsMethod.HttpMethod,
		ResourceId:        apiGatewayResourceId,
		RestApi:           apiGatewayId,
		StatusCode:        response200.StatusCode,
		ResponseTemplates: pulumi.StringMap{"application/json": pulumi.String("{}")},
		ResponseParameters: pulumi.StringMap{
			"method.response.header.Access-Control-Allow-Origin":  pulumi.String("'*'"),
			"method.response.header.Access-Control-Allow-Methods": pulumi.String("'*'"),
			"method.response.header.Access-Control-Allow-Headers": pulumi.String("'*'"),
		},
	})
	if err != nil {
		return fmt.Errorf("Error creating the CORS Integration Response: %v", err)
	}

	return nil
}

func createApiGatewayDeployment(
	ctx *pulumi.Context, functions []APIEndpointFunction, apiResource *apigateway.Resource, gateway *apigateway.RestApi,
	permissions []*lambda.Permission, name string, environment string,
) error {
	apiDeploymentName := fmt.Sprintf("%s-api-deployment", name)
	stage := pulumi.String(environment)

	// Create the depends-on array
	dependsOn := []pulumi.Resource{gateway, apiResource}
	for _, function := range functions {
		dependsOn = append(dependsOn, function.Function)
	}

	for _, permission := range permissions {
		dependsOn = append(dependsOn, permission)
	}

	_, err := apigateway.NewDeployment(ctx, apiDeploymentName, &apigateway.DeploymentArgs{
		RestApi:          gateway.ID(),
		StageDescription: stage,
		StageName:        stage,
	}, pulumi.DependsOn(dependsOn))
	if err != nil {
		return fmt.Errorf("Error creating API Gateway Deployment for %s: %v", name, err)
	}

	return nil
}

type APIEndpointFunction struct {
	Function *lambda.Function
	Method   string
}

func CreateAPIEndpoint(
	ctx *pulumi.Context, gateway *apigateway.RestApi, environment string,
	route APIRoute, methods []string, tableNames []DynamoDBTable,
) (pulumi.StringOutput, error) {
	// Get the AWS account.
	account, err := aws.GetCallerIdentity(ctx)
	if err != nil {
		return pulumi.StringOutput{}, fmt.Errorf("Error getting AWS identity: %v", err)
	}

	// Get the AWS region.
	region, err := aws.GetRegion(ctx, &aws.GetRegionArgs{})
	if err != nil {
		return pulumi.StringOutput{}, fmt.Errorf("Error getting AWS region: %v", err)
	}

	// Create the lambdas functions.
	var lambdaFunctions []APIEndpointFunction
	for _, method := range methods {
		function, err := CreateRouteHandler(ctx, route, method, tableNames)
		if err != nil {
			return pulumi.StringOutput{}, err
		}

		lambdaFunctions = append(lambdaFunctions, APIEndpointFunction{
			Function: function,
			Method:   method,
		})
	}

	// Add a resource to the API Gateway.
	apiResource, err := createAPIGatewayResource(ctx, gateway, route.Name, route.Route)
	if err != nil {
		return pulumi.StringOutput{}, err
	}

	// Add the methods to the API Gateway.
	err = createAPIGatewayRouteMethods(ctx, apiResource, gateway, route.Name, lambdaFunctions)
	if err != nil {
		return pulumi.StringOutput{}, err
	}

	// Add integrations for each lambda to the API Gateway.
	permissions, err := createAPIGatewayIntegration(
		ctx, apiResource, gateway, route.Name, region.Name,
		account.Id, lambdaFunctions,
	)
	if err != nil {
		return pulumi.StringOutput{}, err
	}

	// Create a new deployment
	err = createApiGatewayDeployment(
		ctx, lambdaFunctions, apiResource, gateway, permissions, route.Name,
		environment,
	)
	if err != nil {
		return pulumi.StringOutput{}, err
	}

	endpointURL := pulumi.Sprintf("https://%s.execute-api.%s.amazonaws.com/%s%s", gateway.ID(), region.Name, environment, route.Route)

	return endpointURL, nil
}
