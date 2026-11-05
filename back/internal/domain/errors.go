package domain

import "errors"

type ErrKey int

const (
	ErrConflict ErrKey = iota
	ErrInvalid
	ErrNotFound
	ErrNotModified
)

type Error struct {
	Reason error
	Key    ErrKey
}

// generic error
var (
	ErrEmptyResult = &Error{
		Reason: errors.New("empty result"),
		Key:    ErrNotFound,
	}
)

// Error implements error interface
func (e *Error) Error() string {
	return e.Reason.Error()
}
