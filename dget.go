package main

import (
	"log"
	"os"
)

func (c *config) retrieveSrcPkg(dscURL string) error {
	if err := os.Chdir(c.TempDirpath); err != nil {
		log.Fatal(err)
	}
	command := "dget"
	args := []string{"-d", dscURL}
	return runCommand(command, args...)
}
