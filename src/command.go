package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

func runCmd(cmdArr []string) ([]byte, int, error) {
	var err error
	var exitcode int
	var stdBuffer bytes.Buffer

	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw
	if err = cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// the program has exited with an exit code != 0
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitcode = status.ExitStatus()
			}
		}
	}
	if err != nil {
		fmt.Printf("An error occured: %s\n", err)
	}
	fmt.Printf("")

	return stdBuffer.Bytes(), exitcode, err
}
