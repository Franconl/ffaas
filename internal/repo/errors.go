package repo

import (
	"errors"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrKeyAlreadyUsed = errors.New("flag key already exists")
	ErrKeyRequired    = errors.New("key is required")
	ErrInvalidPercent = errors.New("invalid percentage")
	ErrInvalidBody    = errors.New("invalid JSON body")
)
