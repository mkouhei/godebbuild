package main

import "testing"

func TestReadConfig(t *testing.T) {
	c := readConfig("examples/debbuild.conf")
	if c["passphrase"] != "secret" {
		t.Fatal("parse error [pgp]passphrase")
	}
	if c["rPassphrase"] != "secret" {
		t.Fatal("parse error [reprepro]passphrase")
	}
}
