package event_processor

import (
	"goredis/common/logger"
	"goredis/internal/request"
	"goredis/internal/store"
	"net"
	"sync"
	"time"
)

type (
	Event struct {
		conn net.Conn
		cmd  *request.Request
		err  error
	}

	EventOption func(*Event)

	EventLoop struct {
		queue         chan *Event
		queueClosed   bool
		queueLock     *sync.Mutex
		quitch        chan struct{}
		logger        logger.Logger
		processor     *processor
		persistent    *store.Persistent
		keyvalueStore *store.KeyValueStore
	}

	EventLoopOption func(*EventLoop)
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

func NewEventLoop(opts ...EventLoopOption) *EventLoop {

	ev := &EventLoop{
		queue:       make(chan *Event),
		queueClosed: false,
		queueLock:   &sync.Mutex{},
		quitch:      make(chan struct{}, 1),
	}

	for _, opt := range opts {
		opt(ev)
	}
	return ev
}

func WithProcessor(processor *processor) EventLoopOption {
	return func(ev *EventLoop) {
		ev.processor = processor
	}
}

func WithLogger(logger logger.Logger) EventLoopOption {
	return func(el *EventLoop) {
		el.logger = logger
	}
}

func WithPersistent(persist *store.Persistent) EventLoopOption {
	return func(el *EventLoop) {
		el.persistent = persist
	}
}

func WithKeyValueStore(kv *store.KeyValueStore) EventLoopOption {
	return func(el *EventLoop) {
		el.keyvalueStore = kv
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
		ev.persistent.Close()
		ev.keyvalueStore.Close()

		if !ev.queueClosed {
			ev.queueClosed = true
			close(ev.queue)
			time.Sleep(3 * time.Second)
			ev.logger.Info("event loop closed")
		}
	}

}

func (ev *EventLoop) Start() {

	if err := ev.persistent.LoadData(); err != nil {
		panic(err)
	}

	ev.persistent.PersistData()
	ev.keyvalueStore.Start()

	go func() {

		ev.logger.Info("event loop started")
		for event := range ev.queue {
			ev.processor.Process(event)
		}
	}()
}
