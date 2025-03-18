package command

import (
	"goredis/internal/constants"
	"goredis/internal/request"
	"goredis/internal/response"
	statuscodes "goredis/internal/status_codes"
	"goredis/internal/store"
)

type (
	Command interface {
		Execute(request.Request) *response.Response
	}

	CommandManager struct {
		kv       *store.KeyValueStore
		commands map[string]Command
	}
)

func NewCommandManager(kv *store.KeyValueStore) *CommandManager {

	commands := map[string]Command{
		constants.PING: NewPingCommand(kv),
		constants.SET:  NewAddCommand(kv),
		constants.GET:  NewGetCommand(kv),
		constants.DEL:  NewDelCommand(kv),
		constants.EXPR: NewExprCommand(kv),
		constants.PUSH: NewPushCommand(kv),
		constants.POP:  NewPopCommand(kv),
		constants.INCR: NewIncrCommand(kv),
		constants.DECR: NewDecrCommand(kv),
		constants.KEYS: NewGetKeyCommand(kv),
	}

	return &CommandManager{
		kv:       kv,
		commands: commands,
	}
}

func (cm *CommandManager) Execute(req request.Request) *response.Response {

	command, ok := cm.commands[*req.Op]
	if !ok {
		return response.NewResponse().
			WithCode(statuscodes.UNKNOWN_COMMAND).
			WithOk(false).
			WithRes("unknown command")
	}
	return command.Execute(req)
}
