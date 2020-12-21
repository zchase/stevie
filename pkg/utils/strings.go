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

// stringPartsToCamelCase turns a array of strings into a camelCase word.
func stringPartsToCamelCase(parts []string) string {
	var finalParts []string
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

// stringPartsToSentenceCase turns an array of strings into a sentence case word.
func stringPartsToSentenceCase(parts []string) string {
	var finalParts []string
	for _, part := range parts {
		lowercasePart := strings.ToLower(part)
		lowercasePartChars := strings.Split(lowercasePart, "")
		lowercasePartChars[0] = strings.ToUpper(lowercasePartChars[0])
		finalParts = append(finalParts, strings.Join(lowercasePartChars, ""))
	}

	return strings.Join(finalParts, "")
}

// DashCaseToCamelCase converts a dash-case name to camelCase.
func DashCaseToCamelCase(name string) string {
	parts := strings.Split(name, "-")
	return stringPartsToCamelCase(parts)
}

// DashCaseToSentenceCase converts a dash-case name to SentenceCase
func DashCaseToSentenceCase(name string) string {
	parts := strings.Split(name, "-")
	return stringPartsToSentenceCase(parts)
}

// SentenceToCamelCase turns a sentence into camelCase.
func SentenceToCamelCase(sentence string) string {
	parts := strings.Split(sentence, " ")
	return stringPartsToCamelCase(parts)
}

// SentenceToDashCase takes a sentence and turns it into
// a lowercase string with dash separators.
func SentenceToDashCase(name string) string {
	// Split the string into an array by spaces.
	parts := strings.Split(name, " ")
	lowercaseParts := convertStringsToLowerCase(parts)
	return strings.Join(lowercaseParts, "-")
}
