package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func runCmd(cmdArr []string, varMap tVarMap, print bool) ([]byte, int, error) {
	var err error
	var exitcode int
	var stdBuffer bytes.Buffer

	cmdArr = expandVars(cmdArr, varMap)
	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	mw := io.MultiWriter(&stdBuffer)
	if print == true {
		mw = io.MultiWriter(os.Stdout, &stdBuffer)
	}
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
	fmt.Printf("")

	return stdBuffer.Bytes(), exitcode, err
}

func expandVars(cmdArr []string, varMap tVarMap) []string {
	for idx, el := range cmdArr {
		for key, val := range varMap {
			cmdArr[idx] = strings.Replace(
				el, "{"+key+"}", val, -1,
			)
		}
	}
	return cmdArr
}
