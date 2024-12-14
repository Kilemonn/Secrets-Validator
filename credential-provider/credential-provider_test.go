package credential_provider

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCredentialProviderIdentifier_FromString_InvalidString(t *testing.T) {
	identifier := CredentialProviderIdentifierFromString("Does not exist")
	require.Equal(t, CredentialProviderIdentifierInvalid, identifier)
}

func TestCredentialProviderIdentifier_IsValid(t *testing.T) {
	values := credentialProviderIdentifierValues()
	for _, val := range values {
		require.Equal(t, val != CredentialProviderIdentifierInvalid, val.IsValid())
	}
}

func TestCredentialProviderIdentifier(t *testing.T) {
	values := credentialProviderIdentifierValues()
	labels := credentialProviderIdentifierStrings()

	require.Equal(t, len(values), len(labels))
	for i := range len(values) {
		fromString := CredentialProviderIdentifierFromString(labels[i])
		require.Equal(t, fromString, values[i])
		require.Equal(t, uint(i), fromString.Index())
		require.Equal(t, fromString.String(), labels[i])
	}
}
