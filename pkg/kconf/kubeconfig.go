package kconf

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
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

// open the given file and returns KubeConfig instance.
// If given file is empty string, the function will load kubeconfig
// from KUBECONFIG env. variable. If KUBECONFIG variable is not set,
// then the function open default ~/.kube/config file.
func Open(file string) (kc *KubeConfig, path string, err error) {
	// figureout what file path need to be open
	path = file
	if path == "" {
		path = os.Getenv("KUBECONFIG")
		if path == "" {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, ".kube", "config")
		}
	}

	multipleConfigs := strings.Contains(path, ":")
	if multipleConfigs {
		configs := strings.Split(path, ":")
		path = configs[0]
	}

	// read the file and parse KubeConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, path, &OpenError{original: err, path: path}
	}

	kc, err = OpenData(data)
	if err != nil {
		return nil, path, &OpenError{original: err, path: path}
	}

	return kc, path, nil
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
	}

	// sort contexts
	sort.SliceStable(cfg.Contexts, func(i, j int) bool {
		return strings.Compare(cfg.Contexts[i].Name, cfg.Contexts[j].Name) <= 0
	})

	return &KubeConfig{Config: cfg}, nil
}

// save the KubeConfig into file as YAML
func (c *KubeConfig) Save(filename string) error {
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	return c.WriteTo(fd)
}

// write the KubeConfig into writer as YAML. It's used when
// you need print the KubeConfig YAML into standard output etc.
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
func (c *KubeConfig) Import(src *KubeConfig, opts *ImportOptions) {

	// If 'As' option is set and there is only one context, then it's
	// renamed to 'As'
	// Right now it's hard to do some meaningful and simple renaming
	// if src have more than 1 context. Maybe later will be implemented
	// some logic based on some need.
	if opts.As != "" && len(src.Contexts) == 1 {
		src.Rename(src.Contexts[0].Name, opts.As)
	}

	c.addToContexts(src.Contexts...)
	c.addToClusters(src.Clusters...)
	c.addToUsers(src.AuthInfos...)
}

// Export returns you new KubeConfig where is given context
// with required User and Cluster.
// If contextName is empty string, then the current context will be
// exported
func (cfg *KubeConfig) Export(contextName string, opts *ExportOptions) (*KubeConfig, error) {
	if contextName == "" {
		contextName = cfg.CurrentContext
	}

	exported := New()
	exported.CurrentContext = contextName
	ctx, cluster, user := cfg.getFullContext(contextName)
	if ctx == nil {
		return exported, fmt.Errorf("the '%s' is missing in kubeconfig", contextName)
	}

	exported.addToContexts(*ctx)

	if cluster != nil {
		exported.addToClusters(*cluster)
	}

	if cluster != nil {
		exported.addToUsers(*user)
	}

	// If 'As' option is available, the context will be exported
	// with given name
	if opts.As != "" {
		exported.Rename(contextName, opts.As)
	}

	// set current context to exported context
	exported.CurrentContext = exported.Contexts[0].Name

	return exported, nil
}

// completely remove context by name and context's
// cluster and user
func (c *KubeConfig) Remove(contextName string) {
	ctx := c.GetContext(contextName)
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
	ctx := c.GetContext(src)
	if ctx == nil {
		return
	}

	c.renameCluster(ctx.Context.Cluster, dest)
	ctx.Context.Cluster = dest

	c.renameUser(ctx.Context.AuthInfo, dest)
	ctx.Context.AuthInfo = dest

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

		usr := c.GetUser(c.Contexts[i].Context.AuthInfo)
		result[i].addToUsers(*usr)

		cluster := c.GetCluster(c.Contexts[i].Context.Cluster)
		result[i].addToClusters(*cluster)
	}
	return result
}
