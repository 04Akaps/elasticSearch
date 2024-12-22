package service

import (
	"github.com/04Akaps/elasticSearch.git/cache"
	"github.com/04Akaps/elasticSearch.git/config"
	v1 "github.com/04Akaps/elasticSearch.git/service/v1"
)

type Manager struct {
	v1 v1.V1
}

func NewManager(cfg config.Config, cache cache.CacheManager) Manager {
	m := Manager{
		v1: v1.NewV1(cfg, cache),
	}

	return m
}

func (m Manager) V1() v1.V1 {
	return m.v1
}
