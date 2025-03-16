package command

import (
	"goredis/internal/request"
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

func (dc *DecrCommand) Execute(req request.Request) (*string, error) {
	return dc.kv.Decr(req)
}
