package twitter

import (
	"context"
	"fmt"
	"github.com/04Akaps/elasticSearch.git/config"
	_twitter "github.com/04Akaps/elasticSearch.git/types/twitter"
	"github.com/g8rswimmer/go-twitter/v2"
	"net/http"
)

type _authorize struct {
	Token string
}

func (a _authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

type Twitter struct {
	client *twitter.Client
}

func NewTwitterClient(cfg config.Twitter) Twitter {

	client := &twitter.Client{
		Authorizer: _authorize{
			Token: cfg.BearerToken,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	return Twitter{client}
}

func (t Twitter) SearchTweets(key string, opts twitter.TweetSearchOpts) ([]_twitter.SearchResult, error) {
	ctx := context.Background()

	response, err := t.client.TweetSearch(ctx, key, opts)

	if err != nil {
		return nil, err
	}

	res := make([]_twitter.SearchResult, len(response.Raw.Tweets))

	for i, tweet := range response.Raw.Tweets {
		fmt.Println(i)
		fmt.Println(tweet.Language, tweet.CreatedAt, tweet.ID, tweet.AuthorID)

		//res[i] = _twitter.SearchResult{
		//	UserName:  tweet.User.Name,
		//	Text:      tweet.Text,
		//	CreatedAt: unixTime,
		//}
	}

	return res, nil
}
