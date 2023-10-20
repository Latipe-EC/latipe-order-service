package cache

import (
	"fmt"
	"github.com/google/wire"
	"order-rest-api/config"
	"order-rest-api/pkg/cache/redis"
	"time"
)

var Set = wire.NewSet(
	NewCacheEngine,
)

func NewCacheEngine(config *config.Config) (*redis.CacheEngine, error) {
	cfg := redis.RedisConfig{
		Address:               fmt.Sprintf("%v:%v", config.Cache.Redis.Address, config.Cache.Redis.Port),
		DB:                    config.Cache.Redis.DB,
		Password:              config.Cache.Redis.Password,
		ContextTimeoutEnabled: true,
		PoolSize:              5,
		PoolTimeout:           5,
		DialTimeout:           5,
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          5 * time.Second,
		ConnectTimeout:        5 * time.Second,
	}
	client, err := redis.NewCacheEngine(cfg)
	if err != nil {
		return nil, err
	}
	return client, err
}
