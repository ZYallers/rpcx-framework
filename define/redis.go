package define

import "time"

type Redis struct {
	CommonExpiration time.Duration
	TTL              TTLType
}

type TTLType struct {
	Forever, NotExist float64
}

type RedisKey struct {
	String, Hash, ZSet, List map[string]string
}
