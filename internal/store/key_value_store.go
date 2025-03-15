package store

import (
	"errors"
	"goredis/internal/request"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	Value struct {
		value    string
		datatype string
	}

	KeyValueStore struct {
		store      map[string]*Value
		storeLock  *sync.Mutex
		keyTracker map[int][]string
	}
)

var (
	allowedUnits = map[string]struct{}{
		"s": {},
		"m": {},
		"h": {},
	}

	INF         = (1 << 63) - 1
	defaultUnit = "s"
)

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		store:      make(map[string]*Value, 1000),
		keyTracker: make(map[int][]string),
		storeLock:  &sync.Mutex{},
	}
}

func (kv *KeyValueStore) Add(req request.Request) (*string, error) {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	kv.store[*req.Key] = &Value{
		value:    *req.Value,
		datatype: *req.Datatype,
	}
	success := "OK 1"
	err := kv.setExpireTime(*req.Key, *req.Expr)
	if err != nil {
		return nil, err
	}
	return &success, nil
}

func (kv *KeyValueStore) Get(req request.Request) (*string, error) {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	if value, ok := kv.store[*req.Key]; ok {
		return &value.value, nil
	}

	return nil, errors.New("key does not exists")
}

func (kv *KeyValueStore) setExpireTime(key, expr string) error {

	if len(expr) == 0 {
		kv.keyTracker[INF] = append(kv.keyTracker[INF], key)
		return nil
	}

	unit := string(expr[len(expr)-1])
	if _, ok := allowedUnits[unit]; !ok {
		unit = defaultUnit
	}

	duration := expr[:len(expr)-1]
	if unit >= "0" && unit <= "9" {
		duration = expr
	}

	ti, err := strconv.Atoi(duration)
	if err != nil {
		return err
	}

	currentTime := int(time.Now().UnixMilli())
	switch strings.ToLower(unit) {
	case "s":
		currentTime += (1000 * ti)
	case "m":
		currentTime += (1000 * 60 * ti)
	case "h":
		currentTime += (1000 * 60 * 60 * ti)
	}

	kv.keyTracker[currentTime] = append(kv.keyTracker[currentTime], key)
	return nil

}
