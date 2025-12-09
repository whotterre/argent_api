package customErrors

import "errors"

var (
	ErrorActiveAPIKeysExceeded = errors.New("you have already 5 active API keys")
	ErrInvalidPermission = errors.New("invalid permission passed")
)