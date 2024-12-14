package cache

import (
	"github.com/04Akaps/elasticSearch.git/cache/localCache"
	"github.com/04Akaps/elasticSearch.git/cache/redis"
	"github.com/04Akaps/elasticSearch.git/config"
	"time"
)

type CacheManager struct {
	redis      redis.Redis
	localCache localCache.LocalCache
}

func NewCacheManager(cfg config.Config) CacheManager {
	return CacheManager{
		redis:      redis.NewRedis(cfg),
		localCache: localCache.NewLocalCache(time.Duration(cfg.Cache.Local.LocalCacheTTL)),
	}
}

func (c CacheManager) GetValueWithTTL(key string, buffer interface{}) (time.Duration, error) {
	return c.redis.GetValueWithTTL(key, buffer)
}

func (c CacheManager) SetNX(key string, v interface{}, ttl time.Duration) error {
	return c.redis.SetNX(key, v, ttl)
}
