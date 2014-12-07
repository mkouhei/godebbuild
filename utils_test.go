package main

import (
	"os"
	"path"
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

func TestCurdir(t *testing.T) {
	if cwd := path.Base(curdir()); cwd != "godebbuild" {
		t.Fatalf("%v, want: godebbuild", cwd)
	}
}

func TestRunCommand(t *testing.T) {
	cmd := "foo"
	args := []string{}
	if err := runCommand(cmd, args...); err == nil {
		t.Fatal("want: <fail>")
	}
	cmd = "true"
	if err := runCommand(cmd, args...); err != nil {
		t.Fatal(err)
	}
}

func TestDebError(t *testing.T) {
	if err := debError("test"); err == nil {
		t.Fatal("want: <fail>")
	}
}
