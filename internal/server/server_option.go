package server

import (
	"goredis/common/config"
	"goredis/common/logger"
	"goredis/internal/event_processor"
	"goredis/internal/protocol"
)

type (
	ServerOption interface {
		apply(*serverOptions)
	}

	applyFunc func(*serverOptions)
)

func (f applyFunc) apply(s *serverOptions) { f(s) }

func WithConfig(cfg *config.Config) ServerOption {
	return applyFunc(func(so *serverOptions) {
		so.config = cfg
	})
}

func WithLogger(logger logger.Logger) ServerOption {
	return applyFunc(func(so *serverOptions) {
		so.logger = logger
	})
}

func WithEventLoop(ev *event_processor.EventLoop) ServerOption {
	return applyFunc(func(so *serverOptions) {
		so.eventLoop = ev
	})
}

func WithParser(parser protocol.Parser) ServerOption {
	return applyFunc(func(so *serverOptions) {
		so.parser = parser
	})
}
