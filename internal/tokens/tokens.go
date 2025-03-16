package tokens

import (
	"goredis/internal/request"
	"strings"
)

// commands
var (
	GRESP          string = "GRESP"
	KEY            string = "KEY"
	OP             string = "OP"
	TTL            string = "TTL"
	DATA_TYPE      string = "DATA_TYPE"
	CONTENT_LENGTH string = "CONTENT_LENGTH"
)

var (
	MANDATORY_TOKENS = map[string]struct{}{
		GRESP: {},
		OP:    {},
	}

	COMMANDS = map[string]func(string) request.RequestOptions{
		KEY:       request.WithKey,
		OP:        request.WithOp,
		DATA_TYPE: request.WithDataType,
		TTL:       request.WithExpr,
	}
)

// instruction

type (
	INST string
)

var (
	GET   INST = "GET"
	SET   INST = "SET"
	DEL   INST = "DEL"
	OEXPR INST = "EXPR"
	PUSH  INST = "PUSH"
	POP   INST = "POP"
	INCR  INST = "INCR"
	DECR  INST = "DECR"
)

func (it INST) ToLower() string {
	return strings.ToLower(string(it))
}
