package credential_provider

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type CredentialProviderIdentifier uint

const (
	CredentialProviderIdentifierInvalid    CredentialProviderIdentifier = iota
	CredentialProviderIdentifierGCP        CredentialProviderIdentifier = iota
	CredentialProviderIdentifierAWS        CredentialProviderIdentifier = iota
	CredentialProviderIdentifierENV        CredentialProviderIdentifier = iota
	CredentialProviderIdentifierKubernetes CredentialProviderIdentifier = iota

	provider_identifier_invalid_string    = "invalid"
	provider_identifier_gcp_string        = "gcp"
	provider_identifier_aws_string        = "aws"
	provider_identifier_env_string        = "env"
	provider_identifier_kubernetes_string = "kubernetes"
)

var (
	ErrCredentialNotFound = errors.New("credential with provided name was not found")
)

func credentialProviderIdentifierValues() []CredentialProviderIdentifier {
	return []CredentialProviderIdentifier{
		CredentialProviderIdentifierInvalid,
		CredentialProviderIdentifierGCP,
		CredentialProviderIdentifierAWS,
		CredentialProviderIdentifierENV,
		CredentialProviderIdentifierKubernetes,
	}
}

func credentialProviderIdentifierStrings() []string {
	return []string{
		provider_identifier_invalid_string,
		provider_identifier_gcp_string,
		provider_identifier_aws_string,
		provider_identifier_env_string,
		provider_identifier_kubernetes_string,
	}
}

func CredentialProviderIdentifierFromString(input string) CredentialProviderIdentifier {
	index := slices.Index(credentialProviderIdentifierStrings(), strings.ToLower(input))
	if index == -1 {
		return CredentialProviderIdentifierInvalid
	} else {
		return CredentialProviderIdentifier(index)
	}
}

func (c CredentialProviderIdentifier) String() string {
	return credentialProviderIdentifierStrings()[c.Index()]
}

func (c CredentialProviderIdentifier) Index() uint {
	return uint(c)
}

func (c CredentialProviderIdentifier) IsValid() bool {
	return c.Index() != CredentialProviderIdentifierInvalid.Index()
}

func IsValidProvider(providerName string) bool {
	return CredentialProviderIdentifierFromString(providerName).IsValid()
}

type CredentialProvider struct {
	Identifier CredentialProviderIdentifier
	Provider   Provider
}

func NewCredentialProvider(id CredentialProviderIdentifier, properties map[string]interface{}) (provider CredentialProvider, err error) {
	provider.Identifier = id

	if id == CredentialProviderIdentifierENV {
		provider.Provider = NewEnvironmentProvider()
	} else if id == CredentialProviderIdentifierGCP {
		provider.Provider, err = NewGcpProvider(properties)
	} else if id == CredentialProviderIdentifierAWS {
		provider.Provider, err = NewAwsProvider(properties)
	} else if id == CredentialProviderIdentifierKubernetes {
		provider.Provider, err = NewKubernetesProvider(properties)
	} else {
		err = fmt.Errorf("invalid provider identifier found [%s]", id.String())
	}

	return
}

// A credential provider's provider interface this determines how credentials are retrieved
type Provider interface {
	GetCredentialNames() ([]string, error)
	GetCredentialWithName(string) (string, error)
	Close() error
}
