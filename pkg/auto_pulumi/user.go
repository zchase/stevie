package auto_pulumi

import (
	"fmt"
	"os"
	"strings"

	"github.com/zchase/stevie/pkg/utils"
)

var pulumiAccessTokenEnvName = "PULUMI_ACCESS_TOKEN"

// validateAccessToken validates a Pulumi Access Token is set
// and is able to authenticate to a Pulumi account.
func validateAccessToken(token string) error {
	// Check the token is set.
	_, err := os.LookupEnv(pulumiAccessTokenEnvName)
	if !err {
		return fmt.Errorf("Pulumi Access Token is not set")
	}

	// Check that we are able to authenticate with Pulumi.
	_, cErr := utils.RunCommand("pulumi", []string{"whoami"})
	if cErr != nil {
		return fmt.Errorf("Invalid Pulumi Access Token provided")
	}

	return nil
}

// setPulumiAccessTokenEvnVariable sets the Pulumi Access Token
// as an environment variable.
func setPulumiAccessTokenEvnVariable(token string) error {
	// Set the token.
	err := os.Setenv(pulumiAccessTokenEnvName, token)
	if err != nil {
		return fmt.Errorf("Error setting Pulumi Access Token Environment Variable: %v", err)
	}

	// Validate the token is set.
	err = validateAccessToken(token)
	if err != nil {
		err = os.Unsetenv(pulumiAccessTokenEnvName)
		if err != nil {
			return fmt.Errorf("Error unsetting invalid Pulumi Access Token: %v", err)
		}

		return err
	}

	return nil
}

func PromptForPulumiAccessToken() error {
	// Prompt for the token.
	utils.Print("Please enter your Pulumi Access Token.")
	token, err := utils.PromptSecretString("Pulumi Access Token")
	if err != nil {
		return err
	}

	return setPulumiAccessTokenEvnVariable(token)
}

// GetCurrentPulumiUser gets the current Pulumi username.
func GetCurrentPulumiUser() (string, error) {
	output, err := utils.RunCommand("pulumi", []string{"whoami"})
	if err != nil {
		err = PromptForPulumiAccessToken()
		if err != nil {
			return "", fmt.Errorf("Couldn't successfully authenticate the Pulumi users: %v", err)
		}
	}

	return strings.TrimSpace(output), nil
}
