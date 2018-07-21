package store

import (
	"errors"
)

var (
	// ErrNotFound is returned when no element was found under a criteria
	ErrNotFound = errors.New("not found")
	// ErrDuplicate is returned when attempting to insert an item which breaks
	// index key uniqueness
	ErrDuplicate = errors.New("duplicate")
)
