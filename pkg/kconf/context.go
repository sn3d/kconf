package kconf

import apiv1 "k8s.io/client-go/tools/clientcmd/api/v1"

// get context by name. If name is empty string, then it returns you
// current context
func (kc *KubeConfig) GetContext(name string) *apiv1.NamedContext {
	if name == "" {
		name = kc.CurrentContext
	}

	for i := range kc.Contexts {
		if kc.Contexts[i].Name == name {
			return &kc.Contexts[i]
		}
	}
	return nil
}

func (kc *KubeConfig) GetCurrentContext() *apiv1.NamedContext {
	return kc.GetContext(kc.CurrentContext)
}

// rename context and context's cluster and user. Both
// cluster and user will have name same as is new name
// of context
func (kc *KubeConfig) Rename(src, dest string) {
	ctx := kc.GetContext(src)
	if ctx == nil {
		return
	}

	kc.RenameCluster(ctx.Context.Cluster, dest)
	ctx.Context.Cluster = dest

	kc.RenameUser(ctx.Context.AuthInfo, dest)
	ctx.Context.AuthInfo = dest

	kc.RenameContext(src, dest)
}

// completely remove context by name and context's
// cluster and user
func (kc *KubeConfig) Remove(contextName string) {
	ctx := kc.GetContext(contextName)
	if ctx == nil {
		return
	}

	kc.removeFromClusters(ctx.Context.Cluster)
	kc.removeFromUsers(ctx.Context.AuthInfo)
	kc.removeFromContexts(contextName)
}

func (kc *KubeConfig) addToContexts(contexts ...apiv1.NamedContext) {
	kc.Contexts = append(kc.Contexts, contexts...)
}

func (kc *KubeConfig) getFullContext(name string) (*apiv1.NamedContext, *apiv1.NamedCluster, *apiv1.NamedAuthInfo) {
	var (
		ctx     *apiv1.NamedContext
		cluster *apiv1.NamedCluster
		user    *apiv1.NamedAuthInfo
	)

	ctx = kc.GetContext(name)
	if ctx != nil {
		cluster = kc.GetCluster(ctx.Context.Cluster)
		user = kc.GetUser(ctx.Context.AuthInfo)
	}

	return ctx, cluster, user
}

func (kc *KubeConfig) removeFromContexts(name string) {
	for idx, k := range kc.Contexts {
		if k.Name == name {
			kc.Contexts[idx] = kc.Contexts[len(kc.Contexts)-1] // copy last element to index
			kc.Contexts = kc.Contexts[:len(kc.Contexts)-1]     // truncate slice
			return
		}
	}
}

func (kc *KubeConfig) RenameContext(src, dest string) {
	context := kc.GetContext(src)
	if context != nil {
		context.Name = dest
	}
}
