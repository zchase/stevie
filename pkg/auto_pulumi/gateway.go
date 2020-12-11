package auto_pulumi

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/apigateway"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func CreateNewAPIGateway(ctx *pulumi.Context, apiName string) (*apigateway.RestApi, error) {
	gateway, err := apigateway.NewRestApi(ctx, apiName, &apigateway.RestApiArgs{
		Name: pulumi.String(apiName),
		Policy: pulumi.String(`{
"Version": "2012-10-17",
"Statement": [
{
  "Action": "sts:AssumeRole",
  "Principal": {
	"Service": "lambda.amazonaws.com"
  },
  "Effect": "Allow",
  "Sid": ""
},
{
  "Action": "execute-api:Invoke",
  "Resource": "*",
  "Principal": "*",
  "Effect": "Allow",
  "Sid": ""
}
]
}`)})
	if err != nil {
		return nil, fmt.Errorf("Error creating API Gateway: %v", err)
	}

	return gateway, nil
}
