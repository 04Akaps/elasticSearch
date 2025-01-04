package elasticSearch

import (
	"context"
	"github.com/04Akaps/elasticSearch.git/common/json"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/types/cerr"
	"github.com/04Akaps/elasticSearch.git/types/nlp"
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

func FindLatestNlpDoc[T nlp.NlpDoc](
	client *elastic.Client,
	index string,
	buffer T,
) error {
	ctx := context.Background()

	result, err := client.Search(index).
		Sort("createdAt", false). // 내림차순
		Size(1).Do(ctx)           // 1개만 조회

	if err != nil {
		return err
	}

	if result.Hits.TotalHits.Value == 0 {
		return cerr.NoDoc
	}

	err = json.JsonHandler.Unmarshal(result.Hits.Hits[0].Source, &buffer)

	if err != nil {
		return err
	}

	return nil
}

func FindByKey[T any](
	client *elastic.Client,
	index string,
	offset, limit int,
	buffer []T, // 제네릭 타입 배열로 받기
) error {
	ctx := context.Background()

	// Elasticsearch 검색 요청
	result, err := client.Search(index).
		From(offset). // offset(시작 위치)
		Size(limit).  // limit(가져올 문서의 개수)
		Do(ctx)       // 실제 실행

	if err != nil {
		return err
	}

	// 결과에서 각 히트를 처리하고, buffer에 추가
	for _, hit := range result.Hits.Hits {
		var item T

		err = json.JsonHandler.Unmarshal(hit.Source, &item)

		if err != nil {
			return err
		}

		buffer = append(buffer, item)
	}

	return nil
}
