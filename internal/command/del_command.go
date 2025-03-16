package command

import (
	"goredis/internal/request"
	"goredis/internal/response"
	"goredis/internal/store"
)

type (
	DelCommand struct {
		kv *store.KeyValueStore
	}
)

func NewDelCommand(kv *store.KeyValueStore) Command {
	return &DelCommand{
		kv: kv,
	}
}

func (dc *DelCommand) Execute(req request.Request) *response.Response {
	return dc.kv.Delete(req)
}
