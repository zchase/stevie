package auto_pulumi

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/zchase/stevie/pkg/utils"
)

// createIAMLambdaRole creates an IAM role for a Lambda
func createIAMLambdaRole(ctx *pulumi.Context, name string) (*iam.Role, error) {
	roleName := fmt.Sprintf("%s-task-exec-role", name)

	role, err := iam.NewRole(ctx, roleName, &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Sid": "",
				"Effect": "Allow",
				"Principal": {
					"Service": "lambda.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}]
		}`),
	})
	if err != nil {
		return nil, fmt.Errorf("Error create IAM role for %s: %v", name, err)
	}

	return role, nil
}

// createLambdaLogPolicy creates a RolePolicy for collecting logs for a Lambda.
func createLambdaLogPolicy(ctx *pulumi.Context, role *iam.Role, name string) (*iam.RolePolicy, error) {
	logPolicyName := fmt.Sprintf("%s-lambda-log-policy", name)
	logPolicy, err := iam.NewRolePolicy(ctx, logPolicyName, &iam.RolePolicyArgs{
		Role: role.Name,
		Policy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Action": [
					"logs:CreateLogGroup",
					"logs:CreateLogStream",
					"logs:PutLogEvents"
				],
				"Resource": "arn:aws:logs:*:*:*"
			}]
		}`),
	})
	if err != nil {
		return nil, fmt.Errorf("Error creating role policy for %s: %v", name, err)
	}

	return logPolicy, nil
}

// detectLambdaLanguage figures out the language of lambda by looking at the
// the file extension.
func detectLambdaLanguage(dirPath string) (string, error) {
	dirContents, err := utils.ReadDirectoryContents(dirPath)
	if err != nil {
		return "", err
	}

	// Check the first file that is not a directory.
	var fileNameParts []string
	for i, content := range dirContents {
		contentParts := strings.Split(content, ".")
		if len(contentParts) >= 2 {
			fileNameParts = contentParts
			break
		}

		if i == (len(dirContents) - 1) {
			return "", fmt.Errorf("Invalid file name provided.")
		}
	}

	switch fileNameParts[len(fileNameParts)-1] {
	case "ts":
		return "typescript", nil
	case "go":
		return "go", nil
	case "cs", "csproj":
		return "dotnet", nil
	default:
		return "", fmt.Errorf("Unsupported language file detected.")
	}
}

// createLambdaFunction creates a Lambda function
func createLambdaFunction(ctx *pulumi.Context, role *iam.Role, logPolicy *iam.RolePolicy, route APIRoute, method string) (*lambda.Function, error) {
	// Determine the language of the function.
	lambdaFilePath := path.Join(route.PathToFiles, method)
	lambdaLanguage, err := detectLambdaLanguage(lambdaFilePath)
	if err != nil {
		return nil, err
	}

	// Determine the Lambda runtime to use.
	var lambdaRuntime string
	var handlerName string
	var handlerZipFile string
	switch lambdaLanguage {
	case "typescript":
		lambdaRuntime = "nodejs12.x"
		handlerName = fmt.Sprintf("%s-%s-handler.%sHandler", route.Name, method, method)
		handlerZipFile, err = PackageTypeScriptLambda(tmpDirName, route.Name, method)
		if err != nil {
			return nil, err
		}
		break
	case "go":
		lambdaRuntime = "go1.x"
		handlerName = fmt.Sprintf("%s-%s-handler", route.Name, method)
		handlerZipFile, err = PackageGoLambda(tmpDirName, route.Name, method)
		if err != nil {
			return nil, err
		}
	case "dotnet":
		lambdaRuntime = "dotnetcore3.1"
		handlerName = fmt.Sprintf("app::app.Functions::%s", utils.DashCaseToSentenceCase(method))
		handlerZipFile, err = PackageDotNetLambda(tmpDirName, route.Name, method)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unsupported runtime detected.")
	}

	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	handlerFileName := path.Join(currentWorkingDirectory, handlerZipFile)
	args := &lambda.FunctionArgs{
		Handler: pulumi.String(handlerName),
		Role:    role.Arn,
		Runtime: pulumi.String(lambdaRuntime),
		Code:    pulumi.NewFileArchive(handlerFileName),
	}

	// Create the lambda using the args.
	function, err := lambda.NewFunction(
		ctx,
		fmt.Sprintf("%s-%s-lambda-function", route.Name, method),
		args,
		pulumi.DependsOn([]pulumi.Resource{logPolicy}),
	)
	if err != nil {
		return nil, fmt.Errorf("Error updating lambda for file [%s]: %v", handlerFileName, err)
	}

	return function, nil
}

// CreateRouteHandler creates a Lambda function used for handling API Gateway requests.
func CreateRouteHandler(ctx *pulumi.Context, route APIRoute, method string) (*lambda.Function, error) {
	lambdaName := fmt.Sprintf("%s-%s", route.Name, method)

	// Create the role.
	role, err := createIAMLambdaRole(ctx, lambdaName)
	if err != nil {
		return nil, err
	}

	// Create the log policy.
	logPolicy, err := createLambdaLogPolicy(ctx, role, lambdaName)
	if err != nil {
		return nil, err
	}

	// Create the function.
	function, err := createLambdaFunction(ctx, role, logPolicy, route, method)
	if err != nil {
		return nil, err
	}

	// Return the function.
	return function, nil
}
