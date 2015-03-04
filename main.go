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
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

const (
	dirPerm  = 0755
	filePerm = 0644
)

var version string
var showVersion = flag.Bool("version", false, "show_version")
var rnr runner

func main() {

	var (
		dsc  string
		pass map[string]string
		err  error
	)

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
	if *showVersion {
		fmt.Printf("version: %s\n", version)
		return
	}
	rnr = realRunner{}

	subcmd := flag.Args()
	if len(subcmd) == 0 {
		log.Fatal("usage: debbuild [options] <backport|original>")
	}
	if subcmd[0] != "backport" && subcmd[0] != "original" {
		log.Fatal("usage: debbuild [options] <backport|original>")
	}

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
	wd, err := workDirpath()
	if err != nil {
		log.Fatal(err)
	}
	bldDepsPkgName := ""
	os.Chdir(wd)
	cfg := &config{wd,
		path.Dir(fmt.Sprintf("%s/temp/", wd)),
		path.Dir(fmt.Sprintf("%s/results/", wd)),
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

	if subcmd[0] == "backport" {
		// backport
		dsc, err = dscName(*u)
		if err != nil {
			log.Fatal(err)
		}
		dscPath := fmt.Sprintf("%s/%s", cfg.TempDirpath, dsc)
		cfg.retrieveSrcPkg(*u)
		buildPkg(pbuilderrcPath, cfg.Basepath, dscPath)

	} else if subcmd[0] == "original" {
		// original
		os.Chdir(initDirpath)
		mkBuildDeps("debian/control")
		cfg.gitBuildPkg()
		if dsc, err = cfg.findDscName(); err != nil {
			log.Fatal(err)
		}
		bldDepsPkgName = fmt.Sprintf("%s-build-deps", pkgName(dsc))
	}

	arch := architecture()
	changesName := changesName(dsc, arch)
	cfg.updatePbuilder()
	changesPath := fmt.Sprintf("%s/%s", cfg.ResultsDirpath, changesName)
	cfg.piuparts(changesPath, *m, *n)
	cfg.changeOwner()
	debsign(changesPath, pass["passphrase"])
	if *b == true {
		if err := dputCheck(changesPath, *w); err != nil {
			log.Println(err)
		}
	} else {
		cfg.dput(changesPath, pass["rPassphrase"], *w)
	}
	if bldDepsPkgName != "" {
		purgeBuildDeps(bldDepsPkgName)
	}
}
