package client

import "errors"

var (
	ErrInvalidCommand = errors.New("invalid command received")
	ErrInvalidArgs    = errors.New("invalid number of arguements received")
)
