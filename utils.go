package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
)

func (c *config) cleanDirs() {
	if err := os.RemoveAll(c.TempDirpath); err != nil {
		log.Fatal(err)
	}
	if err := os.RemoveAll(c.ResultsDirpath); err != nil {
		log.Fatal(err)
	}
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

func workDirpath() (string, error) {
	wd := os.Getenv("WORKSPACE")
	if wd == "" {
		return "", Error("Not set WORKSPACE environment variable")
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

func Error(err string) error {
	return errors.New(err)
}
