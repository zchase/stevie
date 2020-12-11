package utils

import (
	"github.com/manifoldco/promptui"
)

// PromptRequiredString prompst the user for required string.
func PromptRequiredString(label string) (string, error) {
	// Create the prompt.
	prompt := promptui.Prompt{
		Label: label,
	}

	// Run the prompt.
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	// Return the result.
	return result, nil
}

// PromptSecretString prompts the user for a string that is a secret.
func PromptSecretString(label string) (string, error) {
	// Create the prompt.
	prompt := promptui.Prompt{
		Label:       label,
		HideEntered: true,
	}

	// Run the prompt.
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	// Return the result.
	return result, nil
}
