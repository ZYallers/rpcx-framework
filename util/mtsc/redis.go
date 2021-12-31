package mtsc

import (
	"errors"
	"github.com/ZYallers/rpcx-framework/helper"
	"github.com/ZYallers/zgin/libraries/json"
	"github.com/ZYallers/zgin/libraries/mvcs"
	"github.com/go-redis/redis"
	"math/rand"
	"strings"
	"time"
)

type Redis struct {
	mvcs.Redis
	Client func() *redis.Client
}

const hashAllFieldKey = "all"

// 数据不存在情况下，为防止缓存雪崩，随机返回一个30到60秒的有效时间
func (r *Redis) NoDataExpiration() time.Duration {
	// 将时间戳设置成种子数
	rand.Seed(time.Now().UnixNano())
	return time.Duration(30+rand.Intn(30)) * time.Second
}

// 从String类型的缓存中读取数据，如没则重新调用指定方法重新从数据库中读取并写入缓存
func (r *Redis) CacheWithString(key string, output interface{}, expiration time.Duration, fn func() (interface{}, bool)) error {
	if val := r.Client().Get(key).Val(); val != "" {
		return json.Unmarshal(helper.String2Bytes(val), &output)
	}

	var (
		isNull bool
		data   interface{}
	)

	if data, isNull = fn(); isNull {
		expiration = r.NoDataExpiration()
	}

	var value string
	bte, err := json.Marshal(data)
	if err != nil {
		value = "null"
	} else {
		value = helper.Bytes2String(bte)
		_ = json.Unmarshal(bte, &output)
	}
	return r.Client().Set(key, value, expiration).Err()
}

//  DeleteCache 根据key删除对应缓存
//  @receiver r *Redis
//  @author Cloud|2021-12-07 13:56:08
//  @param key ...string ...
//  @return int64 ...
//  @return error ...
func (r *Redis) DeleteCache(key ...string) (int64, error) {
	return r.Client().Del(key...).Result()
}

//  HashGetAll ...
//  @receiver r *Redis
//  @author Cloud|2021-12-15 10:10:46
//  @param key string ...
//  @return result []interface{} ...
func (r *Redis) HashGetAll(key string) (result []interface{}) {
	all := r.Client().HGet(key, hashAllFieldKey).Val()
	if all == "" {
		return
	}
	keys := helper.RemoveDuplicateWithString(strings.Split(all, ","))
	if len(keys) == 0 {
		return
	}
	result = r.Client().HMGet(key, keys...).Val()
	return
}

//  HashMultiSet ...
//  @receiver r *Redis
//  @author Cloud|2021-12-15 10:10:49
//  @param key string ...
//  @param data map[string]interface{} ...
//  @return error ...
func (r *Redis) HashMultiSet(key string, data map[string]interface{}) error {
	fields := make([]string, 0)
	fieldValues := make(map[string]interface{}, 0)
	for k, v := range data {
		if k == "" || v == nil {
			continue
		}
		if b, err := json.Marshal(v); err == nil {
			fieldValues[k] = helper.Bytes2String(b)
			fields = append(fields, k)
		}
	}

	if len(fields) == 0 {
		return errors.New("the data that can be saved is empty")
	}

	if val := r.Client().HGet(key, hashAllFieldKey).Val(); val != "" {
		fields = append(fields, strings.Split(val, ",")...)
	}

	var allFieldValue string
	if len(fields) > 0 {
		allFieldValue = strings.Join(helper.RemoveDuplicateWithString(fields), ",")
	}
	fieldValues[hashAllFieldKey] = allFieldValue
	return r.Client().HMSet(key, fieldValues).Err()
}

//  HashMultiDelete ...
//  @receiver r *Redis
//  @author Cloud|2021-12-15 10:49:15
//  @param key string ...
//  @param fields ...string ...
//  @return error ...
func (r *Redis) HashMultiDelete(key string, fields ...string) error {
	newFields := make([]string, 0)
	if val := r.Client().HGet(key, hashAllFieldKey).Val(); val != "" {
		newFields = append(newFields, strings.Split(val, ",")...)
	}
	if len(newFields) > 0 {
		for _, field := range fields {
			newFields = helper.RemoveWithString(newFields, field)
		}
	}

	var allFieldValue string
	if len(newFields) > 0 {
		allFieldValue = strings.Join(newFields, ",")
	}

	pl := r.Client().Pipeline()
	pl.HDel(key, fields...)
	pl.HSet(key, hashAllFieldKey, allFieldValue)
	_, err := pl.Exec()
	return err
}
