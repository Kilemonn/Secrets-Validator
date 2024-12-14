package validator

import (
	"fmt"

	"github.com/Kilemonn/Secrets-Validator/constraint"
	credential_provider "github.com/Kilemonn/Secrets-Validator/credential-provider"
)

func ExecuteConstraintsAgainstProviders(providers []credential_provider.CredentialProvider, constraints []constraint.Constraint, debugLog bool) map[string][]string {
	failed := make(map[string][]string)

	for _, provider := range providers {
		credentialNames, err := provider.Provider.GetCredentialNames()
		if err != nil {
			fmt.Printf("Failed to get credential names from provider [%s]. With error: [%s].\n", provider.Identifier.String(), err.Error())
			return failed
		}
		for _, credentialName := range credentialNames {
			credential, err := provider.Provider.GetCredentialWithName(credentialName)
			if err != nil {
				fmt.Printf("Failed to retrieve credential with name [%s] from provider [%s] with error [%s].", credentialName, provider.Identifier.String(), err.Error())
				return failed
			}
			for _, constraint := range constraints {
				if constraint.Pattern.Matches(credentialName) {
					if debugLog {
						fmt.Printf("Credential [%s] matched pattern for constraint [%s], applying condition...\n", credentialName, constraint.Name)
					}
					if !constraint.Condition.ApplyCondition(constraint.Name, credential) {
						if debugLog {
							fmt.Printf("Fail - Provider [%s], Constraint [%s], Credential [%s].\n", provider.Identifier.String(), constraint.Name, credentialName)
						}

						if _, exists := failed[constraint.Name]; !exists {
							failed[constraint.Name] = make([]string, 0)
						}
						failed[constraint.Name] = append(failed[constraint.Name], credentialName)
					} else {
						if debugLog {
							fmt.Printf("Pass - Provider [%s], Constraint [%s], Credential [%s].\n", provider.Identifier.String(), constraint.Name, credentialName)
						}
					}
				} else {
					if debugLog {
						fmt.Printf("Credential with name [%s] did not match on pattern for constraint [%s].\n", credentialName, constraint.Name)
					}
				}
			}
		}
	}
	return failed
}
