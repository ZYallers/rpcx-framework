package types

import (
	"time"

	"github.com/ZYallers/golib/utils/redis"
)

type Redis struct {
	redis.Redis
}

type RedisVars struct {
	CommonExpiration time.Duration
	TTL              TTLType
}

type TTLType struct {
	Forever, NotExist float64
}

type RedisKey struct {
	String, Hash, Set, ZSet, List, Geo, Hyper, Bitmap map[string]string
}
