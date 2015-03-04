package main

import (
	"log"
)

func (c *config) piuparts(changesPath string, mirror string, noUpgradeTest bool) error {
	command := "sudo"
	args := []string{"piuparts", "-d", c.Codename, "-D", c.Flavor,
		"--basetgz", c.Basetgz}
	if mirror != "" {
		args = append(args, "-m")
		args = append(args, mirror)
	}
	if noUpgradeTest == true {
		args = append(args, "--no-upgrade-test")
	}
	args = append(args, changesPath)
	if msg, err := rnr.runCommand(command, args...); err != nil {
		log.Println(msg)
		return err
	}
	return nil
}
