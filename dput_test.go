package main

import (
	"testing"
)

func TestDputCheck(t *testing.T) {
	rnr = testRunner{}
	dputCheck("/path/to/somepkg.changes", false)
}
