package command

import (
	"goredis/internal/request"
	"goredis/internal/response"
	"goredis/internal/store"
)

type (
	AddCommand struct {
		kvStore *store.KeyValueStore
	}
)

func NewAddCommand(kv *store.KeyValueStore) Command {
	return &AddCommand{
		kvStore: kv,
	}
}

func (ac *AddCommand) Execute(req request.Request) *response.Response {
	return ac.kvStore.Add(req)
}
