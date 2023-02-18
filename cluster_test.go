package kconf

import (
	"testing"
)

func Test_ParseClusters(t *testing.T) {
	c, err := OpenFile("testdata/clusters.yaml")
	if err != nil {
		t.FailNow()
	}

	if c.Clusters[0].Name != "c1" {
		t.FailNow()
	}

	if c.Clusters[0].Cluster.CertificateAuthorityData != "DATA" {
		t.FailNow()
	}
}
