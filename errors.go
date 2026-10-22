package engage

import "errors"

var (
	ErrPostNotFound = errors.New("post not found")
	ErrInvalidPost  = errors.New("invalid post")
	ErrNoPoster     = errors.New("a post needs a caller to stand behind it")
)
