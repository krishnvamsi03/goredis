package tokens

import (
	"goredis/internal/request"
)

// commands
var (
	GRESP          string = "GRESP"
	KEY            string = "KEY"
	OP             string = "OP"
	EXPR           string = "EXPR"
	DATA_TYPE      string = "DATA_TYPE"
	CONTENT_LENGTH string = "CONTENT_LENGTH"
)

var (
	MANDATORY_TOKENS = map[string]struct{}{
		GRESP: {},
		KEY:   {},
		OP:    {},
	}

	COMMANDS = map[string]func(string) request.RequestOptions{
		KEY:       request.WithKey,
		OP:        request.WithOp,
		DATA_TYPE: request.WithDataType,
		EXPR:      request.WithExpr,
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
)
