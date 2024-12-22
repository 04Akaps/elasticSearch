package v1

import "github.com/04Akaps/elasticSearch.git/types/request"

func (v1 V1) InsertTest(req request.InsertTestRequest) {
	v := make(map[string]interface{})
	v[req.Key] = req.Value

	v1.elasticSearch.InsertTest(req.Index, v)
}
