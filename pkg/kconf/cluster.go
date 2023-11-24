package kconf

import (
	"fmt"

	apiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

// change cluster for given context. If context is empty string, then
// the current context will be used
func (kc *KubeConfig) ChangeCluster(context, cluster string) error {
	ctx := kc.GetContext(context)
	if ctx == nil {
		return fmt.Errorf("no context %s in kubeconfing", context)
	}

	usr := kc.GetCluster(cluster)
	if usr == nil {
		return fmt.Errorf("no cluster %s in kubeconfig", cluster)
	}

	ctx.Context.Cluster = cluster
	return nil
}

func (kc *KubeConfig) GetCluster(name string) *apiv1.NamedCluster {
	for i := range kc.Clusters {
		if kc.Clusters[i].Name == name {
			return &kc.Clusters[i]
		}
	}
	return nil
}

func (kc *KubeConfig) addToClusters(clusters ...apiv1.NamedCluster) {
	kc.Clusters = append(kc.Clusters, clusters...)
}

func (kc *KubeConfig) removeFromClusters(name string) {
	for idx, k := range kc.Clusters {
		if k.Name == name {
			kc.Clusters[idx] = kc.Clusters[len(kc.Clusters)-1] // copy last element to index
			kc.Clusters = kc.Clusters[:len(kc.Clusters)-1]     // truncate slice
			return
		}
	}
}

func (c *KubeConfig) RenameCluster(src, dest string) {
	cluster := c.GetCluster(src)
	if cluster != nil {
		cluster.Name = dest
	}
}
