package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TesstRetrieveSrcPkg(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
	}))
	defer ts.Close()

	c := &config{}
	c.TempDirpath = "temp"
	os.Mkdir(c.TempDirpath, dirPerm)
	if err := c.retrieveSrcPkg(ts.URL); err != nil {
		t.Fatal(err)
	}
}
