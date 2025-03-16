package store

import (
	"errors"
	gerrors "goredis/errors"
	"goredis/internal/constants"
	"goredis/internal/request"
	"goredis/internal/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	Value struct {
		value    string
		values   []string
		datatype string
	}

	KeyValueStore struct {
		store      map[string]*Value
		storeLock  *sync.Mutex
		ticker     *time.Ticker
		ttlTracker map[string]int64
		ttlDone    chan bool
	}
)

var (
	allowedDataTypes = map[string]struct{}{
		constants.STR:  {},
		constants.INT:  {},
		constants.LIST: {},
	}
)

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		store:      make(map[string]*Value, 1000),
		storeLock:  &sync.Mutex{},
		ticker:     time.NewTicker(time.Second * 1),
		ttlTracker: make(map[string]int64),
		ttlDone:    make(chan bool),
	}
}

func (kv *KeyValueStore) InitKvStore() {
	kv.clearExpiredKeys()
}

func (kv *KeyValueStore) Close() {
	kv.ticker.Stop()
	kv.ttlDone <- true
}

func (kv *KeyValueStore) Add(req request.Request) (*string, error) {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	if req.Key == nil || req.Ttl == nil || req.Value == nil || req.Datatype == nil {
		return nil, gerrors.ErrRequiredParamsMissingSet
	}

	dt := strings.ToUpper(*req.Datatype)
	if _, ok := allowedDataTypes[dt]; !ok {
		return nil, gerrors.ErrInvalidDatatypes
	}

	value := &Value{}
	success := "1"
	if strings.EqualFold(*req.Datatype, constants.LIST) {
		value.values = strings.Split(*req.Value, ",")
		success = strconv.Itoa(len(value.values))
	} else {
		if strings.EqualFold(*req.Datatype, constants.INT) {
			_, err := strconv.Atoi(*req.Value)
			if err != nil {
				return nil, errors.New("not an integer for integer datatype")
			}
		}

		value.value = *req.Value
	}

	value.datatype = dt

	kv.store[*req.Key] = value

	err := kv.setExpireTime(req.Key, req.Ttl)
	if err != nil {
		return nil, err
	}
	return &success, nil
}

func (kv *KeyValueStore) Get(req request.Request) (*string, error) {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	kv.deleteKeyAtTime(*req.Key)

	if value, ok := kv.store[*req.Key]; ok {
		if value.datatype == constants.LIST {
			res := strings.Join(value.values, ",")
			return &res, nil
		}
		return &value.value, nil
	}

	return nil, errors.New("key does not exists")
}

func (kv *KeyValueStore) Delete(req request.Request) (*string, error) {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	kv.deleteKeyAtTime(*req.Key)

	if _, ok := kv.store[*req.Key]; ok {
		kv.deleteKeys(*req.Key)
	}

	res := "1"
	return &res, nil
}

func (kv *KeyValueStore) SetExpiration(req request.Request) (*string, error) {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	kv.deleteKeyAtTime(*req.Key)

	if _, ok := kv.store[*req.Key]; ok {
		err := kv.setExpireTime(req.Key, req.Ttl)
		if err != nil {
			return nil, err
		}

		res := "1"
		return &res, nil
	}

	return nil, errors.New("key does not exists")
}

func (kv *KeyValueStore) Push(req request.Request) (*string, error) {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	kv.deleteKeyAtTime(*req.Key)

	if req.Value == nil {
		return nil, errors.New("expected value to push")
	}

	if _, ok := kv.store[*req.Key]; !ok {
		return nil, errors.New("key does not exists")
	}

	storedValue := kv.store[*req.Key]
	if storedValue.datatype != constants.LIST {
		return nil, errors.New("invalid operation push allowd on list only")
	}

	storedValue.values = append(storedValue.values, *req.Value)
	kv.store[*req.Value] = storedValue
	res := "1"
	return &res, nil
}

func (kv *KeyValueStore) Pop(req request.Request) (*string, error) {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	if req.Value == nil {
		return nil, errors.New("expected no of elements to pop")
	}

	kv.deleteKeyAtTime(*req.Key)

	if _, ok := kv.store[*req.Key]; !ok {
		return nil, errors.New("key does not exists")
	}

	storedValue := kv.store[*req.Key]
	if storedValue.datatype != constants.LIST {
		return nil, errors.New("invalid operation pop allowed on list only")
	}

	values := strings.Split(*req.Value, " ")
	dir, ele := strings.TrimSpace(values[0]), strings.TrimSpace(values[1])
	if utils.IsEmpty(dir) || utils.IsEmpty(ele) {
		return nil, errors.New("no of elements or direction to pop")
	}
	noOfEle, err := strconv.Atoi(ele)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(dir, "L") {
		for noOfEle > 0 && len(storedValue.values) > 0 {
			storedValue.values = storedValue.values[1:]
		}
	} else {
		for noOfEle > 0 && len(storedValue.values) > 0 {
			storedValue.values = storedValue.values[:len(storedValue.values)-1]
		}
	}

	res := "1"
	return &res, nil
}

func (kv *KeyValueStore) Incr(req request.Request) (*string, error) {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	kv.deleteKeyAtTime(*req.Key)

	if _, ok := kv.store[*req.Key]; !ok {
		return nil, errors.New("key does not exists")
	}

	storedValue := kv.store[*req.Key]
	if storedValue.datatype != constants.INT {
		return nil, errors.New("invalid operation incr allowed on int only")
	}

	val, _ := strconv.Atoi(storedValue.value)
	val += 1
	storedValue.value = strconv.Itoa(val)
	kv.store[*req.Key] = storedValue
	res := "1"
	return &res, nil
}

func (kv *KeyValueStore) Decr(req request.Request) (*string, error) {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	kv.deleteKeyAtTime(*req.Key)

	if _, ok := kv.store[*req.Key]; !ok {
		return nil, errors.New("key does not exists")
	}

	storedValue := kv.store[*req.Key]
	if storedValue.datatype != constants.INT {
		return nil, errors.New("invalid operation decr allowed on int only")
	}

	val, _ := strconv.Atoi(storedValue.value)
	val -= 1
	storedValue.value = strconv.Itoa(val)
	kv.store[*req.Key] = storedValue
	res := "1"

	return &res, nil
}

func (kv *KeyValueStore) setExpireTime(key, ttl *string) error {

	if ttl == nil {
		return nil
	}

	duration, err := strconv.Atoi(*ttl)
	if err != nil {
		return err
	}

	kv.ttlTracker[*key] = time.Now().Unix() + int64(duration)
	return nil

}

func (kv *KeyValueStore) clearExpiredKeys() {

	go func() {

		for {
			select {
			case <-kv.ttlDone:
				return
			case t := <-kv.ticker.C:
				for key, value := range kv.ttlTracker {
					if value < t.Unix() {
						kv.deleteKeys(key)
					}
				}
			}
		}
	}()
}

func (kv *KeyValueStore) deleteKeyAtTime(key string) {

	if ttl, ttlOk := kv.ttlTracker[key]; ttlOk && ttl < time.Now().Unix() {
		kv.deleteKeys(key)
	}
}

func (kv *KeyValueStore) deleteKeys(keys ...string) {
	for _, key := range keys {
		delete(kv.store, key)
		delete(kv.ttlTracker, key)
	}
}
