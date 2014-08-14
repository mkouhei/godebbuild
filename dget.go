package main

import (
	"log"
	"os"
)

func (c *config) retrieveSrcPkg(dscUrl string) error {
	if err := os.Chdir(c.TempDirpath); err != nil {
		log.Fatal(err)
	}
	command := "dget"
	args := []string{"-d", dscUrl}
	return runCommand(command, args...)
}
