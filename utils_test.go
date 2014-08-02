package main

import (
	"os"
	"testing"
)

func TestWorkDirPath(t *testing.T) {
	os.Setenv("WORKSPACE", "/tmp")
	wd := workDirpath()
	if wd != "/tmp" {
		t.Fatal(wd)
	}
}
