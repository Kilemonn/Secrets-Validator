package credential_provider

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func withSecretConfig(t *testing.T, filepath string, testFunc func()) {
	cmd := exec.Command("kubectl", "apply", "-f", filepath)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to apply kubernetes config file with path: [%s]. With output [%s]. And error: [%s].", filepath, output, err.Error())
		t.FailNow()
		return
	}
	// Assuming the configuration was applied correctly with the provided file, I am assuming it will be delete successfully too. For now
	defer exec.Command("kubectl", "delete", "-f", filepath)
	testFunc()
}

func TestKubernetesProvider(t *testing.T) {
	m := make(map[string]interface{})
	m[property_namespace] = "default"
	m[property_secret_name] = "my-kubernetes-secret"

	withSecretConfig(t, "./kubernetes-test-secret.yaml", func() {
		provider, err := NewKubernetesProvider(m)
		require.NoError(t, err)
		defer provider.Close()

		val, err := provider.GetCredentialWithName("credential-name")
		require.NoError(t, err)
		require.NotEmpty(t, val)
	})
}
