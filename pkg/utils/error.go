package utils

import (
	"log"
)

// HandleError outputs error messages and exits the program.
func HandleError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
		return
	}

	log.Fatal(message)
}
