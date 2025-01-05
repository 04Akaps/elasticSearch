package elasticSearch

import (
	"context"
	"github.com/04Akaps/elasticSearch.git/common/json"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/types/cerr"
	"github.com/04Akaps/elasticSearch.git/types/twitter"
	"github.com/olivere/elastic/v7"
	"log"
	"net/http"
	"strings"
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
		log.Panic("Failed to connect elasticSearch", "cerr", err)
	}

	_, _, err = client.Ping(config.URI).Do(context.Background())

	if err != nil {
		log.Panic("Failed to ping to elasticSearch node", "cerr", err)
	}

	log.Println("Success to connect elasticSearch")

	return ElasticSearch{cfg: cfg, client: client}
}

func (e ElasticSearch) Bulk() *elastic.BulkService {
	return e.client.Bulk()
}

func (e ElasticSearch) Client() *elastic.Client {
	return e.client
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

func (e ElasticSearch) FindTweetsText(
	index string,
	offset, limit int,
	startTime int64,
) (string, int64, error) {
	ctx := context.Background()

	query := elastic.NewBoolQuery().
		Must(elastic.NewRangeQuery("createdAt").Gt(startTime))

	result, err := e.client.Search(index).
		Query(query).
		Sort("createdAt", true).
		From(offset).
		Size(limit).
		Do(ctx)

	if err != nil {
		return "", 0, err
	}

	if result.Hits.TotalHits.Value == 0 {
		return "", 0, cerr.NoDoc
	}

	var builder strings.Builder

	for _, hit := range result.Hits.Hits {
		var item twitter.SearchResult

		err = json.JsonHandler.Unmarshal(hit.Source, &item)

		if err != nil {
			return "", 0, err
		}

		builder.WriteString(item.Text)
	}

	return builder.String(), result.Hits.TotalHits.Value, nil
}
