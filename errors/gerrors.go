package gerrors

import "errors"

var (
	ErrInvalidProtocol    = errors.New("invalid protocol received, some content of data are missing or invalid")
	ErrIncomReqBody       = errors.New("invalid or incompleted request body")
	ErrInvalidCmds        = errors.New("invalid commands received")
	ErrInvalidContentType = errors.New("invalid content type received empty or not received")
	ErrContentMismatch    = errors.New("content length provide and content mismatch")
)
