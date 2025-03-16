package command

import (
	"goredis/internal/request"
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

func (dc *DelCommand) Execute(req request.Request) (*string, error) {
	return dc.kv.Delete(req)
}
