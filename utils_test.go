package main

import (
	"os"
	"testing"
)

func TestWorkDirPath(t *testing.T) {
	os.Setenv("WORKSPACE", "/tmp/foo")
	os.Mkdir("/tmp/foo", 0600)
	os.Mkdir("/tmp/bar", 0600)
	wd := workDirpath()
	if wd != "/tmp/foo" {
		t.Fatal(wd)
	}

	c := &config{}
	c.TempDirpath = "/tmp/foo"
	c.ResultsDirpath = "/tmp/bar"

	c.cleanDirs()
}
