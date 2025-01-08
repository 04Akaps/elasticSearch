package main

import (
	"github.com/04Akaps/elasticSearch.git/cmd/dependency"
	"github.com/04Akaps/elasticSearch.git/network"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		dependency.Cfg,
		dependency.CacheManager,
		dependency.ElasticSearch,
		dependency.Service,
		dependency.Ollama,
		fx.Provide(network.NewRouter),
		fx.Invoke(func(network.Router) {}),
	).Run()
}
