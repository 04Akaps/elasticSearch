package v1

import "github.com/04Akaps/elasticSearch.git/types/request"

func (v1 V1) ReadTest(req request.ReadTestRequest) {
	v1.elasticSearch.ReadTest(req.Index, req.Key, req.Value)
}
