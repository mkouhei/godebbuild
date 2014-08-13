package main

import (
	"testing"
)

var (
	exampleDscUrl  = "http://http.debian.net/debian/pool/main/e/example/example_0.1-1.dsc"
	exampleDscName = "example_0.1-1.dsc"
)

func TestDscName(t *testing.T) {
	if dscName, err := DscName(exampleDscUrl); err != nil {
		t.Fatalf("%v, want: %s is example_0.1-1.dsc", err, dscName)
	}
}

func TestPkgName(t *testing.T) {
	if pkgName := PkgName(exampleDscName); pkgName != "example" {
		t.Fatalf("%v, want: example", pkgName)
	}
}

func TestArchitecure(t *testing.T) {
	if arch := Architecture(); arch != "amd64" {
		t.Fatalf("%v, want: amd64", arch)
	}
}
