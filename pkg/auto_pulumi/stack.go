package auto_pulumi

import (
	"context"

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
)

// createPulumiProject creates a new Pulumi project and the appropriate
// stacks for each environment.
//
// TODO: Check to see if a project exists before trying to create it. Otherwise I am guessing some
// pretty funky stuff will happen if the similar project names are specified.
func CreatePulumiProject(ctx context.Context, owner string, projectName string, env string, description string) (string, error) {
	// Loop over the environments are create
	stackName := auto.FullyQualifiedStackName(owner, projectName, env)
	_, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, nil)
	if err != nil {
		return "", err
	}

	return stackName, nil
}
