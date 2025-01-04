package localCache

import (
	"context"
	"github.com/04Akaps/elasticSearch.git/common/json"
	"github.com/allegro/bigcache/v3"
	"log"
	"strconv"
	"strings"
	"time"
)

type LocalCache struct {
	cache     *bigcache.BigCache
	marshal   func(v interface{}) ([]byte, error)
	unMarshal func(data []byte, v interface{}) error
}

func NewLocalCache(ttl time.Duration) LocalCache {
	log.Println("local cache connect started")

	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(ttl))

	if err != nil {
		log.Panic("Failed to connect cache", "cerr", err)
	}

	log.Println("Success to connect local cache")

	return LocalCache{
		cache:     cache,
		marshal:   json.JsonHandler.Marshal,
		unMarshal: json.JsonHandler.Unmarshal,
	}
}

func (l *LocalCache) Set(key string, value interface{}) error {
	v, err := l.marshal(value)

	if err != nil {
		return err
	}

	return l.cache.Set(key, v)
}

func (l *LocalCache) Get(key string, buffer interface{}) error {
	result, err := l.cache.Get(key)

	if err != nil {
		return err
	}

	err = l.unMarshal(result, buffer)

	if err != nil {
		return err
	}

	return nil
}

func GetLCacheKey(v ...interface{}) string {
	var builder strings.Builder

	for i, val := range v {
		if i > 0 {
			builder.WriteString(":") // 값 사이에 하이픈 추가
		}

		switch value := val.(type) {
		case int64:
			builder.WriteString(strconv.FormatInt(value, 10)) // int64를 string으로 변환하여 추가
		case string:
			builder.WriteString(value) // string은 그대로 추가
		}
	}

	return builder.String()
}
