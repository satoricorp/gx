//go:build darwin || linux || freebsd || netbsd || openbsd || dragonfly

package cli

import (
	"os/exec"
	"syscall"
)

func configureDetachedCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
