package kconf

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/sn3d/tdata"
)

func Test_Open(t *testing.T) {
	InitTestdata()

	cfg, _, err := Open(Abs("open-test.yaml"))
	if err != nil {
		t.FailNow()
	}

	fmt.Println(cfg)
	fmt.Println(cfg.Clusters[0].Cluster.CertificateAuthorityData)

	err = cfg.Save(Abs("save-test.yaml"))
	if err != nil {
		t.FailNow()
	}

	same := CompareFiles("save-test.yaml", "save-expected.yaml")
	if !same {
		t.FailNow()
	}
}

func Test_OpenSorted(t *testing.T) {
	InitTestdata()

	cfg, _, err := Open(Abs("sorted-test.yaml"))
	if err != nil {
		t.FailNow()
	}

	if !strings.HasPrefix(cfg.Contexts[0].Name, "01-") {
		t.FailNow()
	}

	if !strings.HasPrefix(cfg.Contexts[1].Name, "02-") {
		t.FailNow()
	}

	if !strings.HasPrefix(cfg.Contexts[2].Name, "03-") {
		t.FailNow()
	}
}

func Test_Import(t *testing.T) {
	InitTestdata()

	// add configuration 2 into configuration 1
	cfg1, _, _ := Open(Abs("import-1.yaml"))
	cfg2, _, _ := Open(Abs("import-2.yaml"))

	cfg1.Import(cfg2, &ImportOptions{})
	cfg1.Save(Abs("import-result.yaml"))

	// validate users
	users := cfg1.AuthInfos
	if len(users) != 2 {
		t.FailNow()
	}

	clientKey := strings.TrimRight(string(users[1].AuthInfo.ClientKeyData), "\n")
	if clientKey != "userdata2" {
		t.FailNow()
	}

	// validate clusters
	clusters := cfg1.Clusters
	if len(clusters) != 2 {
		t.FailNow()
	}

	if clusters[1].Name != "cluster-2" {
		t.FailNow()
	}

	// validate contexts
	contexts := cfg1.Contexts
	if len(contexts) != 2 {
		t.FailNow()
	}
}

func Test_RemoveContext(t *testing.T) {
	InitTestdata()

	kcfg, _, err := Open(Abs("remove-test.yaml"))
	if err != nil {
		t.FailNow()
	}

	kcfg.Remove("green")

	if len(kcfg.Contexts) != 2 {
		t.FailNow()
	}

	if len(kcfg.AuthInfos) != 2 {
		t.FailNow()
	}

	if len(kcfg.Clusters) != 2 {
		t.FailNow()
	}
}

func Test_RenameContet(t *testing.T) {
	InitTestdata()

	kcfg, _, err := Open(Abs("rename-test.yaml"))
	if err != nil {
		t.FailNow()
	}

	kcfg.Rename("blue", "cyan")

	// blue resources (context, user, cluster) should not exist
	ctx := kcfg.GetContext("blue")
	if ctx != nil {
		t.FailNow()
	}

	cluster := kcfg.GetCluster("blue-cluster")
	if cluster != nil {
		t.FailNow()
	}

	user := kcfg.GetUser("John")
	if user != nil {
		t.FailNow()
	}

	// 'cyan' resource (contex, user, cluster) must exist
	ctx = kcfg.GetContext("cyan")
	if ctx == nil {
		t.FailNow()
	}

	user = kcfg.GetUser("cyan")
	if user == nil {
		t.FailNow()
	}

	cluster = kcfg.GetCluster("cyan")
	if cluster == nil {
		t.FailNow()
	}

	// cyan context must refer to correct user and cluster
	if ctx.Context.Cluster != "cyan" {
		t.FailNow()
	}

	if ctx.Context.AuthInfo != "cyan" {
		t.FailNow()
	}
}

func Test_Split(t *testing.T) {
	InitTestdata()

	kcfg, _, err := Open(Abs("split.yaml"))
	if err != nil {
		t.FailNow()
	}

	splitted := kcfg.Split()
	if len(splitted) != 3 {
		t.FailNow()
	}

	testingKcfg := splitted[1]
	if testingKcfg.CurrentContext != testingKcfg.Contexts[0].Name {
		t.FailNow()
	}

	if testingKcfg.Clusters[0].Name != testingKcfg.Contexts[0].Context.Cluster {
		t.FailNow()
	}
}
