package main

import (
	"flag"
	"fmt"

	"github.com/Kilemonn/Secrets-Validator/config"
	"github.com/Kilemonn/Secrets-Validator/consts"
	"github.com/Kilemonn/Secrets-Validator/validator"
)

func main() {
	var configFilePath string
	var debugLog bool
	flag.StringVar(&configFilePath, consts.ARG_FILE_PATH, "", "configuration file path")
	flag.BoolVar(&debugLog, consts.ARG_DEBUG, false, "enable debug logging")
	flag.Parse()

	if configFilePath == "" {
		fmt.Printf("Expected value for flag [%s].\n", consts.ARG_FILE_PATH)
		return
	}

	providers, constraints, err := config.ValidateConfiguration(configFilePath)
	if err != nil {
		fmt.Printf("Failed to initialise constraints and/or providers. Error: [%s].\n", err.Error())
	}

	failedConstraints := validator.ExecuteConstraintsAgainstProviders(providers, constraints, debugLog)
	if len(failedConstraints) > 0 {
		fmt.Printf("Validation failed, the following constraints failed on the following entries:\n")
		for key, val := range failedConstraints {
			fmt.Printf("Constraint name: [%s]:\n", key)
			for _, v := range val {
				fmt.Printf("\tFailed - Credential name: [%s]\n", v)
			}
			fmt.Println("")
		}
	}

	for _, provider := range providers {
		provider.Provider.Close()
	}
}
