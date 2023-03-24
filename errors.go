package mongodbr

import "errors"

var (
	ErrInvalidType = errors.New("invalid type")
	ErrNoCursor    = errors.New("no cursor")
)
