package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
)

type testRunner struct{}

func (r testRunner) runCommand(command string, args ...string) error {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Print(string(stdout))
	return nil
}

func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	fmt.Println("testing helper process")
}

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
	rnr = realRunner{}
	cmd := "foo"
	args := []string{}
	if err := rnr.runCommand(cmd, args...); err == nil {
		t.Fatal("want: <fail>")
	}
	cmd = "true"
	if err := rnr.runCommand(cmd, args...); err != nil {
		t.Fatal(err)
	}
}

func TestDebError(t *testing.T) {
	if err := debError("test"); err == nil {
		t.Fatal("want: <fail>")
	}
}

func TestChangeOwner(t *testing.T) {
	rnr = testRunner{}
	c := config{}
	c.ResultsDirpath = "/path/to/test"
	c.changeOwner()
}
