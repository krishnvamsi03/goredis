package store

import (
	"encoding/json"
	"goredis/common/config"
	"goredis/common/logger"
	"goredis/proto/persistent"
	"os"
	"time"

	"google.golang.org/protobuf/proto"
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
	defaultRetry        int = 3
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
				per.intervalTicker.Stop()
				per.logger.Info("persitent data snapshot closed")
				return
			case <-per.intervalTicker.C:
				per.logger.Info("taking kv store snapshot")

				for i := 0; i < defaultRetry; i++ {
					err := per.kvStore.Persist(per.persistentOptions.Path)
					if err != nil {
						per.logger.Error(err)
						continue
					}
					per.logger.Info("snapshot succesfull")
					break
				}
			}
		}
	}()
}

func (per *Persistent) LoadData() error {

	byteData, err := os.ReadFile(per.persistentOptions.Path)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return nil
		}
		per.logger.Error(err)
		return err
	}

	var kvData persistent.PersistentStore
	err = proto.Unmarshal(byteData, &kvData)
	if err != nil {
		per.logger.Error(err)
		return err
	}

	store := kvData.GetKv().GetStore()
	storeBytes, err := json.Marshal(store)
	if err != nil {
		per.logger.Error(err)
		return err
	}

	var kvStoreDeserial map[string]*Value
	err = json.Unmarshal(storeBytes, &kvStoreDeserial)
	if err != nil {
		per.logger.Error(err)
		return err
	}

	ttlTracker := kvData.GetKv().GetTtlTracker()
	if len(kvStoreDeserial) > 0 || len(ttlTracker) > 0 {
		per.kvStore.LoadFromSnapshot(kvStoreDeserial, ttlTracker)
		per.logger.Info("successfully loaded data from disk")
	}
	return nil
}

func (per *Persistent) Close() {
	close(per.exit)
}
