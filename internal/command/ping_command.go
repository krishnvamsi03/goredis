package command

import (
	"goredis/internal/request"
	"goredis/internal/response"
	"goredis/internal/store"
)

type (
	PingCommand struct {
		kv *store.KeyValueStore
	}
)

func NewPingCommand(kv *store.KeyValueStore) Command {
	return &PingCommand{
		kv: kv,
	}
}

func (pc *PingCommand) Execute(req request.Request) *response.Response {
	return pc.kv.Ping(req)
}
