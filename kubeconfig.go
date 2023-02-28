package kconf

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type KubeConfig struct {
	ApiVersion     string       `yaml:"apiVersion"`
	Kind           string       `yaml:"kind"`
	CurrentContext string       `yaml:"current-context"`
	Clusters       []ClusterElm `yaml:"clusters",omitempty`
	Users          []UserElm    `yaml:"users",omitempty`
	Contexts       []ContextElm `yaml:"contexts",omitempty`
}

// NewKubeConfig creates empty config with pre-filled kind and version
func New() *KubeConfig {
	cfg := new(KubeConfig)
	cfg.ApiVersion = "v1"
	cfg.Kind = "Config"
	return cfg
}

// OpenDefault load kubeconfig data from file defined in
// KUBECONFIG env. variable. Because KUBECONFIG might contain
// multiple files, we took the first file
func OpenDefault() (*KubeConfig, error) {
	filename := getDefaultKubeConfig()
	return OpenFile(filename)
}

// OpenKubeConfigFile load data from file and parse them into
// KubeConfig instance.
func OpenFile(file string) (*KubeConfig, error) {
	filename, _ := filepath.Abs(file)
	yamlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return Open(yamlData)
}

// OpenBase64 will decode input data from base64 and parse it
func OpenBase64(b64Data []byte) (*KubeConfig, error) {
	// we need to deal with data as string because sometimes
	// from stdin we get weird bytes on end of buffer and YAML
	// cannot be parsed. I'm converting it to string and decode
	// as string. This ensure the decoded YAML can be parsed.
	//
	// ... I'm sure there is a better way how to deal with it
	encodedData, err := base64.StdEncoding.DecodeString(string(b64Data))
	if err != nil {
		return nil, err
	}
	return Open([]byte(encodedData))
}

// OpenKubeConfig will parse data and returns you KubeConfig instance
func Open(data []byte) (*KubeConfig, error) {
	cfg := New()
	err := yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Import all users, contexts and clusters from src kubeconfig
// to current kubeconfig
func (cfg *KubeConfig) Import(src *KubeConfig) {
	cfg.Users = append(cfg.Users, src.Users...)
	cfg.Contexts = append(cfg.Contexts, src.Contexts...)
	cfg.Clusters = append(cfg.Clusters, src.Clusters...)
}

// Export returns you new KubeConfig where is given context
// with required User and Cluster.
func (cfg *KubeConfig) Export(contextName string) (*KubeConfig, error) {

	exportedKubeConfig := New()
	exportedKubeConfig.CurrentContext = contextName

	for _, ctx := range cfg.Contexts {
		if ctx.Name == contextName {
			exportedKubeConfig.Contexts = []ContextElm{ctx}
			break
		}
	}

	if len(exportedKubeConfig.Contexts) == 0 {
		return nil, fmt.Errorf("context '%s' not found", contextName)
	}

	// get user
	usrName := exportedKubeConfig.Contexts[0].Context.User
	for _, usr := range cfg.Users {
		if usr.Name == usrName {
			exportedKubeConfig.Users = []UserElm{usr}
			break
		}
	}

	// get cluster
	clusterName := exportedKubeConfig.Contexts[0].Context.Cluster
	for _, cluster := range cfg.Clusters {
		if cluster.Name == clusterName {
			exportedKubeConfig.Clusters = []ClusterElm{cluster}
		}
	}

	return exportedKubeConfig, nil
}

// SaveDefault store the given KubeConfig into file defined
// in KUBECONFIG env. variable.
func (cfg *KubeConfig) SaveDefault() error {
	filename := getDefaultKubeConfig()
	return cfg.Save(filename)
}

// Save write KubeConfig's YAML into file. If given file is
// empty string, then it's saved into existing config in KUBECONFIG
func (cfg *KubeConfig) Save(file string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if file == "" {
		file = getDefaultKubeConfig()
	}

	err = ioutil.WriteFile(file, data, 0)
	if err != nil {
		return err
	}

	return nil
}

// ToString render KubeConfig YAML as string
func (cfg *KubeConfig) ToString() string {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err.Error()
	} else {
		return string(data)
	}
}

// getDefaultKubeConfig parse KUBECONFIG and returns you
// path to first kubeconfig. We need to parse it because
// KUBECONFIG env. variable might contain multiple config
// files.
func getDefaultKubeConfig() string {
	envValue := os.Getenv("KUBECONFIG")
	configs := strings.Split(envValue, ":")
	return configs[0]
}

// RemoveContext by his name
func (cfg *KubeConfig) RemoveContext(context string) {
	idx := -1
	var user string
	var cluster string
	for i, ctx := range cfg.Contexts {
		if ctx.Name == context {
			cluster = ctx.Context.Cluster
			user = ctx.Context.User
			idx = i
			break
		}
	}

	// not found
	if idx == -1 {
		return
	}

	cfg.removeFromContexts(idx)

	// find and remove user
	for i, usr := range cfg.Users {
		if usr.Name == user {
			cfg.removeFromUsers(i)
			break
		}
	}

	// find and remove cluster
	for i, clst := range cfg.Clusters {
		if clst.Name == cluster {
			cfg.removeFromClusters(i)
			break
		}
	}
}

func (cfg *KubeConfig) removeFromContexts(idx int) {
	s := cfg.Contexts
	s[idx] = s[len(s)-1] // Copy last element to index i.
	s = s[:len(s)-1]     // Truncate slice.
	cfg.Contexts = s
}

func (cfg *KubeConfig) removeFromUsers(idx int) {

	s := cfg.Users
	s[idx] = s[len(s)-1] // Copy last element to index i.
	s = s[:len(s)-1]     // Truncate slice.
	cfg.Users = s
}

func (cfg *KubeConfig) removeFromClusters(idx int) {
	s := cfg.Clusters
	s[idx] = s[len(s)-1] // Copy last element to index i.
	s = s[:len(s)-1]     // Truncate slice.
	cfg.Clusters = s
}
