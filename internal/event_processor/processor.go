package event_processor

import (
	"goredis/internal/command"
	"goredis/internal/response"
	statuscodes "goredis/internal/status_codes"
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

func (p *processor) Start() {
	p.commandManager.Start()
}

func (p *processor) Stop() {
	p.commandManager.Stop()
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
