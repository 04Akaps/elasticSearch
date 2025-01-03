package loop

import (
	"fmt"
	"github.com/04Akaps/elasticSearch.git/common/http"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/repository/elasticSearch"
	"time"
)

// TODO -> NLP AI를 붙여서 특정 구간에서 현재 Tweets의 내용이 긍정적인지
// -> 추가로 어떤 이야기를 하고 있는지를 요약하는 정보로 추가로 시계열 데이터를 생성

type NlpLoop struct {
	cfg                   config.Config
	ElasticSearch         elasticSearch.ElasticSearch
	HuggingFaceHttpClient *http.Client
}

func RunNlpLoop(
	cfg config.Config,
	elasticSearch elasticSearch.ElasticSearch,
) {
	l := NlpLoop{
		cfg:                   cfg,
		ElasticSearch:         elasticSearch,
		HuggingFaceHttpClient: http.NewClient(10*time.Second, ""),
	}

	fmt.Println(l)
}
