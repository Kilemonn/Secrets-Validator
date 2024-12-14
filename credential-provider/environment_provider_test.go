package credential_provider

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithLoadedEnvironmentVars(t *testing.T) {
	initialProvider := NewEnvironmentProvider()
	defer initialProvider.Close()
	initialLen := len(initialProvider.GetCredentials())

	propertyName := "TestWithLoadedEnvironmentVars"
	propertyValue := "A Value!"
	_, err := initialProvider.GetCredentialWithName(propertyName)
	require.Error(t, err)

	withLoadedEnvironmentVars(t, map[string]string{propertyName: propertyValue}, func() {
		newProvider := NewEnvironmentProvider()
		defer newProvider.Close()
		require.Greater(t, len(newProvider.GetCredentials()), initialLen)

		cred, err := newProvider.GetCredentialWithName(propertyName)
		require.NoError(t, err)
		require.Equal(t, propertyValue, cred)
	})

	provider := NewEnvironmentProvider()
	defer provider.Close()
	_, err = provider.GetCredentialWithName(propertyName)
	require.Error(t, err)
	require.Equal(t, initialLen, len(provider.GetCredentials()))
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

	require.Equal(t, 0, count)
}
