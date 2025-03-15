package protocol

import (
	"bufio"
	"goredis/internal/request"
)

type (
	Parser interface {
		Parse(reader *bufio.Reader) (*request.Request, error)
	}
)
