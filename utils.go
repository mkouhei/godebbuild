package main

import (
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