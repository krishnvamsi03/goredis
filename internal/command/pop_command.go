package command

import (
	"goredis/internal/request"
	"goredis/internal/response"
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

func (pp *PopCommand) Execute(req request.Request) *response.Response {
	return pp.kv.Pop(req)
}
