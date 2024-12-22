package main

import (
	"flag"
	"github.com/04Akaps/elasticSearch.git/cmd/dependency"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/docker"
	"github.com/04Akaps/elasticSearch.git/network"
	"go.uber.org/fx"
)

var confFlag = flag.String("config", "./config.toml", "configuration toml file path")

func main() {
	flag.Parse()
	cfg := config.NewConfig(*confFlag)
	docker.Initialize(cfg)

	fx.New(
		dependency.Cfg,
		dependency.CacheManager,
		dependency.ElasticSearch,
		dependency.Service,
		fx.Provide(network.NewRouter),
		fx.Invoke(func(network.Router) {}),
	).Run()
}
