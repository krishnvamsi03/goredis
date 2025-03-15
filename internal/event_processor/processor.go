package event_processor

import (
	"goredis/internal/command"
	"goredis/internal/response"
)

type (
	processor struct {
		commandManager *command.CommandManager
	}
)

func newProcessor() *processor {
	return &processor{
		commandManager: command.NewCommandManager(),
	}
}

func (p *processor) Process(event *Event) {
	if event.err != nil {
		event.conn.Write([]byte(response.BuildResponseWithError(event.err)))
		return
	}

	res, err := p.commandManager.Execute(*event.cmd)

	if err != nil {
		msg := response.BuildResponseWithError(err)
		event.conn.Write([]byte(msg))
		return
	}

	event.conn.Write([]byte(response.BuildResponseWithMsg(*res)))
}
