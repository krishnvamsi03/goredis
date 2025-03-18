package server

import (
	"goredis/common/config"
	"goredis/common/logger"
	"goredis/internal/event_processor"
	"goredis/internal/protocol"
)

type (
	serverOptions struct {
		config    *config.Config
		logger    logger.Logger
		eventLoop *event_processor.EventLoop
		parser    protocol.Parser
	}
)
