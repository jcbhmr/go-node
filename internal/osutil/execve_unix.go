//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package osutil

import (
	"golang.org/x/sys/unix"
)

func Execve(cmd string, argv []string, envv []string) error {
	return unix.Exec(cmd, argv, envv)
}
