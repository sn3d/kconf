package kconf

// Represents each element in 'clusters' list. Each element have name and
// nested 'cluster' structure
type ClusterElm struct {
	Name    string `yaml:"name"`
	Cluster struct {
		Server                   string `yaml:"server",omitempty`
		CertificateAuthority     string `yaml:"certificate-authority",omitempty`
		CertificateAuthorityData string `yaml:"certificate-authority-data",omitempty`
	} `yaml:"cluster",omitempty`
}
