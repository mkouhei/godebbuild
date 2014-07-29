package main

import (
	"log"
	"os"
)

func (c *config) retrieveSrcPkg(dscUrl string) {
	if err := os.Chdir(c.TempDirpath); err != nil {
		log.Fatal(err)
	}
	command := "dget"
	args := []string{"-d", dscUrl}
	runCommand(command, args...)
}
