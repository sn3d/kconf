package kconf

import apiv1 "k8s.io/client-go/tools/clientcmd/api/v1"

func (c *KubeConfig) getContext(name string) *apiv1.NamedContext {
	for i := range c.Contexts {
		if c.Contexts[i].Name == name {
			return &c.Contexts[i]
		}
	}
	return nil
}

func (c *KubeConfig) addToContexts(contexts ...apiv1.NamedContext) {
	c.Contexts = append(c.Contexts, contexts...)
}

func (c *KubeConfig) getFullContext(name string) (*apiv1.NamedContext, *apiv1.NamedCluster, *apiv1.NamedAuthInfo) {
	var (
		ctx     *apiv1.NamedContext
		cluster *apiv1.NamedCluster
		user    *apiv1.NamedAuthInfo
	)

	ctx = c.getContext(name)
	if ctx != nil {
		cluster = c.getCluster(ctx.Context.AuthInfo)
		user = c.getUser(ctx.Context.Cluster)
	}

	return ctx, cluster, user
}

func (c *KubeConfig) removeFromContexts(name string) {
	for idx, k := range c.Contexts {
		if k.Name == name {
			c.Contexts[idx] = c.Contexts[len(c.Contexts)-1] // copy last element to index
			c.Contexts = c.Contexts[:len(c.Contexts)-1]     // truncate slice
			return
		}
	}
}

func (c *KubeConfig) renameContext(src, dest string) {
	context := c.getContext(src)
	if context != nil {
		context.Name = dest
	}
}
