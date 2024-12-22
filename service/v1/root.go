package v1

import (
	"github.com/04Akaps/elasticSearch.git/cache"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/repository/elasticSearch"
)

type V1 struct {
	cfg           config.Config
	cache         cache.CacheManager
	elasticSearch elasticSearch.ElasticSearch
}

func NewV1(
	cfg config.Config,
	cache cache.CacheManager,
	elasticSearch elasticSearch.ElasticSearch,
) V1 {
	return V1{cfg, cache, elasticSearch}
}
