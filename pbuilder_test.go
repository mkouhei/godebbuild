package main

import (
	"os"
	"testing"
)

var (
	basepath       = "/var/cache/pbuilder/base.cow"
	basepathJessie = "/var/cache/pbuilder/base-jessie.cow"
	basetgz        = "/var/cache/pbuilder/base.tgz"
	basetgzJessie  = "/var/cache/pbuilder/base-jessie.tgz"
)

func TestUpdateCowbuilder(t *testing.T) {
	c := config{}
	rnr = testRunner{}
	c.updateCowbuilder()
}

func TestUpdatePbuilder(t *testing.T) {
	c := config{}
	rnr = testRunner{}
	c.updatePbuilder()
}

func TestBuildPkg(t *testing.T) {
	rnr = testRunner{}
	buildPkg("/path/to/pbuilderrc", "/path/to/base.cow", "/path/to/some.dsc")
}

func TestMkBuildDeps(t *testing.T) {
	rnr = testRunner{}
	mkBuildDeps("/path/to/control")
}

func TestPurgeBuildDeps(t *testing.T) {
	rnr = testRunner{}
	purgeBuildDeps("somepkg")
}

func TestGitBuildPkg(t *testing.T) {
	c := config{}
	rnr = testRunner{}
	c.gitBuildPkg()
}

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

func TestPreparePbuilderrc(t *testing.T) {
	c := &config{}
	c.TempDirpath = "temp"
	c.WorkDirpath = "temp"
	c.Codename = "sid"
	c.ResultsDirpath = "temp"

	os.Mkdir(c.WorkDirpath, dirPerm)
	if p := c.preparePbuilderrc(); p != "temp/.pbuilderrc" {
		t.Fatalf("%v, want: temp/.pbuilderrc", p)
	}
	c.cleanDirs()
}
