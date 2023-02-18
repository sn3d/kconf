package kconf

// Represents each element in 'contexts' list in kubeconfig.
type ContextElm struct {
	Name    string `yaml:"name"`
	Context struct {
		Cluster   string `yaml:"cluster",omitempty`
		Namespace string `yaml:"namespace",omitempty`
		User      string `yaml:"user",omitempty`
	} `yaml:"context",omitempty`
}
