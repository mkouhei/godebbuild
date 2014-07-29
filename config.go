package main

import (
	"log"

	"github.com/miguel-branco/goconfig"
)

type config struct {
	WorkDirpath    string
	TempDirpath    string
	ResultsDirpath string
	Flavor         string
	Codename       string
	Basetgz        string
	Basepath       string
}

func readConfig(configPath string) map[string]string {
	c, err := goconfig.ReadConfigFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	p, err := c.GetString("pgp", "passphrase")
	if err != nil {
		p = ""
	}
	r, err := c.GetString("reprepro", "passphrase")
	if err != nil {
		r = ""
	}
	pass := map[string]string{
		"passphrase":  p,
		"rPassphrase": r,
	}
	return pass
}
