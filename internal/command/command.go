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
		commands map[string]Command
	}
)

func NewCommandManager() *CommandManager {

	kv := store.NewKeyValueStore()

	commands := map[string]Command{
		constants.SET: NewAddCommand(kv),
		constants.GET: NewGetCommand(kv),
	}

	return &CommandManager{
		commands: commands,
	}
}

func (cm *CommandManager) Execute(req request.Request) (*string, error) {

	command := cm.commands[*req.Op]
	return command.Execute(req)
}
