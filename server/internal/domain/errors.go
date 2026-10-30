package domain

import "errors"

var (
	ErrConflict         = errors.New("domain: conflict")
	ErrNotFound         = errors.New("domain: not found")
	ErrUnprocessable    = errors.New("domain: unprocessable")
	ErrValidationFailed = errors.New("domain: validation failed")
)
