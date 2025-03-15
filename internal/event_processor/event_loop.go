package event_processor

import (
	"goredis/common/logger"
	"goredis/internal/request"
	"net"
	"sync"
)

type (
	Event struct {
		conn net.Conn
		cmd  *request.Request
		err  error
	}

	EventOption func(*Event)

	EventLoop struct {
		queue       chan *Event
		queueClosed bool
		queueLock   *sync.Mutex
		quitch      chan struct{}
		logger      logger.Logger
		processor   *processor
	}
)

func NewEvent(conn net.Conn) *Event {
	return &Event{
		conn: conn,
	}
}

func (e *Event) WithRequest(cmd *request.Request) *Event {
	e.cmd = cmd
	return e
}

func (e *Event) WithError(err error) *Event {
	e.err = err
	return e
}

func NewEventLoop(logger logger.Logger) *EventLoop {
	return &EventLoop{
		queue:       make(chan *Event),
		queueClosed: false,
		queueLock:   &sync.Mutex{},
		quitch:      make(chan struct{}, 1),
		logger:      logger,
		processor:   newProcessor(),
	}
}

func (ev *EventLoop) AddEvent(event *Event) {
	if ev.queueClosed {
		return
	}
	ev.queue <- event
}

func (ev *EventLoop) CloseLoop() {
	if !ev.queueClosed {
		ev.queueLock.Lock()
		defer ev.queueLock.Unlock()

		if !ev.queueClosed {
			ev.queueClosed = true
			close(ev.queue)
			ev.logger.Info("event loop closed")
		}
	}

}

func (ev *EventLoop) Start() {
	go func() {
		ev.logger.Info("event loop started")
		for event := range ev.queue {
			ev.processor.Process(event)
		}
	}()
}
