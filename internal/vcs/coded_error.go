package vcs

// CodedError carries a process exit code for CLI callers, so a command can
// fail with a specific status instead of a generic 1.
type CodedError struct {
	Code int
	Err  error
}

func (e *CodedError) Error() string { return e.Err.Error() }
func (e *CodedError) Unwrap() error { return e.Err }

// CodedErrorf wraps err so the CLI exits with code.
func CodedErrorf(code int, err error) error {
	if err == nil {
		return nil
	}
	return &CodedError{Code: code, Err: err}
}
