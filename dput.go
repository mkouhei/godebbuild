package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ThomasRooney/gexpect"
)

func dputCheck(changesPath string, withoutLintian bool) error {
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
	if msg, err := rnr.runCommand(command, args...); err != nil {
		log.Println(msg)
		return err
	}
	return nil
}

func (c *config) dput(changesPath string, passphrase string, withoutLintian bool) {
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
