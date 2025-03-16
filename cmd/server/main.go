package main

import (
	"goredis/common/config"
	"goredis/common/logger"
	"goredis/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
)

func main() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	configPath := viper.GetString("GOREDIS_CONFIG_PATH")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalln("failed to load config due to ", err)
		os.Exit(0)
	}

	zapLogger, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalln("failed to build logger due to ", err)
		os.Exit(0)
	}

	zapLogger.Info("go redis server is ready to accept connections")

	srv := server.NewTcpServer(server.WithConfig(cfg), server.WithLogger(zapLogger))

	go func() {
		if err := srv.Start(); err != nil {
			zapLogger.Error(err)
		}
	}()

	<-quit

	zapLogger.Info("shutdown signal received closing all connections")
	srv.Stop()
}
