package main

func (c *config) piuparts(changesPath string, mirror string, noUpgradeTest bool) {
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
	runCommand(command, args...)
}
