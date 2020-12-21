package utils

import (
	"github.com/manifoldco/promptui"
)

// PromptSelection prompts the user to pick from a list of choices.
func PromptSelection(label string, choices []string) (string, error) {
	// Create the prompt.
	prompt := promptui.Select{
		Label: label,
		Items: choices,
	}

	// Run the prompt.
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	// Return the choice.
	return result, nil
}

// PromptRequiredString prompts the user for required string.
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
