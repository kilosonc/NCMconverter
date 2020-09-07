package ncm

import (
	"errors"
)

// ErrNcmFormat represents the file is not a ncm file
var ErrNcmFormat = errors.New("file is not NCM format")

// ErrExtNcm represents the file has not ext .ncm
var ErrExtNcm = errors.New("file should have ext .ncm")

// ErrMagicHeader represents the file's header doesn't match required header
var ErrMagicHeader = errors.New("magic header does not match")
