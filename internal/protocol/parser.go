package protocol

import (
	"bufio"
	"goredis/internal/command"
)

type (
	Parser interface {
		Parse(reader *bufio.Reader) (*command.Command, error)
	}
)
