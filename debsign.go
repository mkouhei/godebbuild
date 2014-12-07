package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ThomasRooney/gexpect"
)

func debsign(changesPath string, passphrase string) {
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
