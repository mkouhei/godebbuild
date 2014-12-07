package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var (
	exampleDscURL      = "http://http.debian.net/debian/pool/main/e/example/example_0.1-1.dsc"
	exampleDscName     = "example_0.1-1.dsc"
	exampleChangesName = "example_0.1-1_amd64.changes"
)

func TestDscName(t *testing.T) {
	if dscName, err := dscName(exampleDscURL); err != nil {
		t.Fatalf("%v, want: %s is example_0.1-1.dsc", err, dscName)
	}
}

func TestPkgName(t *testing.T) {
	if pkgName := pkgName(exampleDscName); pkgName != "example" {
		t.Fatalf("%v, want: example", pkgName)
	}
}

func TestArchitecure(t *testing.T) {
	if arch := architecture(); arch != "amd64" {
		t.Fatalf("%v, want: amd64", arch)
	}
}

func TestFindDscName(t *testing.T) {
	c := &config{}
	c.ResultsDirpath = "temp"
	os.Mkdir(c.ResultsDirpath, dirPerm)
	if dscName, err := c.findDscName(); dscName != "" || err == nil {
		t.Fatalf("%v, want: %s", dscName, "<empty>")
	}
	ioutil.WriteFile(fmt.Sprintf("temp/%s", exampleDscName), []byte(""), filePerm)
	if dscName, err := c.findDscName(); dscName == "" || err != nil {
		t.Fatalf("%v, want: %s\n", err, exampleDscName)
	}
	c.TempDirpath = "temp"
	c.cleanDirs()
}

func TestChangesName(t *testing.T) {
	if changes := changesName(exampleDscName, "amd64"); changes != exampleChangesName {
		t.Fatalf("%v, want: %s", changes, exampleChangesName)
	}
}
