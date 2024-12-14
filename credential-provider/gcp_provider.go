package credential_provider

import (
	"context"
	"fmt"
	"hash/crc32"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/Kilemonn/Secrets-Validator/util"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	property_project_id           string = "project-id"
	property_credential_file_path string = "credential-file-path"
)

type GcpProvider struct {
	client             *secretmanager.Client
	projectId          string
	ctx                context.Context
	credentialFilePath string
}

func NewGcpProvider(properties map[string]interface{}) (provider GcpProvider, err error) {
	requiredProperties := []string{property_credential_file_path, property_project_id}
	notContained := util.ContainsAllKeys(requiredProperties, properties)
	if len(notContained) > 0 {
		err = fmt.Errorf("missing properties %s, for GCP provider", notContained)
		return
	}

	provider.projectId = properties[property_project_id].(string)
	provider.credentialFilePath = properties[property_credential_file_path].(string)

	provider.ctx = context.Background()
	var client *secretmanager.Client
	client, err = secretmanager.NewClient(provider.ctx, option.WithCredentialsFile(provider.credentialFilePath))
	if err != nil {
		return
	}
	provider.client = client
	return
}

func (p GcpProvider) GetCredentialNames() (names []string, err error) {
	listRequest := &secretmanagerpb.ListSecretsRequest{
		Parent: "projects/" + p.projectId,
	}

	it := p.client.ListSecrets(p.ctx, listRequest)
	for {
		var resp *secretmanagerpb.Secret
		resp, err = it.Next()
		if err == iterator.Done {
			// Iteration completing is not an error that needs to be returned to the caller
			return names, nil
		}

		if err != nil {
			fmt.Printf("Failed to retrieve credential names, error: [%s].", err.Error())
			return
		}
		names = append(names, resp.Name)
	}
}

func (p GcpProvider) GetCredentialWithName(key string) (string, error) {
	credentialRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: key + "/versions/latest",
	}
	result, err := p.client.AccessSecretVersion(p.ctx, credentialRequest)
	if err != nil {
		fmt.Printf("Failed to retrieve credential with name [%s] with error: [%s].\n", key, err.Error())
		return "", err
	}
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		return "", fmt.Errorf("check sum verification failed on property with name [%s]", key)
	}

	return string(result.Payload.Data), nil
}

func (p GcpProvider) Close() (err error) {
	if err = p.client.Close(); err != nil {
		fmt.Printf("Failed to close GCP provider with error: [%s].\n", err.Error())
	}
	return err
}
