package constants

const (
	STR  string = "STR"
	LIST string = "LIST"
	INT  string = "INT"
)

var (
	AllowedDataTypes = map[string]struct{}{
		STR:  {},
		INT:  {},
		LIST: {},
	}
)
