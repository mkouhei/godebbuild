package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"text/template"
)

func (c *config) setBasepath() {
	if c.Codename == "sid" {
		c.Basepath = path.Dir(fmt.Sprintf("/var/cache/pbuilder/base.cow/"))
	} else {
		c.Basepath = path.Dir(fmt.Sprintf("/var/cache/pbuilder/base-%s.cow/", c.Codename))
	}
}

func (c *config) setBasetgz() {
	if c.Codename == "sid" {
		c.Basetgz = fmt.Sprintf("/var/cache/pbuilder/base.tgz")
	} else {
		c.Basetgz = fmt.Sprintf("/var/cache/pbuilder/base-%s.tgz", c.Codename)
	}
}

func (c *config) updateCowbuilder() {
	command := "sudo"
	args := []string{"cowbuilder", "--update", "--basepath", c.Basepath}
	runCommand(command, args...)
}

func (c *config) updatePbuilder() {
	command := "sudo"
	args := []string{"pbuilder", "--update", "--basetgz", c.Basetgz}
	runCommand(command, args...)
}

func (c config) preparePbuilderrc() string {
	const content = `DISTRIBUTION={{.Codename}}
DEBBUILDOPTS="-sa"
BUILDRESULT={{.ResultsDirpath}}
`
	pbuilderrcTempl := template.Must(template.New("").Parse(content))
	buf := &bytes.Buffer{}
	pbuilderrcTempl.Execute(buf, &c)
	pbuilderrcPath := fmt.Sprintf("%s/.pbuilderrc", c.WorkDirpath)
	if err := ioutil.WriteFile(pbuilderrcPath, buf.Bytes(), filePerm); err != nil {
		log.Fatal(err)
	}
	return pbuilderrcPath
}

func buildPkg(pbuilderrcPath string, basepath string, dscPath string) {
	command := "sudo"
	args := []string{"cowbuilder", "--build", "--configfile",
		pbuilderrcPath, "--basepath", basepath, dscPath}
	runCommand(command, args...)
}

func mkBuildDeps(controlFilePath string) {
	command := "sudo"
	args := []string{"mk-build-deps", "-i", "-r", controlFilePath, "-t",
		"'apt-get --force-yes -y'"}
	runCommand(command, args...)
}

func purgeBuildDeps(bldDepsPkgName string) {
	command := "sudo"
	args := []string{"apt-get", "purge", "-y", bldDepsPkgName}
	runCommand(command, args...)
}

func (c *config) gitBuildPkg() {
	command := "sudo"
	exportDirOpt := fmt.Sprintf("--git-export-dir=%s", c.ResultsDirpath)
	gitDistOpt := fmt.Sprintf("--git-dist=%s", c.Codename)
	args := []string{"git-buildpackage", "--git-ignore-branch",
		"--git-pbuilder", exportDirOpt, "-sa", "--git-ignore-new",
		gitDistOpt}
	runCommand(command, args...)
}
