package main

import (
	"goredis/common/config"
	"goredis/common/logger"
	"goredis/internal/command"
	"goredis/internal/event_processor"
	"goredis/internal/protocol"
	"goredis/internal/server"
	"goredis/internal/store"
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

	kvStore := store.NewKeyValueStore(zapLogger)
	commanManger := command.NewCommandManager(kvStore)
	persist := store.NewPeristent(
		store.WithLogger(zapLogger),
		store.WithPersistenOpts(&cfg.PersistentOptions),
		store.WithKv(kvStore),
	)

	processor := event_processor.NewProcessor(commanManger)

	el := event_processor.NewEventLoop(
		event_processor.WithProcessor(processor),
		event_processor.WithLogger(zapLogger),
		event_processor.WithPersistent(persist),
		event_processor.WithKeyValueStore(kvStore),
	)

	srv := server.NewTcpServer(server.WithConfig(cfg),
		server.WithLogger(zapLogger),
		server.WithParser(protocol.NewGrespParser(zapLogger)),
		server.WithEventLoop(el),
	)

	zapLogger.Info("go redis server is ready to accept connections")

	go func() {
		if err := srv.Start(); err != nil {
			zapLogger.Error(err)
		}
	}()

	<-quit

	zapLogger.Info("shutdown signal received closing all connections")
	srv.Stop()
}
