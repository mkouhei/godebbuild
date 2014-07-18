package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/ThomasRooney/gexpect"
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
	WorkDirpath     string
	TempDirpath     string
	UnsignedDirpath string
	SignedDirpath   string
	Flavor          string
	Codename        string
	Basetgz         string
	Basepath        string
}

func workDirpath() string {
	workDirpath := os.Getenv("WORKSPACE")
	if workDirpath == "" {
		log.Fatal("Not set WORKSPACE environment variable")
	}
	if _, err := ioutil.ReadDir(workDirpath); err != nil {
		log.Fatal(err)
	}
	return path.Dir(workDirpath)
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
		log.Print(string(buf[0:n]))
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
BINDMOUNTS={{.TempDirpath}}
DEBBUILDOPTS="-sa"
BUILDRESULT={{.UnsignedDirpath}}
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

func retrieveSrcPkg(dscUrl string) {
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
	exportDirOpt := fmt.Sprintf("--git-exportdir=%s", c.UnsignedDirpath)
	gitDistOpt := fmt.Sprintf("--git-dist=%s", c.Codename)
	args := []string{"git-buildpackage", "--git-ignore-branch",
		"--git-pbuilder", exportDirOpt, "-sa", "--git-ignore-new",
		gitDistOpt}
	runCommand(command, args...)
}

func (c *config) findDscName() string {
	var dscName string
	fis, err := ioutil.ReadDir(c.UnsignedDirpath)
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
	flag.Parse()

	subcmd := flag.Args()
	if subcmd[0] != "backport" && subcmd[0] != "original" {
		log.Fatal("usage: debbuild [options] <backport|original>")
	}

	workDirpath := workDirpath()
	cfg := &config{workDirpath,
		path.Dir(fmt.Sprintf("%s/temp/", workDirpath)),
		path.Dir(fmt.Sprintf("%s/unsigned_results/", workDirpath)),
		path.Dir(fmt.Sprintf("%s/signed_results/", workDirpath)),
		*f, *c, "", ""}

	os.Mkdir(cfg.TempDirpath, dirPerm)
	os.Mkdir(cfg.UnsignedDirpath, dirPerm)
	os.Mkdir(cfg.SignedDirpath, dirPerm)

	cfg.setBasepath()
	cfg.setBasetgz()

	cfg.updateCowbuilder()
	pbuilderrcPath := cfg.preparePbuilderrc()

	var dscName string
	if subcmd[0] == "backport" {
		// backport
		dscName = DscName(*u)
		dscPath := fmt.Sprintf("%s/%s", cfg.WorkDirpath, dscName)
		retrieveSrcPkg(*u)
		buildPkg(pbuilderrcPath, cfg.Basepath, dscPath)
	} else if subcmd[0] == "original" {
		// original
		cfg.gitBuildPkg()
		dscName = cfg.findDscName()
	}

	arch := Architecture()
	changesName := ChangesName(dscName, arch)
	cfg.updatePbuilder()
	unsignedChangesPath := fmt.Sprintf("%s/%s", cfg.UnsignedDirpath, changesName)
	cfg.Piuparts(unsignedChangesPath, *m, *n)
	cfg.changeOwner(cfg.UnsignedDirpath)
	Debsign(unsignedChangesPath, *p)
	if *b == true {
		DputCheck(unsignedChangesPath, *w)
	} else {
		cfg.Dput(unsignedChangesPath, *r, *w)
	}
}
