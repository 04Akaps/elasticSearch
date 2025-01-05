package elasticSearch

import (
	"context"
	"github.com/04Akaps/elasticSearch.git/common/json"
	"github.com/04Akaps/elasticSearch.git/types/cerr"
	"github.com/04Akaps/elasticSearch.git/types/nlp"
	"github.com/olivere/elastic/v7"
)

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

	result, err := client.Search(index).
		From(offset). // offset(시작 위치)
		Size(limit).  // limit(가져올 문서의 개수)
		Do(ctx)       // 실제 실행

	if err != nil {
		return err
	}

	if result.Hits.TotalHits.Value == 0 {
		return cerr.NoDoc
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
