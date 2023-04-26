package kconf

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
	apilatest "k8s.io/client-go/tools/clientcmd/api/latest"
	apiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

type KubeConfig struct {
	*apiv1.Config
}

func New() *KubeConfig {
	return &KubeConfig{
		Config: &apiv1.Config{
			Kind:       "Config",
			APIVersion: "v1",
			AuthInfos:  make([]apiv1.NamedAuthInfo, 0),
			Clusters:   make([]apiv1.NamedCluster, 0),
			Contexts:   make([]apiv1.NamedContext, 0),
			Preferences: apiv1.Preferences{
				Extensions: make([]apiv1.NamedExtension, 0),
			},
		},
	}
}

// load data from given file and parse them into
// KubeConfig instance. If file is not set (empty string),
// the function will load and parse first kubeconfig in
// KUBECONFIG env. variable.
func Open(file string) (*KubeConfig, error) {
	if file == "" {
		envValue := os.Getenv("KUBECONFIG")
		configs := strings.Split(envValue, ":")
		file = configs[0]
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return OpenData(data)
}

// OpenBase64 will decode input data from base64 and parse it
func OpenBase64(b64Data []byte) (*KubeConfig, error) {
	// I need to deal with data as string because sometimes
	// from stdin we get weird bytes on end of buffer and YAML
	// cannot be parsed. I'm converting it to string and decode
	// as string. This ensure the decoded YAML can be parsed.
	//
	// ... I'm sure there is a better way how to deal with it
	encodedData, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		return nil, err
	}
	return OpenData([]byte(encodedData))
}

// parse given data as YAML into new KubeConfig
func OpenData(data []byte) (*KubeConfig, error) {
	cfg := &apiv1.Config{}
	_, _, err := apilatest.Codec.Decode(data, &schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Config",
	}, cfg)

	if err != nil {
		return nil, err
	} else {
		return &KubeConfig{Config: cfg}, nil
	}
}

func (c *KubeConfig) Save(filename string) error {
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	return c.WriteTo(fd)
}

func (c *KubeConfig) WriteTo(w io.Writer) error {
	err := apilatest.Codec.Encode(c.Config, w)
	if err != nil {
		return err
	} else {
		return nil
	}
}

// Import all users, contexts and clusters from src kubeconfig
// to current kubeconfig
func (c *KubeConfig) Import(src *KubeConfig) {
	c.addToContexts(src.Contexts...)
	c.addToClusters(src.Clusters...)
	c.addToUsers(src.AuthInfos...)
}

// Export returns you new KubeConfig where is given context
// with required User and Cluster.
func (cfg *KubeConfig) Export(contextName string) (*KubeConfig, error) {
	exported := New()
	exported.CurrentContext = contextName
	ctx, cluster, user := cfg.getFullContext(contextName)
	if ctx == nil {
		return exported, fmt.Errorf("the '%s' is missing in kubeconfig", contextName)
	}

	exported.addToContexts(*ctx)
	exported.addToClusters(*cluster)
	exported.addToUsers(*user)

	return exported, nil
}

// completely remove context by name and context's
// cluster and user
func (c *KubeConfig) Remove(contextName string) {
	ctx := c.getContext(contextName)
	if ctx == nil {
		return
	}

	c.removeFromClusters(ctx.Context.Cluster)
	c.removeFromUsers(ctx.Context.AuthInfo)
	c.removeFromContexts(contextName)
}

// rename context and context's cluster and user. Both
// cluster and user will have name same as is new name
// of context
func (c *KubeConfig) Rename(src, dest string) {
	ctx := c.getContext(src)
	if ctx == nil {
		return
	}

	c.renameCluster(ctx.Context.Cluster, dest)
	c.renameUser(ctx.Context.AuthInfo, dest)
	c.renameContext(src, dest)
}

// split one kubeconfig into smaller kubeconfig pieces for
// each context
func (c *KubeConfig) Split() []*KubeConfig {
	result := make([]*KubeConfig, len(c.Contexts))
	for i := range c.Contexts {
		result[i] = New()
		result[i].addToContexts(c.Contexts[i])
		result[i].CurrentContext = c.Contexts[i].Name

		usr := c.getUser(c.Contexts[i].Context.AuthInfo)
		result[i].addToUsers(*usr)

		cluster := c.getCluster(c.Contexts[i].Context.Cluster)
		result[i].addToClusters(*cluster)
	}
	return result
}
