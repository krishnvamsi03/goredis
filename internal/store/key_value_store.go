package store

import (
	"encoding/json"
	"fmt"
	"goredis/common/logger"
	"goredis/internal/constants"
	"goredis/internal/request"
	"goredis/internal/response"
	statuscodes "goredis/internal/status_codes"
	"goredis/internal/utils"
	"goredis/proto/persistent"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

type (
	Value struct {
		Value    string   `json:"value"`
		Values   []string `json:"Values"`
		Datatype string   `json:"datatype"`
	}

	KeyValueStore struct {
		store      map[string]*Value
		storeLock  *sync.Mutex
		ticker     *time.Ticker
		ttlTracker map[string]int64
		ttlDone    chan bool
		logger     logger.Logger
	}

	KeyValueStoreOpt func(*KeyValueStore)
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

func NewKeyValueStore(logger logger.Logger) *KeyValueStore {
	return &KeyValueStore{
		store:      make(map[string]*Value, 1000),
		storeLock:  &sync.Mutex{},
		ticker:     time.NewTicker(time.Second * 1),
		ttlTracker: make(map[string]int64),
		ttlDone:    make(chan bool),
		logger:     logger,
	}
}

func (kv *KeyValueStore) Start() {
	kv.clearExpiredKeys()
}

func (kv *KeyValueStore) Close() {
	kv.ticker.Stop()
	kv.ttlDone <- true
}

func (kv *KeyValueStore) Ping(req request.Request) *response.Response {
	return response.NewResponse().
		WithCode(statuscodes.SUCCESS).
		WithOk(true).
		WithRes("PONG")
}

func (kv *KeyValueStore) Add(req request.Request) *response.Response {
	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	response := response.NewResponse()

	if req.Key == nil || req.Value == nil || req.Datatype == nil {
		return response.WithCode(statuscodes.REQUIRED_INP_MISSING).
			WithOk(false).
			WithRes("either key, Datatype or value is missing for setting value")
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
		value.Values = strings.Split(*req.Value, ",")
	case constants.INT:
		_, err := strconv.Atoi(*req.Value)
		if err != nil {
			return response.WithCode(statuscodes.INVALID_VALUE_FOR_DATATYPE).
				WithOk(false).
				WithRes("invalid value provided for given data type")
		}
	default:
		value.Value = *req.Value
	}

	value.Datatype = dt
	kv.store[*req.Key] = value

	if req.Ttl != nil {
		err := kv.setExpireTime(req.Key, req.Ttl)
		if err != nil {
			return response.WithCode(statuscodes.INVALID_INP).
				WithOk(false).
				WithRes("failed to set ttl")
		}
	}
	return response.WithCode(statuscodes.SUCCESS).
		WithOk(true).
		WithRes("1")
}

func (kv *KeyValueStore) GetKey(req request.Request) *response.Response {

	if req.Key == nil {
		return response.NewResponse().
			WithCode(statuscodes.REQUIRED_INP_MISSING).
			WithOk(false).
			WithRes("key is missing")
	}

	if *req.Key == "*" {
		keys := ""
		for key := range kv.store {
			keys += fmt.Sprintf("%s\n", key)
		}

		keys = strings.TrimSpace(keys)
		return response.NewResponse().WithCode(statuscodes.SUCCESS).WithOk(true).WithRes(keys)
	}

	validPattPattern := `^\*?[a-zA-Z0-9_]+(:[a-zA-Z0-9_]*)?\*?$`

	re := regexp.MustCompile(validPattPattern)
	if !re.MatchString(*req.Key) {
		return response.NewResponse().
			WithCode(statuscodes.INVALID_INP).
			WithOk(false).
			WithRes("only * wildcard is allowed")
	}

	res := ""
	keyRe := regexp.MustCompile(fmt.Sprintf("^%s", *req.Key))
	for key := range kv.store {
		if keyRe.MatchString(key) {
			res = fmt.Sprintf("%s\n", key)
		}
	}

	res = strings.TrimSpace(res)
	return response.NewResponse().
		WithCode(statuscodes.SUCCESS).
		WithOk(true).
		WithRes(res)
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
	switch value.Datatype {
	case constants.LIST:
		res = strings.Join(value.Values, ",")
	default:
		res = value.Value
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
	if storedValue.Datatype != constants.LIST {
		return response.WithCode(statuscodes.OPERATION_NOT_ALLOWED_FOR_DATATYPE).
			WithOk(false).
			WithRes("push cannot be performed on Non List type")
	}

	Values := strings.Split(*req.Value, ",")
	for _, v := range Values {
		storedValue.Values = append(storedValue.Values, strings.TrimSpace(v))
	}
	kv.store[*req.Value] = storedValue

	return response.WithCode(statuscodes.SUCCESS).
		WithOk(true).
		WithRes(fmt.Sprintf("%d", len(Values)))
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
	if storedValue.Datatype != constants.LIST {
		return response.WithCode(statuscodes.OPERATION_NOT_ALLOWED_FOR_DATATYPE).
			WithOk(false).
			WithRes("pop cannot be performed on Non List type")
	}

	Values := strings.Split(*req.Value, " ")
	dir, ele := strings.TrimSpace(Values[0]), strings.TrimSpace(Values[1])
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
		for noOfEle > 0 && len(storedValue.Values) > 0 {
			storedValue.Values = storedValue.Values[1:]
			noOfEle--
			cnt++
		}
	} else {
		for noOfEle > 0 && len(storedValue.Values) > 0 {
			storedValue.Values = storedValue.Values[:len(storedValue.Values)-1]
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
	if storedValue.Datatype != constants.INT {
		return response.WithCode(statuscodes.OPERATION_NOT_ALLOWED_FOR_DATATYPE).
			WithOk(false).
			WithRes("incr can apply on int types only")
	}

	val, _ := strconv.Atoi(storedValue.Value)
	val += 1
	storedValue.Value = strconv.Itoa(val)
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
	if storedValue.Datatype != constants.INT {
		return response.WithCode(statuscodes.OPERATION_NOT_ALLOWED_FOR_DATATYPE).
			WithOk(false).
			WithRes("decr can apply on int types only")
	}

	val, _ := strconv.Atoi(storedValue.Value)
	val -= 1
	storedValue.Value = strconv.Itoa(val)
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
				kv.ticker.Stop()
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

func (kv *KeyValueStore) Persist(path string) error {

	kv.storeLock.Lock()
	defer kv.storeLock.Unlock()

	byteData, err := json.Marshal(kv.store)
	if err != nil {
		kv.logger.Error(err)
		return err
	}

	var persistentStore map[string]*persistent.Value
	err = json.Unmarshal(byteData, &persistentStore)
	if err != nil {
		kv.logger.Error(err)
		return err
	}

	store := &persistent.PersistentStore{
		CreatedAtUnix: time.Now().Unix(),
		CreatedAt:     time.Now().Format(time.RFC3339),
		Kv: &persistent.KeyValueStore{
			Store:      persistentStore,
			TtlTracker: kv.ttlTracker,
		},
	}

	serialData, err := proto.Marshal(store)
	if err != nil {
		kv.logger.Error(err)
		return err
	}
	err = os.WriteFile(path, serialData, 0664)
	if err != nil {
		kv.logger.Error(err)
		return err
	}
	return nil
}
