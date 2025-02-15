package config

import (
	"fmt"
	"os"

	"github.com/Kilemonn/Secrets-Validator/constraint"
	credential_provider "github.com/Kilemonn/Secrets-Validator/credential-provider"
	"github.com/Kilemonn/Secrets-Validator/util"
	"gopkg.in/yaml.v3"
)

const (
	yamlPropertyConstraints         = "constraints"
	yamlPropertyCredentialProviders = "credential-providers"

	yamlPropertyCondition = "condition"
	yamlPropertyPattern   = "pattern"
)

func ValidateConfiguration(configFilePath string) (providers []credential_provider.CredentialProvider, constraints []constraint.Constraint, err error) {
	// TODO: Introduce strict/relaxed flag to fail during setup if constraints or providers cannot be resolved

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("Failed to read in file [%s]. Error: [%s].\n", configFilePath, err.Error())
		return
	}

	var unmarshalled = make(map[string]interface{})
	err = yaml.Unmarshal(data, &unmarshalled)
	if err != nil {
		fmt.Printf("Failed to unmarshal data from file [%s]. Error: [%s].\n", configFilePath, err.Error())
		return
	}

	if credProviders, exists := unmarshalled[yamlPropertyCredentialProviders]; exists {
		providers, err = getProviders(credProviders.([]interface{}))
		if err != nil {
			fmt.Printf("Failed to validate [%s]. Error: [%s].\n", yamlPropertyCredentialProviders, err.Error())
			return
		}
	} else {
		fmt.Printf("Failed to validate configuration yaml [%s] does not exist.\n", yamlPropertyCredentialProviders)
		return
	}

	if constraintsMap, exists := unmarshalled[yamlPropertyConstraints]; exists {
		constraints, err = getConstraints(constraintsMap.([]interface{}))
		if err != nil {
			fmt.Printf("Failed to validate [%s]. Error: [%s].\n", yamlPropertyConstraints, err.Error())
			return
		}
	} else {
		fmt.Printf("Failed to validate configuration yaml [%s] does not exist.\n", yamlPropertyConstraints)
		return
	}

	return
}

func getProviders(credProviders []interface{}) (credentialProviders []credential_provider.CredentialProvider, err error) {
	for _, val := range credProviders {
		var provider credential_provider.CredentialProvider
		provider, err = validateProvider(val.(map[string]interface{}))
		if err != nil {
			fmt.Printf("Failed to parse provider with error: [%s].\n", err.Error())
			return
		}

		credentialProviders = append(credentialProviders, provider)
		fmt.Printf("Registered credential provider with ID [%s].\n", provider.Identifier.String())
	}
	return
}

func validateProvider(providerMap map[string]interface{}) (provider credential_provider.CredentialProvider, err error) {
	if len(providerMap) != 1 {
		fmt.Println("Expected only len of 1 holding the provider name.")
		err = fmt.Errorf("found more than one provider name")
		return
	}

	for name, properties := range providerMap {
		providerId := credential_provider.CredentialProviderIdentifierFromString(name)
		if !providerId.IsValid() {
			fmt.Printf("Failed to register credential provider with name [%s], pleaese provide only valid provider names.\n", name)
			err = fmt.Errorf("invalid provider name [%s] provided", name)
			return
		}

		var propertiesAsMap map[string]interface{}
		if properties != nil {
			propertiesAsMap = properties.(map[string]interface{})
		} else {
			propertiesAsMap = make(map[string]interface{})
		}

		provider, err = credential_provider.NewCredentialProvider(providerId, propertiesAsMap)
		if err != nil {
			fmt.Printf("Failed to construct provider because invalid properties were provided. Error: [%s].", err.Error())
			return
		}
	}

	return
}

func getConstraints(constraintsMap []interface{}) (constraints []constraint.Constraint, err error) {
	for _, val := range constraintsMap {
		constraintObj, err := validateConstraint(val.(map[string]interface{}))
		if err != nil {
			fmt.Printf("Failed to validate constraint: [%s].\n", err.Error())
		} else {
			constraints = append(constraints, constraintObj)
			fmt.Printf("Registered constraint with name [%s].\n", constraintObj.Name)
		}
	}
	return
}

func validateConstraint(constraintMap map[string]interface{}) (constraintObj constraint.Constraint, err error) {
	if len(constraintMap) != 1 {
		fmt.Println("Expected only len of 1 holding the constraints name.")
		err = fmt.Errorf("found more than one contraint name")
		return
	}

	for name, properties := range constraintMap {
		notContainedKeys := util.ContainsAllKeys([]string{yamlPropertyCondition, yamlPropertyPattern}, properties.(map[string]interface{}))
		if len(notContainedKeys) != 0 {
			fmt.Printf("Required properties [%s] are not defined for constraint with name [%s].", notContainedKeys, name)
			err = fmt.Errorf("required properties [%s] are not defined for constraint with name [%s]", notContainedKeys, name)
			return
		}

		propertiesAsMap := properties.(map[string]interface{})
		constraintObj, err = constraint.NewConstraint(name, propertiesAsMap[yamlPropertyPattern].(string), propertiesAsMap[yamlPropertyCondition].(string))
		if err != nil {
			fmt.Printf("Failed to construct constraint [%s] due to [%s].\n", name, err.Error())
		}
	}

	return
}
