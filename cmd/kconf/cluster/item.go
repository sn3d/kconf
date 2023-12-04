package cluster

import apiv1 "k8s.io/client-go/tools/clientcmd/api/v1"

type ClusterItem struct {
	Cluster *apiv1.NamedCluster
}

func (i ClusterItem) Title() string       { return i.Cluster.Name }
func (i ClusterItem) FilterValue() string { return i.Cluster.Name }
