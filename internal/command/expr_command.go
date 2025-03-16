package command

import (
	"goredis/internal/request"
	"goredis/internal/store"
)

type (
	ExprCommand struct {
		kv *store.KeyValueStore
	}
)

func NewExprCommand(kv *store.KeyValueStore) Command {
	return &ExprCommand{
		kv: kv,
	}
}

func (ex *ExprCommand) Execute(req request.Request) (*string, error) {
	return ex.kv.SetExpiration(req)
}
