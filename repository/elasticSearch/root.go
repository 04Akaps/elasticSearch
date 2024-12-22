package elasticSearch

import (
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/olivere/elastic/v7"
	"log"
)

type ElasticSearch struct {
	cfg    config.Config
	client *elastic.Client
}

func NewElasticSearch(cfg config.Config) ElasticSearch {
	log.Println("Start to connect elasticSearch")

	config := cfg.Repository.ElasticSearch

	client, err := elastic.NewClient(
		elastic.SetBasicAuth(config.User, config.Password),
		elastic.SetURL(config.URI),
		elastic.SetSniff(false),
	)

	if err != nil {
		log.Panic("Failed to connect elasticSearch", "err", err)
	}

	return ElasticSearch{cfg: cfg, client: client}
}
