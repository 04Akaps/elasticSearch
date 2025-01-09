package loop

import (
	"context"
	"errors"
	"fmt"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/repository/elasticSearch"
	"github.com/04Akaps/elasticSearch.git/repository/ollama"
	"github.com/04Akaps/elasticSearch.git/types/cerr"
	"github.com/04Akaps/elasticSearch.git/types/nlp"
	"github.com/olivere/elastic/v7"
	"github.com/robfig/cron"
	"log"
	"strings"
	"sync"
	"time"
)

// -> NLP AI를 붙여서 특정 구간에서 현재 Tweets의 내용이 긍정적인지
// -> 추가로 어떤 이야기를 하고 있는지를 요약하는 정보로 추가로 시계열 데이터를 생성

const (
	_nlpSuffix = "nlp"
)

type NlpLoop struct {
	cfg           config.Config
	elasticSearch elasticSearch.ElasticSearch
	ollaMa        ollama.Ollama
	c             *cron.Cron

	nlpDataMapper map[string]nlp.NlpDoc
}

func RunNlpLoop(
	cfg config.Config,
	elasticSearch elasticSearch.ElasticSearch,
	ollaMa ollama.Ollama,
) {
	l := NlpLoop{
		cfg:           cfg,
		elasticSearch: elasticSearch,
		ollaMa:        ollaMa,
		c:             cron.New(),
		nlpDataMapper: make(map[string]nlp.NlpDoc, len(cfg.Twitter)),
	}

	go l.runNlpLoop()
}

func (n *NlpLoop) runNlpLoop() {
	n.c.Start()

	// 00, 15, 30, 45분 주기로 실행
	// -> 15분 간격이기 떄문에 lock은 따로 걸지 않았다.
	n.c.AddFunc("0 */15 * * * *", func() {
		n.tweetsSummary()
	})

	n.c.Run()
}

func (n *NlpLoop) tweetsSummary() {
	/*
		값을 조회하고 처리할 기준이 필요하다.
		항상 현재 시점의 값만 처리해야 하는것이 아니라,
		처리가 안되어있는 시간부터 처리를 해야 할 필요가 있다.
	*/

	var works sync.WaitGroup
	works.Add(len(n.cfg.Twitter))

	for key, _ := range n.cfg.Twitter {
		nlpKey := fmt.Sprintf("%s:%s", key, _nlpSuffix)

		var lastNlpDoc nlp.NlpDoc
		err := elasticSearch.FindLatestNlpDoc[nlp.NlpDoc](n.elasticSearch.Client(), nlpKey, lastNlpDoc)

		if err != nil {
			if !errors.Is(err, cerr.NoDoc) {
				log.Println("Can't get last nlp doc", "key", key, "err", err)
			}
			// 만약 최초의 요청이라면, 즉 데이터를 이제 막 수집하기 시작해서 값이 없다면,
			// 그냥 다음 요청에 진행한다.
			continue
		}

		// Redis의 Scan 처럼 데이터를 offet, limit을 적용해서 부분적으로 가져온다.
		go n.processTweetData(nlpKey, lastNlpDoc.CreatedAt, &works)
	}

	works.Wait()

	n.ProcessSummarized()

	// 데이터 모두 업데이트 이후에 초기화
	n.nlpDataMapper = make(map[string]nlp.NlpDoc, len(n.cfg.Twitter))
}

func (n *NlpLoop) processTweetData(
	key string,
	startTIme int64,
	works *sync.WaitGroup,
) {
	defer works.Done()

	offset := 0
	limit := 20

	var builder strings.Builder
	var counts int64

	for {

		texts, count, err := n.elasticSearch.FindTweetsText(key, offset, limit, startTIme)

		if err != nil {
			if !errors.Is(err, cerr.NoDoc) {
				log.Println("Failed to get tweets message", "key", key, "err", err)
			}
			break
		}

		builder.WriteString(texts)

		// 20개씩 순차적으로 가져와서 처리
		limit += 20
		offset += 1
		counts += count
	}

	builder.WriteString("This is tweets some posts summary this text")

	summarized, err := n.ollaMa.Call(builder.String())

	if err != nil {
		log.Println("Failed to call ollaMa", "key", key, "err", err)
		return
	}

	n.nlpDataMapper[key] = nlp.NlpDoc{
		Summary:                  summarized,
		TotalAggregatedDocuments: counts,
		CreatedAt:                time.Now().Unix(),
	}
}

func (n *NlpLoop) ProcessSummarized() {
	bulkClient := n.elasticSearch.Bulk()

	index := 0

	for k, doc := range n.nlpDataMapper {
		req := elastic.NewBulkIndexRequest().
			Index(k).
			Id(string(rune(index + 1))).
			Doc(doc)

		bulkClient = bulkClient.Add(req)
	}

	n.elasticSearch.BulkDo(context.Background(), bulkClient)
}
