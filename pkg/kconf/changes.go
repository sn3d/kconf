package kconf

import (
	"fmt"
)

// change default namespace for given context. If given context
// is empty string, then the current context of kubeconfig will
// be used
func (c *KubeConfig) Chns(context, namespace string) error {
	ctx := c.GetContext(context)
	if ctx == nil {
		return fmt.Errorf("no context %s in kubeconfig", context)
	}

	ctx.Context.Namespace = namespace
	return nil
}

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

// change cluster for given context. If context is empty string, then
// the current context will be used
func (c *KubeConfig) Chclus(context, cluster string) error {
	ctx := c.GetContext(context)
	if ctx == nil {
		return fmt.Errorf("no context %s in kubeconfing", context)
	}

	usr := c.GetCluster(cluster)
	if usr == nil {
		return fmt.Errorf("no cluster %s in kubeconfig", cluster)
	}

	ctx.Context.Cluster = cluster
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

	return nil
}
