package loop

import (
	"errors"
	"fmt"
	"github.com/04Akaps/elasticSearch.git/common/http"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/repository/elasticSearch"
	"github.com/04Akaps/elasticSearch.git/types/cerr"
	"github.com/04Akaps/elasticSearch.git/types/nlp"
	"github.com/robfig/cron"
	"log"
	"sync"
	"time"
)

// -> NLP AI를 붙여서 특정 구간에서 현재 Tweets의 내용이 긍정적인지
// -> 추가로 어떤 이야기를 하고 있는지를 요약하는 정보로 추가로 시계열 데이터를 생성

const (
	_nlpSuffix = "nlp"
)

type NlpLoop struct {
	cfg                   config.Config
	ElasticSearch         elasticSearch.ElasticSearch
	HuggingFaceHttpClient *http.Client
	c                     *cron.Cron

	nlpDataMapper map[string]nlp.NlpDoc
}

func RunNlpLoop(
	cfg config.Config,
	elasticSearch elasticSearch.ElasticSearch,
) {
	l := NlpLoop{
		cfg:                   cfg,
		ElasticSearch:         elasticSearch,
		HuggingFaceHttpClient: http.NewClient(10*time.Second, ""),
		c:                     cron.New(),
		nlpDataMapper:         make(map[string]nlp.NlpDoc, len(cfg.Twitter)),
	}

	go l.runNlpLoop()
}

func (n *NlpLoop) runNlpLoop() {
	n.c.Start()

	// 00, 15, 30, 45분 주기로 실행
	n.c.AddFunc("0 */15 * * * *", func() {
		n.tweetsSummary()
	})

	n.c.Run()
}

func (n *NlpLoop) tweetsSummary() {
	/*
			값을 조회하고 처리할 기준이 필요하다.
			항상 현재 시점의 값만 처리해야 하는것이 아니라,
			처리가 안되어있는 시간부터 처리를 해야 할 필요가 있으니.

			1. 마지막으로 들어간 Doc을 가져오고, 들어간 시점 이후의 시간으로 필터를 건다.
		  	2. 해당 시점 이후의 값을 Redis의 Scan 형태로 가져 온 후에, 데이터를 처리한다.
				-> 한번에 다 가져오면 부하가 걸린다는 생각이 들어서
			3. 그 후 업데이트가 다 되면, 그 후에 조회시점을 업데이트 한다.
	*/

	var works sync.WaitGroup
	works.Add(len(n.cfg.Twitter))

	for key, _ := range n.cfg.Twitter {
		nlpKey := fmt.Sprintf("%s:%s", key, _nlpSuffix)
		// 1. 마지막 NLP Doc의 CreatedAt 시간을 가져오자.

		var lastNlpDoc nlp.NlpDoc
		err := elasticSearch.FindLatestNlpDoc[nlp.NlpDoc](n.ElasticSearch.Client(), nlpKey, lastNlpDoc)

		if err != nil {
			if !errors.Is(err, cerr.NoDoc) {
				log.Println("Can't get last nlp doc", "key", key, "err", err)
			}
			// 만약 최초의 요청이라면, 즉 데이터를 이제 막 수집하기 시작해서 값이 없다면,
			// 그냥 다음 요청에 진행한다.
			continue
		}

		// Redis의 Scan 처럼 데이터를 offet, limit을 적용해서 부분적으로 가져온다.
		go n.processTweetData(nlpKey, lastNlpDoc.CreatedAt, works)

	}

	works.Wait()

	// 데이터 모두 업데이트 이후에 초기화
	n.nlpDataMapper = make(map[string]nlp.NlpDoc, len(n.cfg.Twitter))
}

func (n *NlpLoop) processTweetData(
	key string,
	startTIme int64,
	works sync.WaitGroup,
) {
	defer works.Done()

	offset := 0
	limit := 20

	var totalText string
	// -> 해당 값에 담아서 tweets의 모든 Text 정보를 수집
	// 이후 Hugging API에 해당 Text를 담아서 원하는 데이터 추출

	for {

		// 20개씩 순차적으로 가져와서 처리
		limit += 20
		offset += 1
		break
	}

}

//var buffer []twitter.Tweet
//
//var res twitter.Tweet
////
////index string,
////	offset, limit int,
////	buffer []T, // 제네릭 타입 배열로 받기
//
//elasticSearch.FindByKey[twitter.Tweet](
//	n.ElasticSearch.Client(),
//	)
//
//n.ElasticSearch.
//
//
//k := fmt.Sprintf("%s:%s", key, _nlpSuffix)
