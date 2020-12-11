package utils

import (
	"strings"

	"github.com/fatih/camelcase"
)

func convertStringsToLowerCase(parts []string) []string {
	var lowercaseParts []string
	for _, part := range parts {
		lowercaseParts = append(lowercaseParts, strings.ToLower(part))
	}
	return lowercaseParts
}

// CamelCaseToDashCase converts a camelCase name to dash case.
func CamelCaseToDashCase(name string) string {
	parts := camelcase.Split(name)
	lowercaseParts := convertStringsToLowerCase(parts)
	return strings.Join(lowercaseParts, "-")
}

// SentenceToCamelCase turns a sentence into camelCase.
func SentenceToCamelCase(sentence string) string {
	var finalParts []string
	parts := strings.Split(sentence, " ")

	for i, part := range parts {
		lowercasePart := strings.ToLower(part)

		if i == 0 {
			finalParts = append(finalParts, lowercasePart)
		} else {
			lowercasePartChars := strings.Split(lowercasePart, "")
			lowercasePartChars[0] = strings.ToUpper(lowercasePartChars[0])
			finalParts = append(finalParts, strings.Join(lowercasePartChars, ""))
		}
	}

	return strings.Join(finalParts, "")
}

// SentenceToDashCase takes a sentence and turns it into
// a lowercase string with dash separators.
func SentenceToDashCase(name string) string {
	// Split the string into an array by spaces.
	parts := strings.Split(name, " ")
	lowercaseParts := convertStringsToLowerCase(parts)
	return strings.Join(lowercaseParts, "-")
}
