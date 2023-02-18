package kconf

import "testing"

func Test_ParseContext(t *testing.T) {
	cfg, err := OpenFile("testdata/contexts.yaml")
	if err != nil {
		t.FailNow()
	}

	ctx := cfg.Contexts[0]

	if ctx.Name != "blue" {
		t.FailNow()
	}

	if ctx.Context.Cluster != "blue-cluster" {
		t.FailNow()
	}

	if ctx.Context.User != "John" {
		t.FailNow()
	}

	if ctx.Context.Namespace != "team-a" {
		t.FailNow()
	}
}
