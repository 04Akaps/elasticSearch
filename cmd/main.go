package main

import (
	"flag"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/docker"
)

var confFlag = flag.String("config", "./config.toml", "configuration toml file path")

func main() {
	flag.Parse()
	cfg := config.NewConfig(*confFlag)
	docker.Initialize(cfg)
}
