package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os/exec"
	"strings"
)

func PkgName(dscName string) string {
	p := strings.Split(dscName, "_")
	return p[len(p)-1]
}

func DscName(rawurl string) (string, error) {
	dscUrl, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	var p []string
	p = strings.Split(dscUrl.Path, "/")
	return p[len(p)-1], nil
}

func (c *config) findDscName() string {
	var dscName string
	fis, err := ioutil.ReadDir(c.ResultsDirpath)
	if err != nil {
		log.Fatal(err)
	}
	if len(fis) == 0 {
		log.Fatal(err)
	}
	for _, fi := range fis {
		if strings.HasSuffix(fi.Name(), ".dsc") == true {
			dscName = fi.Name()
			break
		}
	}
	return dscName
}

func Architecture() string {
	cmd := exec.Command("dpkg-architecture", "-qDEB_BUILD_ARCH")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimRight(string(out), "\n")
}

func ChangesName(dscName string, arch string) string {
	pkgName := strings.Split(dscName, ".dsc")[0]
	return fmt.Sprintf("%s_%s.changes", pkgName, arch)
}
