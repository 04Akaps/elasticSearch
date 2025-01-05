package elasticSearch

import (
	"context"
	"fmt"
	_twitter "github.com/04Akaps/elasticSearch.git/types/twitter"
	"github.com/olivere/elastic/v7"
	"log"
)

/*
	일부 테스트용 함수.
*/

// Only Test
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
		log.Println("Failed to insert dummy data", "cerr", err)
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
			log.Println("Failed to unMarshal data", "cerr", err)
			continue
		}

		fmt.Println(string(testRes))
	}
}
