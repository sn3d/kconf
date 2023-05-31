package kconf_test

import (
	"testing"

	"github.com/sn3d/kconf"
	. "github.com/sn3d/tdata"
)

func Test_Chns(t *testing.T) {
	InitTestdata()

	//GIVEN: existing kubeconfig with current context set
	config, err := kconf.Open(Abs("changes-test.yaml"))
	if err != nil {
		t.FailNow()
	}

	//WHEN: we change namespace for empty "" context
	config.Chns("", "changed-ns")

	//THEN: the namespace for current contex should be changed
	ctx := config.GetContext(config.CurrentContext)
	if ctx.Context.Namespace != "changed-ns" {
		t.FailNow()
	}
}

func Test_Chusr(t *testing.T) {
	InitTestdata()

	//GIVEN: existing kubeconfig with 2 users
	config, err := kconf.Open(Abs("changes-test.yaml"))
	if err != nil {
		t.FailNow()
	}

	//WHEN: we change user for current context to another
	config.Chusr("", "John")

	//THEN: the context should have user changed
	ctx := config.GetContext("")
	if ctx.Context.AuthInfo != "John" {
		t.FailNow()
	}
}

func Test_Chclus(t *testing.T) {
	InitTestdata()

	//GIVEN: existing kubeconfig with 2 clusters
	config, err := kconf.Open(Abs("changes-test.yaml"))
	if err != nil {
		t.FailNow()
	}

	//WHEN: we change user for current context to another
	config.Chclus("", "blue-cluster")

	//THEN: the context should have user changed
	ctx := config.GetContext("")
	if ctx.Context.Cluster != "blue-cluster" {
		t.FailNow()
	}
}

func Test_Clustermod(t *testing.T) {
	InitTestdata()

	//GIVEN: existing kubeconfig with cluster
	config, err := kconf.Open(Abs("changes-test.yaml"))
	if err != nil {
		t.FailNow()
	}

	//WHEN: we change cluster's URL
	config.Clustermod("blue-cluster", &kconf.ClustermodOptions{
		ServerURL: "http://changed.com",
	})

	//THEN: the cluster's server should be changed
	cluster := config.Clusters[0]
	if cluster.Cluster.Server != "http://changed.com" {
		t.FailNow()
	}
}

func Test_Usermod(t *testing.T) {
	InitTestdata()

	//GIVEN: existing kubeconfig with cluster
	config, err := kconf.Open(Abs("changes-test.yaml"))
	if err != nil {
		t.FailNow()
	}

	//WHEN: we change John's token
	config.Usermod("John", &kconf.UsermodOptions{
		Token: "token-123",
	})

	//THEN: the context should have user changed
	user := config.AuthInfos[0]
	if user.AuthInfo.Token != "token-123" {
		t.FailNow()
	}
}
