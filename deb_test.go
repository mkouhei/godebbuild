package main

import (
	"testing"
)

var (
	exampleDscUrl = "http://http.debian.net/debian/pool/main/e/example/example_0.1-1.dsc"
)

func TestDscName(t *testing.T) {
	if dscName, err := DscName(exampleDscUrl); err != nil {
		t.Fatal("%v, want: %s is example_0.1-1.dsc", err, dscName)
	}
}
