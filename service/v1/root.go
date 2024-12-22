package v1

import (
	"github.com/04Akaps/elasticSearch.git/cache"
	"github.com/04Akaps/elasticSearch.git/config"
)

type V1 struct {
	cfg   config.Config
	cache cache.CacheManager
}

func NewV1(
	cfg config.Config,
	cache cache.CacheManager,
) V1 {
	return V1{cfg, cache}
}
