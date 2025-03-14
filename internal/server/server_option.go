package server

import (
	"goredis/common/config"
	"goredis/common/logger"
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
