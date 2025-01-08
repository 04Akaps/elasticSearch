package dependency

import (
	"flag"
	"github.com/04Akaps/elasticSearch.git/cache"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/repository/elasticSearch"
	"github.com/04Akaps/elasticSearch.git/repository/ollama"
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
	fx.Provide(func(cfg config.Config) *cache.CacheManager {
		return cache.NewCacheManager(cfg)
	}),
)

var ElasticSearch = fx.Module(
	"elasticSearch",
	fx.Provide(func(cfg config.Config) elasticSearch.ElasticSearch {
		return elasticSearch.NewElasticSearch(cfg)
	}),
)

var Ollama = fx.Module(
	"ollama",
	fx.Provide(func(cfg config.Config) ollama.Ollama {
		return ollama.NewOllama(cfg.OllaMa.Model, cfg.OllaMa.Url)
	}),
)

var Service = fx.Module(
	"service",
	fx.Provide(func(cfg config.Config, cacheManager *cache.CacheManager, elasticSearch elasticSearch.ElasticSearch, ollaMa ollama.Ollama) service.Manager {
		return service.NewManager(cfg, cacheManager, elasticSearch, ollaMa)
	}),
)
