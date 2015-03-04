package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
)

type runner interface {
	runCommand(string, ...string) (string, error)
}
type realRunner struct{}

func (r realRunner) runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	}
	return stdout.String(), nil
}

func (c *config) cleanDirs() {
	if err := os.RemoveAll(c.TempDirpath); err != nil {
		log.Fatal(err)
	}
	if err := os.RemoveAll(c.ResultsDirpath); err != nil {
		log.Fatal(err)
	}
}

func (c *config) changeOwner() error {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	command := "sudo"
	args := []string{"chown", "-R", fmt.Sprintf("%s:", u.Username), c.ResultsDirpath}
	if msg, err := rnr.runCommand(command, args...); err != nil {
		log.Println(msg)
		return err
	}
	return nil
}

func workDirpath() (string, error) {
	wd := os.Getenv("WORKSPACE")
	if wd == "" {
		return "", debError("Not set WORKSPACE environment variable")
	}
	if _, err := ioutil.ReadDir(wd); err != nil {
		return "", err
	}
	return wd, nil
}

func curdir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[cwd]: %s\n", cwd)
	return cwd
}

func debError(err string) error {
	return errors.New(err)
}
