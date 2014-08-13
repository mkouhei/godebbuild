package main

import (
	"os"
	"testing"
)

func TestWorkDirPath(t *testing.T) {
	var (
		wd  string
		err error
		tmp = "temp"
	)

	if wd, err = workDirpath(); wd != "" || err == nil {
		t.Fatalf("%v, want: %s", wd, "<empty>")
	}
	os.Setenv("WORKSPACE", tmp)
	os.Mkdir(tmp, 0600)
	if wd, err = workDirpath(); err != nil {
		t.Fatalf("%v, want: %s", err, "")
	}
	if wd != tmp {
		t.Fatalf("%v, want: %s", wd, tmp)
	}

	c := &config{}
	c.TempDirpath = tmp
	c.cleanDirs()
}
