package command

import (
	"goredis/internal/request"
	"goredis/internal/store"
)

type (
	GetCommand struct {
		kvStore *store.KeyValueStore
	}
)

func NewGetCommand(kv *store.KeyValueStore) Command {
	return &GetCommand{
		kvStore: kv,
	}
}

func (gc *GetCommand) Execute(req request.Request) (*string, error) {
	return gc.kvStore.Get(req)
}
