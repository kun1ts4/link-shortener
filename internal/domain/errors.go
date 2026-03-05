package domain

import (
	"fmt"
)

var (
	ErrNotFound      = fmt.Errorf("not found")
	ErrInvalid       = fmt.Errorf("invalid input")
	ErrAlreadyExists = fmt.Errorf("already exists")
)
