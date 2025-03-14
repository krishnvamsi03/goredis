package gerrors

import "errors"

var (
	ErrInvalidProtocol = errors.New("invalid protocol received")
	ErrIncomReqBody = errors.New("invalid or incompleted request body")
)