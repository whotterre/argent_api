package customErrors

import "errors"

var (
	ErrorActiveAPIKeysExceeded = errors.New("you have already 5 active API keys")
	ErrInvalidPermission = errors.New("invalid permission passed")
	ErrHashingAPIKey = errors.New("failed to hash API key")
	ErrNonExistentAPIKey = errors.New("API key doesn't exist")
	ErrRollingOverNotExpiredKey = errors.New("API key being rolled over hasn't expired yet")
)