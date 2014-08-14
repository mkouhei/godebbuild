package main

import (
	"testing"
)

var (
	basepath       = "/var/cache/pbuilder/base.cow"
	basepathJessie = "/var/cache/pbuilder/base-jessie.cow"
	basetgz        = "/var/cache/pbuilder/base.tgz"
	basetgzJessie  = "/var/cache/pbuilder/base-jessie.tgz"
)

func TestSetBasepath(t *testing.T) {
	c := config{}
	c.Codename = "sid"
	c.setBasepath()
	if c.Basepath != basepath {
		t.Fatalf("%v, want: %s", c.Basepath, basepath)
	}
	c.Codename = "jessie"
	c.setBasepath()
	if c.Basepath != basepathJessie {
		t.Fatalf("%v, want: %s", c.Basepath, basepathJessie)
	}
}

func TestSetBasetgz(t *testing.T) {
	c := config{}
	c.Codename = "sid"
	c.setBasetgz()
	if c.Basetgz != basetgz {
		t.Fatalf("%v, want: %s", c.Basetgz, basetgz)
	}
	c.Codename = "jessie"
	c.setBasetgz()
	if c.Basetgz != basetgzJessie {
		t.Fatalf("%v, want: %s", c.Basetgz, basetgzJessie)
	}
}
