package server

import (
	"goredis/common/config"
	"goredis/common/logger"
)

type (
	serverOptions struct {
		config *config.Config
		logger logger.Logger
	}
)
