package command

import (
	"goredis/internal/request"
	"goredis/internal/response"
	"goredis/internal/store"
)

type (
	DecrCommand struct {
		kv *store.KeyValueStore
	}
)

func NewDecrCommand(kv *store.KeyValueStore) Command {
	return &DecrCommand{
		kv: kv,
	}
}

func (dc *DecrCommand) Execute(req request.Request) *response.Response {
	return dc.kv.Decr(req)
}
