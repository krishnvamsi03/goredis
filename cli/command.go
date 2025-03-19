package client

type (
	command struct {
		op       string
		key      string
		value    string
		datatype string
		ttl      string
	}
)