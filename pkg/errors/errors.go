package errors

import "errors"

var (
	ErrNoRows = errors.New("record not found")
)
