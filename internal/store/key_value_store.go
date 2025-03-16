package store

import (
	"fmt"
	"goredis/internal/constants"
	"goredis/internal/request"
	"goredis/internal/response"
	statuscodes "goredis/internal/status_codes"
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

const (
	KEY_DOES_NOT_EXIST_MSG = "key does not exists or expired"
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

func (kv *KeyValueStore) Add(req request.Request) *response.Response {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	response := response.NewResponse()

	if req.Key == nil || req.Ttl == nil || req.Value == nil || req.Datatype == nil {
		return response.WithCode(statuscodes.REQUIRED_INP_MISSING).
			WithOk(false).
			WithRes("either key, ttl, datatype or value is missing for setting value")
	}

	dt := strings.ToUpper(*req.Datatype)
	if _, ok := allowedDataTypes[dt]; !ok {
		return response.WithCode(statuscodes.DATA_TYPE_NOT_ALLOWED).
			WithOk(false).
			WithRes("data type not supported allowed are str, list, int.")
	}

	value := &Value{}
	switch dt {
	case constants.LIST:
		value.values = strings.Split(*req.Value, ",")
	case constants.INT:
		_, err := strconv.Atoi(*req.Value)
		if err != nil {
			return response.WithCode(statuscodes.INVALID_VALUE_FOR_DATATYPE).
				WithOk(false).
				WithRes("invalid value provided for given data type")
		}
	default:
		value.value = *req.Value
	}

	value.datatype = dt
	kv.store[*req.Key] = value

	return response.WithCode(statuscodes.SUCCESS).
		WithOk(true).
		WithRes("1")
}

func (kv *KeyValueStore) Get(req request.Request) *response.Response {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	kv.deleteKeyAtTime(*req.Key)
	response := response.NewResponse()

	if _, ok := kv.store[*req.Key]; !ok {
		return kv.keyDoesNotExistRes()
	}

	value := kv.store[*req.Key]

	res := ""
	switch value.datatype {
	case constants.LIST:
		res = strings.Join(value.values, ",")
	default:
		res = value.value
	}

	return response.WithCode(statuscodes.SUCCESS).
		WithOk(true).
		WithRes(res)
}

func (kv *KeyValueStore) Delete(req request.Request) *response.Response {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	kv.deleteKeyAtTime(*req.Key)

	response := response.NewResponse()

	if _, ok := kv.store[*req.Key]; !ok {
		return response.WithCode(statuscodes.SUCCESS).
			WithOk(true).
			WithRes("0")
	}

	kv.deleteKeys(*req.Key)
	return response.WithCode(statuscodes.SUCCESS).
		WithOk(true).
		WithRes("1")
}

func (kv *KeyValueStore) SetExpiration(req request.Request) *response.Response {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	kv.deleteKeyAtTime(*req.Key)
	_, ok := kv.store[*req.Key]
	if !ok {
		return kv.keyDoesNotExistRes()
	}

	response := response.NewResponse()

	err := kv.setExpireTime(req.Key, req.Ttl)
	if err != nil {
		return response.WithCode(statuscodes.INVALID_INP).
			WithOk(false).
			WithRes("failed to set ttl")
	}

	return response.WithCode(statuscodes.SUCCESS).
		WithOk(true).
		WithRes("1")
}

func (kv *KeyValueStore) Push(req request.Request) *response.Response {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	kv.deleteKeyAtTime(*req.Key)

	response := response.NewResponse()

	if req.Value == nil {
		return response.WithCode(statuscodes.REQUIRED_INP_MISSING).
			WithOk(false).
			WithRes("value is missing for pushing to list")
	}

	if _, ok := kv.store[*req.Key]; !ok {
		return kv.keyDoesNotExistRes()
	}

	storedValue := kv.store[*req.Key]
	if storedValue.datatype != constants.LIST {
		return response.WithCode(statuscodes.OPERATION_NOT_ALLOWED_FOR_DATATYPE).
			WithOk(false).
			WithRes("push cannot be performed on Non List type")
	}

	values := strings.Split(*req.Value, ",")
	for _, v := range values {
		storedValue.values = append(storedValue.values, strings.TrimSpace(v))
	}
	kv.store[*req.Value] = storedValue

	return response.WithCode(statuscodes.SUCCESS).
		WithOk(true).
		WithRes(fmt.Sprintf("%d", len(values)))
}

func (kv *KeyValueStore) Pop(req request.Request) *response.Response {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	kv.deleteKeyAtTime(*req.Key)
	response := response.NewResponse()

	if req.Value == nil {
		return response.WithCode(statuscodes.REQUIRED_INP_MISSING).
			WithOk(false).
			WithRes("value is missing for pop from list")
	}

	if _, ok := kv.store[*req.Key]; !ok {
		return kv.keyDoesNotExistRes()
	}

	storedValue := kv.store[*req.Key]
	if storedValue.datatype != constants.LIST {
		return response.WithCode(statuscodes.OPERATION_NOT_ALLOWED_FOR_DATATYPE).
			WithOk(false).
			WithRes("pop cannot be performed on Non List type")
	}

	values := strings.Split(*req.Value, " ")
	dir, ele := strings.TrimSpace(values[0]), strings.TrimSpace(values[1])
	if utils.IsEmpty(dir) {
		dir = "R"
	}
	if utils.IsEmpty(ele) {
		ele = "1"
	}

	noOfEle, err := strconv.Atoi(ele)
	if err != nil {
		return response.WithCode(statuscodes.INVALID_INP).
			WithOk(false).
			WithRes("failed to convert to int")
	}
	if !strings.EqualFold(dir, "L") && !strings.EqualFold(dir, "R") {
		return response.WithCode(statuscodes.INVALID_INP).
			WithOk(false).
			WithRes("either l or r is needed")
	}

	cnt := 0
	if strings.EqualFold(dir, "L") {
		for noOfEle > 0 && len(storedValue.values) > 0 {
			storedValue.values = storedValue.values[1:]
			noOfEle--
			cnt++
		}
	} else {
		for noOfEle > 0 && len(storedValue.values) > 0 {
			storedValue.values = storedValue.values[:len(storedValue.values)-1]
			noOfEle--
			cnt++
		}
	}

	kv.store[*req.Key] = storedValue
	return response.WithCode(statuscodes.SUCCESS).
		WithOk(true).
		WithRes(fmt.Sprintf("%d", cnt))
}

func (kv *KeyValueStore) Incr(req request.Request) *response.Response {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	kv.deleteKeyAtTime(*req.Key)

	if _, ok := kv.store[*req.Key]; !ok {
		return kv.keyDoesNotExistRes()
	}

	response := response.NewResponse()
	storedValue := kv.store[*req.Key]
	if storedValue.datatype != constants.INT {
		return response.WithCode(statuscodes.OPERATION_NOT_ALLOWED_FOR_DATATYPE).
			WithOk(false).
			WithRes("incr can apply on int types only")
	}

	val, _ := strconv.Atoi(storedValue.value)
	val += 1
	storedValue.value = strconv.Itoa(val)
	kv.store[*req.Key] = storedValue

	return response.WithCode(statuscodes.SUCCESS).
		WithOk(true).
		WithRes("1")
}

func (kv *KeyValueStore) Decr(req request.Request) *response.Response {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	kv.deleteKeyAtTime(*req.Key)

	if _, ok := kv.store[*req.Key]; !ok {
		return kv.keyDoesNotExistRes()
	}

	response := response.NewResponse()

	storedValue := kv.store[*req.Key]
	if storedValue.datatype != constants.INT {
		return response.WithCode(statuscodes.OPERATION_NOT_ALLOWED_FOR_DATATYPE).
			WithOk(false).
			WithRes("decr can apply on int types only")
	}

	val, _ := strconv.Atoi(storedValue.value)
	val -= 1
	storedValue.value = strconv.Itoa(val)
	kv.store[*req.Key] = storedValue
	return response.WithCode(statuscodes.SUCCESS).
		WithOk(true).
		WithRes("1")
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

func (kv *KeyValueStore) keyDoesNotExistRes() *response.Response {
	return response.NewResponse().WithCode(statuscodes.KEY_DOES_NOT_EXISTS).
		WithOk(false).
		WithRes(KEY_DOES_NOT_EXIST_MSG)
}
