package elasticSearch

import (
	"context"
	"fmt"
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

	_, _, err = client.Ping(config.URI).Do(context.Background())

	if err != nil {
		log.Panic("Failed to ping to elasticSearch node", "err", err)
	}

	log.Println("Success to connect elasticSearch")

	return ElasticSearch{cfg: cfg, client: client}
}

// just query test
func (e ElasticSearch) InsertTest(index string, v interface{}) {

	ctx := context.Background()

	_, err := e.client.
		Index().Index(index).
		BodyJson(v).Do(ctx)

	if err != nil {
		log.Println("Failed to insert dummy data", "err", err)
		return
	}
}

// just query test
func (e ElasticSearch) ReadTest(index, key, value string) {
	query := elastic.NewMatchQuery(key, value)
	ctx := context.Background()

	result, err := e.client.Search(index).Query(query).Do(ctx)

	if err != nil {
		log.Println("Failed get data", "err", err)
		return
	}

	for _, hit := range result.Hits.Hits {
		var testRes []byte

		err = hit.Source.UnmarshalJSON(testRes)

		if err != nil {
			log.Println("Failed to unMarshal data", "err", err)
			continue
		}

		fmt.Println(string(testRes))
	}
}
