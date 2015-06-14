package git

import (
	"errors"
	"io"
)

// ErrorDeltaNotFound should be returned by Backend implementations
var ErrorDeltaNotFound = errors.New("delta not found")

// A Delta is the difference between two commits
type Delta interface{}

// A Backend for git data
type Backend interface {
	FindDelta(from, to string) (Delta, error)

	GetRefs() ([]Ref, error)

	ReadPackfile(d Delta) (io.ReadCloser, error)

	UpdateRef(update RefUpdate) error

	WritePackfile(from, to string, r io.Reader) error
}
