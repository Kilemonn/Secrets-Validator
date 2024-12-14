package credential_provider

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/maps"
)

type EnvironmentProvider struct {
	credentials map[string]string
}

func NewEnvironmentProvider() (provider EnvironmentProvider) {
	provider = EnvironmentProvider{
		credentials: make(map[string]string),
	}
	for _, v := range os.Environ() {
		index := strings.Index(v, "=")
		if index != -1 {
			provider.credentials[v[0:index]] = v[index+1:]
		} else {
			fmt.Printf("No '=' found in environment value: [%s]", v)
		}
	}
	return
}

func (p EnvironmentProvider) GetCredentials() map[string]string {
	return p.credentials
}

func (p EnvironmentProvider) GetCredentialNames() ([]string, error) {
	return maps.Keys(p.credentials), nil
}

func (p EnvironmentProvider) GetCredentialWithName(key string) (string, error) {
	val, exists := p.credentials[key]
	if !exists {
		return "", ErrCredentialNotFound
	}
	return val, nil
}

func (p EnvironmentProvider) Close() error {
	// No-op
	return nil
}
