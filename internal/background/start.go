package background

import (
	"os"
	"os/exec"
)

// StartDetachedGX launches the gx binary with args in a detached background process.
func StartDetachedGX(args ...string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(exe, args...)
	if devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0); err == nil {
		cmd.Stdin = devNull
		cmd.Stdout = devNull
		cmd.Stderr = devNull
	}
	cmd.Env = os.Environ()
	configureDetachedCommand(cmd)
	if err := cmd.Start(); err != nil {
		return err
	}
	if cmd.Process != nil {
		_ = cmd.Process.Release()
	}
	return nil
}
