package kconf

import (
	"fmt"
	"io/ioutil"
)

// change user for given context. If context is empty string, then
// the current context will be used
func (c *KubeConfig) Chusr(context, user string) error {
	ctx := c.GetContext(context)
	if ctx == nil {
		return fmt.Errorf("no context %s in kubeconfing", context)
	}

	usr := c.GetUser(user)
	if usr == nil {
		return fmt.Errorf("no user %s in kubeconfig", user)
	}

	ctx.Context.AuthInfo = user
	return nil
}

// changing clusters's parameters
func (c *KubeConfig) Clustermod(clusterName string, opts *ClustermodOptions) error {
	clst := c.GetCluster(clusterName)
	if clst == nil {
		return fmt.Errorf("no cluster %s in kubeconfig", clusterName)
	}

	if opts.ServerURL != "" {
		clst.Cluster.Server = opts.ServerURL
	}

	return nil
}

// changing user's parameters
func (c *KubeConfig) Usermod(userName string, opts *UsermodOptions) error {
	usr := c.GetUser(userName)
	if usr == nil {
		return fmt.Errorf("no user %s in kubeconfig", userName)
	}

	if opts.Token != "" {
		usr.AuthInfo.Token = opts.Token
	}

	if opts.ClientCertificateFile != "" {
		cert, err := ioutil.ReadFile(opts.ClientCertificateFile)
		if err != nil {
			return fmt.Errorf("cannot read client-cert file %s", opts.ClientCertificateFile)
		}
		usr.AuthInfo.ClientCertificateData = cert
	}

	if opts.ClientKeyFile != "" {
		key, err := ioutil.ReadFile(opts.ClientKeyFile)
		if err != nil {
			return fmt.Errorf("cannot read client-key file %s", opts.ClientKeyFile)
		}
		usr.AuthInfo.ClientKeyData = key
	}

	return nil
}
