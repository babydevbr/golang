// Package app use cases
package app

import "errors"

var (
	// ErrInvalid input.
	ErrInvalid = errors.New("invalid")
	// ErrOnSave situation.
	ErrOnSave = errors.New("on save")
)
