package cache

import (
	"context"
	"fmt"
	"time"
)
import "github.com/go-redis/redis/v8"

type RedisServer struct {
	Client         *redis.Client
	CacheTTLSecond int
}

func NewRedisServer(conOptions map[string]interface{}) (rs RedisServer) {
	ro := redis.Options{
		Addr:     fmt.Sprintf("%v", conOptions["addr"]),
		Username: fmt.Sprintf("%v", conOptions["username"]),
		Password: fmt.Sprintf("%v", conOptions["password"]),
		DB:       conOptions["db"].(int),
	}
	ttl, ok := conOptions["ttl"].(int)
	if ok {
		rs.CacheTTLSecond = ttl
	}
	rs.Client = redis.NewClient(&ro)
	return
}

func (rs *RedisServer) GetCache(key string) ([]byte, error) {
	itm := rs.Client.Get(context.Background(), key)
	if itm.Err() == redis.Nil {
		return nil, nil
	}
	dt, err := itm.Bytes()
	if err != nil {
		return nil, err
	}
	return dt, nil
}

func (rs *RedisServer) SetCache(key string, value []byte) error {
	resp := rs.Client.Set(context.Background(), key, value, time.Duration(rs.CacheTTLSecond)*time.Second)
	if resp.Err() != nil {
		return resp.Err()
	}
	return nil
}
