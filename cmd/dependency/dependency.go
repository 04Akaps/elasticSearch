package dependency

import (
	"flag"
	"github.com/04Akaps/elasticSearch.git/cache"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/repository/elasticSearch"
	"github.com/04Akaps/elasticSearch.git/service"
	"go.uber.org/fx"
)

var configFlag = flag.String("config", "./config.toml", "configuration toml file path")

func init() {
	flag.Parse()
}

var Cfg = fx.Module(
	"config",
	fx.Provide(func() config.Config {
		return config.NewConfig(*configFlag)
	}),
)

var CacheManager = fx.Module(
	"cacheManager",
	fx.Provide(func(cfg config.Config) cache.CacheManager {
		return cache.NewCacheManager(cfg)
	}),
)

var ElasticSearch = fx.Module(
	"elasticSearch",
	fx.Provide(func(cfg config.Config) elasticSearch.ElasticSearch {
		return elasticSearch.NewElasticSearch(cfg)
	}),
)

var Service = fx.Module(
	"service",
	fx.Provide(func(cfg config.Config, cacheManager cache.CacheManager, search elasticSearch.ElasticSearch) service.Manager {
		return service.NewManager(cfg, cacheManager, search)
	}),
)
