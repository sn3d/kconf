package kconf

import "fmt"

// this error is occured while Open() execution
type OpenError struct {
	original error
	path     string
}

func (e *OpenError) Error() string {
	return fmt.Sprintf("Cannot open kubeconfig file '%s'. Check if --kubeconfig or KUBECONFIG is defined properly", e.path)
}

func (e *OpenError) Unwrap() error {
	return e.original
}
