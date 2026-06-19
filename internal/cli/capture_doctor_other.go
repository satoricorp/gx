//go:build !darwin && !linux

package cli

type syscallStatfs struct{}

func statfs(string, *syscallStatfs) error {
	return errDiskStatUnsupported
}

func (s syscallStatfs) FreeBytes() uint64 {
	return 0
}

var errDiskStatUnsupported = errUnsupported("disk stat")

type errUnsupported string

func (e errUnsupported) Error() string { return string(e) }
