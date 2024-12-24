package v1

import (
	"github.com/04Akaps/elasticSearch.git/types/request"
	_twitter "github.com/04Akaps/elasticSearch.git/types/twitter"
)

func (v1 V1) InsertTest(req request.InsertTestRequest) {
	v := make(map[string]interface{})
	v[req.Key] = req.Value

	v1.elasticSearch.InsertTest(req.Index, v)
}

func (v1 V1) InsertMapperTest() {
	v := make([]_twitter.SearchResult, 2)

	v = append(v, _twitter.SearchResult{
		Text: "test",
	})

	v = append(v, _twitter.SearchResult{
		Text: "test two",
	})

	v1.elasticSearch.InsertBulkTest("insert-mapper-type", v)
}
