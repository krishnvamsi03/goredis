package command

import (
	"goredis/internal/request"
	"goredis/internal/response"
	"goredis/internal/store"
)

type (
	PushCommand struct {
		kv *store.KeyValueStore
	}
)

func NewPushCommand(kv *store.KeyValueStore) Command {
	return &PushCommand{
		kv: kv,
	}
}

func (pc *PushCommand) Execute(req request.Request) *response.Response {
	return pc.kv.Push(req)
}
