package service

import (
	"github.com/04Akaps/elasticSearch.git/cache"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/repository/elasticSearch"
	"github.com/04Akaps/elasticSearch.git/repository/ollama"
	"github.com/04Akaps/elasticSearch.git/service/loop"
	v1 "github.com/04Akaps/elasticSearch.git/service/v1"
)

type Manager struct {
	cfg config.Config

	v1 v1.V1
}

func NewManager(
	cfg config.Config,
	cache *cache.CacheManager,
	elasticSearch elasticSearch.ElasticSearch,
	ollaMa ollama.Ollama,
) Manager {
	m := Manager{
		cfg: cfg,
		v1:  v1.NewV1(cfg, cache, elasticSearch),
	}

	loop.RunTwitterLoop(cfg, elasticSearch)
	loop.RunValidatorLoop(cfg, elasticSearch)
	loop.RunNlpLoop(cfg, elasticSearch, ollaMa)

	return m
}

func (m Manager) V1() v1.V1 {
	return m.v1
}
