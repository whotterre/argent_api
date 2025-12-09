package customErrors

import "errors"

var (
	ErrInsufficientFunds = errors.New("balance less than 0")
)