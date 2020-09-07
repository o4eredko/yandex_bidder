package domain

import "errors"

var (
	ErrGroupNotFound       = errors.New("group not found")
	ErrJobNotFound         = errors.New("job not found")
	ErrJobAlreadyScheduled = errors.New("job already scheduled")
)
