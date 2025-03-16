package command

import (
	"goredis/internal/request"
	"goredis/internal/response"
	"goredis/internal/store"
)

type (
	GetKeyCommand struct {
		kv *store.KeyValueStore
	}
)

func NewGetKeyCommand(kv *store.KeyValueStore) Command {
	return &GetKeyCommand{
		kv: kv,
	}
}

func (gkc *GetKeyCommand) Execute(req request.Request) *response.Response {
	return gkc.kv.GetKey(req)
}
