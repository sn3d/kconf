package kconf

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sn3d/testdata"
)

func Test_Open(t *testing.T) {
	testdata.Setup()

	cfg, err := Open(testdata.Abs("open-test.yaml"))
	if err != nil {
		t.FailNow()
	}

	fmt.Println(cfg)
	fmt.Println(cfg.Clusters[0].Cluster.CertificateAuthorityData)

	err = cfg.Save(testdata.Abs("save-test.yaml"))
	if err != nil {
		t.FailNow()
	}

	same := testdata.CompareFiles("save-test.yaml", "save-expected.yaml")
	if !same {
		t.FailNow()
	}
}

func Test_Import(t *testing.T) {

	testdata.Setup()

	// add configuration 2 into configuration 1
	cfg1, _ := Open(testdata.Abs("import-1.yaml"))
	cfg2, _ := Open(testdata.Abs("import-2.yaml"))

	cfg1.Import(cfg2)
	cfg1.Save(testdata.Abs("import-result.yaml"))

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
	testdata.Setup()

	kcfg, err := Open(testdata.Abs("remove-test.yaml"))
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
	testdata.Setup()

	kcfg, err := Open(testdata.Abs("rename-test.yaml"))
	if err != nil {
		t.FailNow()
	}

	kcfg.Rename("blue", "cyan")

	// blue resources (context, user, cluster) should not exist
	ctx := kcfg.getContext("blue")
	if ctx != nil {
		t.FailNow()
	}

	cluster := kcfg.getCluster("blue-cluster")
	if cluster != nil {
		t.FailNow()
	}

	user := kcfg.getUser("John")
	if user != nil {
		t.FailNow()
	}

	// 'cyan' resource (contex, user, cluster) must exist
	ctx = kcfg.getContext("cyan")
	if ctx == nil {
		t.FailNow()
	}

	user = kcfg.getUser("cyan")
	if user == nil {
		t.FailNow()
	}

	cluster = kcfg.getCluster("cyan")
	if cluster == nil {
		t.FailNow()
	}
}

func Test_Split(t *testing.T) {
	testdata.Setup()

	kcfg, err := Open(testdata.Abs("split.yaml"))
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
