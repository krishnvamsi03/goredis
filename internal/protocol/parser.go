package protocol

import "bufio"

type (
	Parser interface {
		Parse(reader *bufio.Reader) error
	}
)
