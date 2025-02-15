package credential_provider

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Kilemonn/Secrets-Validator/util"
	"golang.org/x/exp/maps"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	property_namespace   = "namespace"
	property_secret_name = "secret-name"
)

type KubernetesProvider struct {
	ctx        context.Context
	client     kubernetes.Interface
	namespace  string
	secretName string
}

func NewKubernetesProvider(properties map[string]interface{}) (provider KubernetesProvider, err error) {
	requiredProperties := []string{property_namespace, property_secret_name}
	notContained := util.ContainsAllKeys(requiredProperties, properties)
	if len(notContained) > 0 {
		err = fmt.Errorf("missing properties %s, for Kubernetes provider", notContained)
		return
	}
	provider.ctx = context.Background()
	provider.namespace = properties[property_namespace].(string)
	provider.secretName = properties[property_secret_name].(string)

	if home := homedir.HomeDir(); home != "" {
		kubeConfigPath := filepath.Join(home, ".kube", "config")
		var config *rest.Config
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return
		}
		provider.client, err = kubernetes.NewForConfig(config)
	} else {
		err = fmt.Errorf("unable to determine home directory for the current user to find the path to /.kube/config")
	}

	return
}

func (p KubernetesProvider) GetCredentialNames() ([]string, error) {
	client := p.client.CoreV1().Secrets(p.namespace)
	secret, err := client.Get(p.ctx, p.secretName, v1.GetOptions{})
	if err != nil {
		return []string{}, err
	}
	return maps.Keys(secret.Data), nil
}

func (p KubernetesProvider) GetCredentialWithName(key string) (string, error) {
	client := p.client.CoreV1().Secrets(p.namespace)
	secret, err := client.Get(p.ctx, p.secretName, v1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(secret.Data[key]), nil
}

func (p KubernetesProvider) Close() error {
	// No-op
	return nil
}
