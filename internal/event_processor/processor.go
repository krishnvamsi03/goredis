package event_processor

import (
	"goredis/internal/command"
	"goredis/internal/response"
	statuscodes "goredis/internal/status_codes"
)

type (
	Processor interface {
		Process(event *Event)
	}

	processor struct {
		commandManager *command.CommandManager
	}
)

var _ Processor = (*processor)(nil)

func NewProcessor(cm *command.CommandManager) *processor {
	return &processor{
		commandManager: cm,
	}
}

func (p *processor) Process(event *Event) {
	if event.err != nil {
		res := response.NewResponse().WithCode(statuscodes.INVALID_PROTOCOL).WithOk(false).WithRes(event.err.Error())
		event.conn.Write([]byte(res.Build()))
		return
	}

	res := p.commandManager.Execute(*event.cmd)

	event.conn.Write([]byte(res.Build()))
}
