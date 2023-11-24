package ctx

import (
	"fmt"

	"github.com/sn3d/kconf/pkg/kconf"
	apiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

type ContextItem struct {
	Context *apiv1.NamedContext
	Kconf   *kconf.KubeConfig
}

func (i ContextItem) Title() string { return i.Context.Name }

func (i ContextItem) FilterValue() string { return i.Context.Name }

func (i ContextItem) Description() string {
	ctx := i.Context.Context

	// get the URL of cluster
	url := ""
	if ctx.Cluster != "" {
		cluster := i.Kconf.GetCluster(ctx.Cluster)
		if cluster != nil {
			url = cluster.Cluster.Server
		}
	}

	if ctx.Namespace == "" {
		return fmt.Sprintf("url: %s", url)
	} else {
		return fmt.Sprintf("url: %s namespace: %s", url, ctx.Namespace)
	}
}
