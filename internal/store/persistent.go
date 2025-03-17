package store

import (
	"goredis/common/config"
	"goredis/common/logger"
	"time"
)

type (
	Persistent struct {
		persistentOptions *config.PersistentOptions
		logger            logger.Logger
		exit              chan struct{}
		intervalTicker    *time.Ticker
		kvStore           *KeyValueStore
	}

	persistentOpt func(*Persistent)
)

const (
	defaultSnapInterval int = 10
)

func WithKv(store *KeyValueStore) persistentOpt {
	return func(p *Persistent) {
		p.kvStore = store
	}
}

func WithLogger(logger logger.Logger) persistentOpt {
	return func(per *Persistent) {
		per.logger = logger
	}
}

func WithPersistenOpts(opt *config.PersistentOptions) persistentOpt {
	return func(p *Persistent) {
		p.persistentOptions = opt
		if p.persistentOptions.Interval == 0 {
			p.persistentOptions.Interval = defaultSnapInterval
		}

		switch p.persistentOptions.Unit {
		case "h":
			p.intervalTicker = time.NewTicker(time.Hour * time.Duration(p.persistentOptions.Interval))
		case "m":
			p.intervalTicker = time.NewTicker(time.Minute * time.Duration(p.persistentOptions.Interval))
		default:
			p.intervalTicker = time.NewTicker(time.Second * time.Duration(p.persistentOptions.Interval))
		}
	}
}

func NewPeristent(opts ...persistentOpt) *Persistent {

	persistent := &Persistent{
		exit: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(persistent)
	}
	return persistent
}

func (per *Persistent) PersistData() {

	go func() {
		for {
			select {
			case <-per.exit:
				return
			case <-per.intervalTicker.C:
				per.logger.Info("taking kv snapshot")

			}
		}
	}()
}

func (per *Persistent) Close() {
	close(per.exit)
}
