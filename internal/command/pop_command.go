package command

import (
	"goredis/internal/request"
	"goredis/internal/store"
)

type (
	PopCommand struct {
		kv *store.KeyValueStore
	}
)

func NewPopCommand(kv *store.KeyValueStore) Command {
	return &PopCommand{
		kv: kv,
	}
}

func (pp *PopCommand) Execute(req request.Request) (*string, error) {
	return pp.kv.Pop(req)
}
