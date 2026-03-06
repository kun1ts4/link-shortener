package domain

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalid       = errors.New("invalid input")
	ErrAlreadyExists = errors.New("already exists")
	ErrStorageFull   = errors.New("storage is full")
)
