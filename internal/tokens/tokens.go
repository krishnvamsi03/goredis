package tokens

var (
	GRESP string = "GRESP"
	KEY   string = "KEY"
	OP    string = "OP"
	EXPR  string = "EXPR"
)

var (
	INIT_TOKENS = map[string]struct{}{
		GRESP: {},
		KEY:   {},
		OP:    {},
		EXPR:  {},
	}
)
