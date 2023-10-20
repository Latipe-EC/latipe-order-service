package redis

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisConfig struct {
	Address               string
	DB                    int
	Password              string
	ContextTimeoutEnabled bool
	PoolSize              int
	PoolTimeout           time.Duration
	DialTimeout           time.Duration
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	ConnectTimeout        time.Duration
}

func NewCacheEngine(cfg RedisConfig) (*CacheEngine, error) {
	client := redis.NewClient(&redis.Options{
		DB:   cfg.DB,
		Addr: cfg.Address,
		//PoolSize: cfg.PoolSize,
		//Password:     cfg.Password,
		//DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	cacheEngine := &CacheEngine{
		client: client,
		ctx:    context.Background(),
	}
	if err := cacheEngine.Ping(); err != nil {
		return nil, err
	}

	return cacheEngine, nil
}

type CacheEngine struct {
	client *redis.Client
	ctx    context.Context
}

// WithContext for operate
func (c *CacheEngine) WithContext(ctx context.Context) *CacheEngine {
	cp := *c
	cp.ctx = ctx
	return &cp
}

// Get gets the value for the given key.
func (c *CacheEngine) Get(key string) ([]byte, error) {
	result := c.client.Get(c.ctx, key)
	val, err := result.Bytes()
	if redis.Nil == err {
		return val, fiber.ErrNotFound
	}
	return val, err
}

// Get gets the value for the given key.
func (c *CacheEngine) Expire(key string) error {
	result := c.client.Expire(c.ctx, key, 10*time.Minute)
	return result.Err()
}

// Set stores the given value for the given key along with a
func (c *CacheEngine) Set(key string, val []byte) error {
	ttl := 1 * time.Hour
	result := c.client.Set(c.ctx, key, val, ttl)
	return result.Err()
}

// Delete deletes the value for the given key.
func (c *CacheEngine) Delete(key string) error {
	result := c.client.Del(c.ctx, key)
	return result.Err()
}

// Reset resets the storage and delete all keys.
func (c *CacheEngine) Reset() error {
	result := c.client.FlushAll(c.ctx)
	return result.Err()
}

// Close closes the storage and will stop any running garbage
func (c *CacheEngine) Close() error {
	return c.client.Close()
}

// Ping check connection
func (c *CacheEngine) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.client.Ping(ctx).Err()
	if err != nil {
		return err
	}

	return nil
}
