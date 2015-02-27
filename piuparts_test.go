package main

import (
	"testing"
)

func TestPiuparts(t *testing.T) {
	c := &config{}
	c.Codename = "jessie"
	c.Flavor = "debian"
	c.Basetgz = "/path/to/base-jessie.tgz"
	rnr = testRunner{}
	c.piuparts("/path/to/changes", "http://example.org/mirror", false)
}
