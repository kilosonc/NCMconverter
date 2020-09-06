package ncm

import (
	"errors"
)

var ErrNcmFormat = errors.New("File is not NCM format!")
var ErrExtNcm = errors.New("File should have ext ncm!")
var ErrMagicHeader = errors.New("Magic header does not match!")
