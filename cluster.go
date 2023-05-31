package kconf

import (
	apiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

func (c *KubeConfig) GetCluster(name string) *apiv1.NamedCluster {
	for i := range c.Clusters {
		if c.Clusters[i].Name == name {
			return &c.Clusters[i]
		}
	}
	return nil
}

func (c *KubeConfig) addToClusters(clusters ...apiv1.NamedCluster) {
	c.Clusters = append(c.Clusters, clusters...)
}

func (c *KubeConfig) removeFromClusters(name string) {
	for idx, k := range c.Clusters {
		if k.Name == name {
			c.Clusters[idx] = c.Clusters[len(c.Clusters)-1] // copy last element to index
			c.Clusters = c.Clusters[:len(c.Clusters)-1]     // truncate slice
			return
		}
	}
}

func (c *KubeConfig) renameCluster(src, dest string) {
	cluster := c.GetCluster(src)
	if cluster != nil {
		cluster.Name = dest
	}
}
