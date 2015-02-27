package main

import (
	"os"
	"testing"
)

func TestRetrieveSrcPkg(t *testing.T) {
	c := &config{}
	c.TempDirpath = "temp"
	os.Mkdir(c.TempDirpath, dirPerm)
	rnr = testRunner{}
	if err := c.retrieveSrcPkg("http://example.org/dummy.dsc"); err != nil {
		t.Fatal(err)
	}
	os.RemoveAll(c.TempDirpath)
}
