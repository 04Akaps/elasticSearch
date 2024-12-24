package twitter

import (
	"context"
	"fmt"
	"github.com/04Akaps/elasticSearch.git/config"
	_twitter "github.com/04Akaps/elasticSearch.git/types/twitter"
	"github.com/g8rswimmer/go-twitter"
	"github.com/olivere/elastic/v7"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type _authorize struct {
	Token string
}

func (a _authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

type Twitter struct {
	client *twitter.Tweet
}

func NewTwitterClient(cfg config.Twitter) Twitter {

	tweet := &twitter.Tweet{
		Authorizer: _authorize{
			Token: cfg.BearerToken,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	return Twitter{tweet}
}

func (t Twitter) SearchTweets(
	ctx context.Context,
	key string,
	opts twitter.TweetRecentSearchOptions,
	field twitter.TweetFieldOptions,
) ([]*elastic.BulkIndexRequest, error) {

	searchResult, err := t.client.RecentSearch(ctx, key, opts, field)

	if err != nil {
		return nil, err
	}

	var res []*elastic.BulkIndexRequest
	index := 0

	for _, l := range searchResult.LookUps {
		info := l.Tweet

		doc := _twitter.SearchResult{
			Text:     normalizeSpaces(info.Text),
			Language: info.Language,
			AuthorID: info.AuthorID,
			Geo: _twitter.Geo{
				PlaceID:     info.Geo.PlaceID,
				Coordinates: info.Geo.Coordinates.Coordinates,
			},
			ID:        info.ID,
			Source:    info.Source,
			CreatedAt: convertToUnix(info.CreatedAt),
		}

		req := elastic.NewBulkIndexRequest().
			Index(key).
			Id(string(rune(index + 1))).
			Doc(doc)

		res = append(res, req)

		index++
	}

	return res, nil
}

func normalizeSpaces(input string) string {
	// 정규식으로 연속된 공백을 하나로 치환
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(input, " "))
}

func convertToUnix(timeString string) int64 {
	parsedTime, err := time.Parse(time.RFC3339, timeString)

	if err != nil {
		return 0
	}

	return parsedTime.Unix()
}
