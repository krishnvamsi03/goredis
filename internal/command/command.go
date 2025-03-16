package command

import (
	"goredis/internal/constants"
	"goredis/internal/request"
	"goredis/internal/store"
)

type (
	Command interface {
		Execute(request.Request) (*string, error)
	}

	CommandManager struct {
		kv       *store.KeyValueStore
		commands map[string]Command
	}
)

func NewCommandManager() *CommandManager {

	kv := store.NewKeyValueStore()

	commands := map[string]Command{
		constants.SET:  NewAddCommand(kv),
		constants.GET:  NewGetCommand(kv),
		constants.DEL:  NewDelCommand(kv),
		constants.EXPR: NewExprCommand(kv),
		constants.PUSH: NewPushCommand(kv),
		constants.POP:  NewPopCommand(kv),
		constants.INCR: NewIncrCommand(kv),
		constants.DECR: NewDecrCommand(kv),
	}

	return &CommandManager{
		kv:       kv,
		commands: commands,
	}
}

func (cm *CommandManager) Start() {
	cm.kv.InitKvStore()
}

func (cm *CommandManager) Stop() {
	cm.kv.Close()
}

func (cm *CommandManager) Execute(req request.Request) (*string, error) {

	command := cm.commands[*req.Op]
	return command.Execute(req)
}
