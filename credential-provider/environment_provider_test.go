package credential_provider

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithLoadedEnvironmentVars(t *testing.T) {
	provider := NewEnvironmentProvider()
	initialLen := len(provider.GetCredentials())

	propertyName := "TestWithLoadedEnvironmentVars"
	propertyValue := "A Value!"
	_, err := provider.GetCredentialWithName(propertyName)
	assert.Error(t, err)

	withLoadedEnvironmentVars(t, map[string]string{propertyName: propertyValue}, func() {
		newProvider := NewEnvironmentProvider()
		assert.Greater(t, len(newProvider.GetCredentials()), initialLen)

		cred, err := newProvider.GetCredentialWithName(propertyName)
		assert.NoError(t, err)
		assert.Equal(t, propertyValue, cred)
	})

	provider = NewEnvironmentProvider()
	_, err = provider.GetCredentialWithName(propertyName)
	assert.Error(t, err)
	assert.Equal(t, initialLen, len(provider.GetCredentials()))
}

func withLoadedEnvironmentVars(t *testing.T, vars map[string]string, testFunc func()) {
	count := 0
	for k, v := range vars {
		if _, exists := os.LookupEnv(k); !exists {
			if os.Setenv(k, v) == nil {
				count += 1
			}
		} else {
			vars[k] = ""
		}
	}

	testFunc()

	for k := range vars {
		if os.Unsetenv(k) == nil {
			count -= 1
		}
	}

	assert.Equal(t, 0, count)
}
