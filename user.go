package kconf

// Represents element in 'users' list. Each element have name and
// nested 'user' structure
type UserElm struct {
	Name string `yaml:"name",omitempty`
	User struct {
		ClientKey             string `yaml:"client-key",omitempty`
		ClientKeyData         string `yaml:"client-key-data",omitempty`
		ClientCertificate     string `yaml:"client-certificate",omitempty`
		ClientCertificateData string `yaml:"client-certificate-data",omitempty`
	} `yaml:"user",omitempty`
}
