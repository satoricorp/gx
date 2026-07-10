//go:build windows

package background

import "os/exec"

func configureDetachedCommand(cmd *exec.Cmd) {}
