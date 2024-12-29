package elasticSearch

import (
	"context"
	"fmt"
	"github.com/04Akaps/elasticSearch.git/config"
	_twitter "github.com/04Akaps/elasticSearch.git/types/twitter"
	"github.com/olivere/elastic/v7"
	"log"
	"net/http"
)

type ElasticSearch struct {
	cfg    config.Config
	client *elastic.Client
}

type APIKeyTransport struct {
	APIKey    string
	Transport http.RoundTripper
}

func (t *APIKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Add the Authorization header with the API Key
	req.Header.Set("Authorization", "ApiKey "+t.APIKey)

	// Continue the HTTP request using the provided Transport
	return t.Transport.RoundTrip(req)
}

func NewElasticSearch(cfg config.Config) ElasticSearch {
	log.Println("Start to connect elasticSearch")

	config := cfg.Repository.ElasticSearch

	var connectorConfig []elastic.ClientOptionFunc
	connectorConfig = append(connectorConfig, elastic.SetURL(config.URI))
	connectorConfig = append(connectorConfig, elastic.SetSniff(false))

	if config.ApiKey != "" {
		// connet using api Key
		// -> Elastic Enterprise를 사용하면 Api Key를 통해서 접속을 해야 한다.
		connectorConfig = append(connectorConfig, elastic.SetHttpClient(&http.Client{
			Transport: &APIKeyTransport{
				APIKey:    config.ApiKey,
				Transport: http.DefaultTransport,
			},
		}))
	} else {
		// connect using user info
		connectorConfig = append(connectorConfig, elastic.SetBasicAuth(config.User, config.Password))
	}

	client, err := elastic.NewClient(connectorConfig...)

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

func (e ElasticSearch) Bulk() *elastic.BulkService {
	return e.client.Bulk()
}

func (e ElasticSearch) Indexes() (map[string]bool, error) {
	indices, err := e.client.CatIndices().Do(context.Background())

	if err != nil {
		return nil, err
	}

	res := make(map[string]bool, len(indices))

	for _, index := range indices {
		res[index.Index] = true
	}

	return res, nil

}

func (e ElasticSearch) InsertBulkTest(index string, v []_twitter.SearchResult) {
	bulkRequest := e.client.Bulk()

	for i, doc := range v {
		req := elastic.NewBulkIndexRequest().
			Index(index).
			Id(string(rune(i + 1))).
			Doc(doc)
		bulkRequest = bulkRequest.Add(req)
	}

	ctx := context.Background()
	bulkResponse, err := bulkRequest.Do(ctx)
	if err != nil {
		log.Fatalf("Failed to execute bulk request: %s", err)
	}

	// 응답 처리
	if bulkResponse.Errors {
		log.Println("Bulk request completed with errors")
	} else {
		log.Println("Bulk request succeeded")
	}
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
