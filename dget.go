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
	if msg, err := rnr.runCommand(command, args...); err != nil {
		log.Println(msg)
		return err
	}
	return nil
}
