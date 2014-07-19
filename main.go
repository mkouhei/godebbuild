package main

/*
  Copyright 2014 Kouhei Maeda

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

	  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/ThomasRooney/gexpect"
	"github.com/miguel-branco/goconfig"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"text/template"
)

const (
	dirPerm  = 0755
	filePerm = 0644
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

func workDirpath() string {
	workDirpath := os.Getenv("WORKSPACE")
	if workDirpath == "" {
		log.Fatal("Not set WORKSPACE environment variable")
	}
	if _, err := ioutil.ReadDir(workDirpath); err != nil {
		log.Fatal(err)
	}
	return workDirpath
}

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

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 1024)
	var n int
	for {
		if n, err = stdout.Read(buf); err != nil {
			break
		}
		fmt.Print(string(buf[0:n]))
	}
	if err == io.EOF {
		err = nil
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

func (c *config) Piuparts(changesPath string, mirror string, noUpgradeTest bool) {
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

func DscName(rawurl string) string {
	dscUrl, err := url.Parse(rawurl)
	if err != nil {
		log.Fatal(err)
	}
	var p []string
	p = strings.Split(dscUrl.Path, "/")
	return p[len(p)-1]
}

func curdir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[cwd]: %s\n", cwd)
	return cwd
}

func (c *config) retrieveSrcPkg(dscUrl string) {
	if err := os.Chdir(c.TempDirpath); err != nil {
		log.Fatal(err)
	}
	command := "dget"
	args := []string{"-d", dscUrl}
	runCommand(command, args...)
}

func buildPkg(pbuilderrcPath string, basepath string, dscPath string) {
	command := "sudo"
	args := []string{"cowbuilder", "--build", "--configfile",
		pbuilderrcPath, "--basepath", basepath, dscPath}
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

func (c *config) changeOwner(dirPath string) {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	command := "sudo"
	args := []string{"chown", "-R", fmt.Sprintf("%s:", u.Username), dirPath}
	runCommand(command, args...)
}

func DputCheck(changesPath string, withoutLintian bool) {
	command := "dput"
	var dputOpts string
	if withoutLintian == true {
		fmt.Println("dput checking without lintian")
		dputOpts = "-o"
	} else {
		fmt.Println("dput checking with lintian")
		dputOpts = "-ol"
	}
	args := []string{dputOpts, changesPath}
	runCommand(command, args...)
}

func (c *config) Dput(changesPath string, passphrase string, withoutLintian bool) {
	os.Setenv("LANG", "C")

	if _, err := ioutil.ReadFile(changesPath); err != nil {
		log.Fatal(err)
	}
	var dputOpts string
	if withoutLintian == true {
		dputOpts = ""
	} else {
		dputOpts = "-l"
	}
	command := fmt.Sprintf("dput %s %s %s", dputOpts, c.Codename, changesPath)
	child, err := gexpect.Spawn(command)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Expecting Please enter passphrase:\n")
	child.Expect("Please enter passphrase:")
	if err := child.SendLine(passphrase); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Expecting Please enter passphrase:\n")
	child.Expect("Please enter passphrase:")
	if err := child.SendLine(passphrase); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Interacting.. \n")
	child.Interact()
	fmt.Printf("Done \n")
	child.Close()
}

func Debsign(changesPath string, passphrase string) {
	os.Setenv("LANG", "C")

	if _, err := ioutil.ReadFile(changesPath); err != nil {
		log.Fatal(err)
	}
	command := fmt.Sprintf("debsign %s", changesPath)
	child, err := gexpect.Spawn(command)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Expecting Enter passphrase:\n")
	child.Expect("Enter passphrase: ")
	if err := child.SendLine(passphrase); err != nil {
		log.Fatal(err)
	}
	child.Expect("Enter passphrase: ")
	if err := child.SendLine(passphrase); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Interacting.. \n")
	child.Interact()
	fmt.Printf("Done \n")
	child.Close()
}

func (c *config) cleanDirs() {
	if err := os.RemoveAll(c.TempDirpath); err != nil {
		log.Fatal(err)
	}
	if err := os.RemoveAll(c.ResultsDirpath); err != nil {
		log.Fatal(err)
	}
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

func main() {

	c := flag.String("c", "sid", "codename")
	f := flag.String("f", "debian", "flavor")
	m := flag.String("m", "", "mirror")
	n := flag.Bool("n", false, "skip tesging upgrade from an existing version in the archive with piuparts")
	w := flag.Bool("w", false, "skip cheking with lintian")
	p := flag.String("p", "", "GPG private key passphfase for debsign")
	r := flag.String("r", "", "GPG private key passphfase for reprepro register")
	b := flag.Bool("b", false, "Build only without dput upload")
	u := flag.String("u", "", ".dsc url for backport")
	cl := flag.Bool("clean", false, "clean results and temp directories.")
	cnf := flag.String("config", "", "configuration file of debbuild")
	flag.Parse()

	subcmd := flag.Args()
	if len(subcmd) == 0 {
		log.Fatal("usage: debbuild [options] <backport|original>")
	}
	if subcmd[0] != "backport" && subcmd[0] != "original" {
		log.Fatal("usage: debbuild [options] <backport|original>")
	}

	var pass map[string]string
	if *cnf != "" {
		pass = readConfig(*cnf)
	}
	if *p != "" {
		pass["passphrase"] = *p
	}
	if *r != "" {
		pass["rPassphrase"] = *r
	}

	initDirpath := curdir()
	workDirpath := workDirpath()
	os.Chdir(workDirpath)
	cfg := &config{workDirpath,
		path.Dir(fmt.Sprintf("%s/temp/", workDirpath)),
		path.Dir(fmt.Sprintf("%s/results/", workDirpath)),
		*f, *c, "", ""}

	if *cl == true {
		cfg.cleanDirs()
	}

	os.Mkdir(cfg.TempDirpath, dirPerm)
	os.Mkdir(cfg.ResultsDirpath, dirPerm)

	cfg.setBasepath()
	cfg.setBasetgz()

	cfg.updateCowbuilder()

	pbuilderrcPath := cfg.preparePbuilderrc()

	var dscName string
	if subcmd[0] == "backport" {
		// backport
		dscName = DscName(*u)
		dscPath := fmt.Sprintf("%s/%s", cfg.TempDirpath, dscName)
		cfg.retrieveSrcPkg(*u)
		buildPkg(pbuilderrcPath, cfg.Basepath, dscPath)

	} else if subcmd[0] == "original" {
		// original
		os.Chdir(initDirpath)
		cfg.gitBuildPkg()
		dscName = cfg.findDscName()
	}

	arch := Architecture()
	changesName := ChangesName(dscName, arch)
	cfg.updatePbuilder()
	changesPath := fmt.Sprintf("%s/%s", cfg.ResultsDirpath, changesName)
	cfg.Piuparts(changesPath, *m, *n)
	cfg.changeOwner(cfg.ResultsDirpath)
	Debsign(changesPath, pass["passphrase"])
	if *b == true {
		DputCheck(changesPath, *w)
	} else {
		cfg.Dput(changesPath, pass["rPassphrase"], *w)
	}
}
