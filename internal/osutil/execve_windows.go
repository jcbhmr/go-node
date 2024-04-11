package osutil

import (
	"errors"
	"os"
	"os/exec"
)

func Execve(cmd string, argv []string, envv []string) error {
	c := exec.Cmd{Path: cmd, Args: argv, Env: envv}
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if exitErr := new(exec.ExitError); errors.As(err, &exitErr) {
		os.Exit(exitErr.ExitCode())
	}
	return err
}