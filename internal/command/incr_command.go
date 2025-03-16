package command

import (
	"goredis/internal/request"
	"goredis/internal/store"
)

type (
	IncrCommand struct {
		kv *store.KeyValueStore
	}
)

func NewIncrCommand(kv *store.KeyValueStore) Command {
	return &IncrCommand{
		kv: kv,
	}
}

func (ic *IncrCommand) Execute(req request.Request) (*string, error) {
	return ic.kv.Incr(req)
}
