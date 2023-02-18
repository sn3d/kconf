package kconf

import (
	"io/ioutil"
	"testing"
)

func Test_Import(t *testing.T) {

	// add configuration 2 into configuration 1
	cfg1, _ := OpenFile("testdata/cfg1.yaml")
	cfg2, _ := OpenFile("testdata/cfg2.yaml")

	cfg1.Import(cfg2)

	// validate users
	users := cfg1.Users
	if len(users) != 2 {
		t.FailNow()
	}

	if users[1].User.ClientKeyData != "userdata2" {
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

func Test_OpenBase64(t *testing.T) {
	b64Data, err := ioutil.ReadFile("testdata/b64.yaml")
	if err != nil {
		t.FailNow()
	}

	kcfg, err := OpenBase64(b64Data)
	if err != nil {
		t.FailNow()
	}

	if len(kcfg.Clusters) != 1 {
		t.FailNow()
	}
}

func Test_Export(t *testing.T) {
	kcfg, err := OpenFile("testdata/big.yaml")
	if err != nil {
		t.FailNow()
	}

	greenCfg, err := kcfg.Export("green")
	if err != nil {
		t.FailNow()
	}

	if greenCfg.Contexts[0].Name != "green" {
		t.FailNow()
	}

	if len(greenCfg.Users) != 1 {
		t.FailNow()
	}

	if len(greenCfg.Clusters) != 1 {
		t.FailNow()
	}
}

func Test_RemoveContext(t *testing.T) {
	kcfg, err := OpenFile("testdata/remove.yaml")
	if err != nil {
		t.FailNow()
	}

	kcfg.RemoveContext("green")

	if len(kcfg.Contexts) != 2 {
		t.FailNow()
	}

	if len(kcfg.Users) != 2 {
		t.FailNow()
	}

	if len(kcfg.Clusters) != 2 {
		t.FailNow()
	}
}
