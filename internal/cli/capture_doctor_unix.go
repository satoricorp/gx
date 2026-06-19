//go:build darwin || linux

package cli

import "syscall"

type syscallStatfs struct {
	buf syscall.Statfs_t
}

func statfs(path string, s *syscallStatfs) error {
	if err := syscall.Statfs(path, &s.buf); err != nil {
		return err
	}
	return nil
}

func (s syscallStatfs) FreeBytes() uint64 {
	return uint64(s.buf.Bavail) * uint64(s.buf.Bsize)
}
