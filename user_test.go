package kconf

import (
	"testing"
)

func Test_ParseUsers(t *testing.T) {
	u, err := OpenFile("testdata/users.yaml")
	if err != nil {
		t.FailNow()
	}

	if u.Users[0].Name != "u1" {
		t.FailNow()
	}

	if u.Users[0].User.ClientCertificateData != "d1" {
		t.FailNow()
	}

	if u.Users[0].User.ClientKeyData != "k1" {
		t.FailNow()
	}
}
