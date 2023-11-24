package kconf

import (
	"bytes"
	gocontext "context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// change default namespace for given context. If given context
// is empty string, then the current context of kubeconfig will
// be used
func (kc *KubeConfig) ChangeNamespace(context, namespace string) error {
	ctx := kc.GetContext(context)
	if ctx == nil {
		return fmt.Errorf("no context %s in kubeconfig", context)
	}

	ctx.Context.Namespace = namespace
	return nil
}

func (kc *KubeConfig) GetAllNamespaces(context string) ([]string, error) {
	if context == "" {
		context = kc.CurrentContext
	}

	// get k8s client for given kubeconfig and context
	restCfg, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*clientcmdapi.Config, error) {
		apiConfig, err := kc.convertToClientcmdApiConfig()
		if err != nil {
			return nil, err
		}

		apiConfig.CurrentContext = context
		return apiConfig, nil
	})

	if err != nil {
		return []string{}, err
	}

	k8sClient, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return []string{}, err
	}

	// get all namespaces
	namespaces, err := k8sClient.CoreV1().Namespaces().List(gocontext.Background(), metav1.ListOptions{})
	if err != nil {
		return []string{}, err
	}

	result := make([]string, len(namespaces.Items))
	for i, ns := range namespaces.Items {
		result[i] = ns.Name
	}

	return result, nil
}

// This function converts KubeConfig into clientcmdapi.Config. The client-go library
// comes with two Config-s. One is in clientcmd/api and one is in clientcmd/api/v1 package.
// The current Kubeconfig is using clientcmd/api/v1. But clientcmd/api is used for establishing
// k8s client connection. So this function converts KubeConfig into clientcmdapi.Config.
//
// Because there is no better way, we do conversion like: we write the current Kubeconfig into
// memory as YAML, and then we load the YAML into clientcmdapi.Config. This is not very effective
// but it works.
func (c *KubeConfig) convertToClientcmdApiConfig() (*clientcmdapi.Config, error) {
	buf := new(bytes.Buffer)
	c.WriteTo(buf)
	clientcmdApiConfig, err := clientcmd.Load(buf.Bytes())
	if err != nil {
		return nil, err
	}
	return clientcmdApiConfig, nil
}
