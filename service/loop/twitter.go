package loop

import (
	"context"
	"github.com/04Akaps/elasticSearch.git/common/twitter"
	"log"
	"time"

	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/repository/elasticSearch"
	. "github.com/g8rswimmer/go-twitter"
)

type TweetsLoop struct {
	cfg           config.Config
	ElasticSearch elasticSearch.ElasticSearch
}

func RunTwitterLoop(
	cfg config.Config,
	elasticSearch elasticSearch.ElasticSearch,
) {
	l := TweetsLoop{cfg, elasticSearch}

	for key, info := range l.cfg.Twitter {
		twitterClient := twitter.NewTwitterClient(info)

		go l.runTwitterClient(twitterClient, key, info)
	}
}

func (t *TweetsLoop) runTwitterClient(client twitter.Twitter, key string, info config.Twitter) {
	startTime := info.StartTime

	if startTime == 0 {
		// if zero start current Time
		startTime = time.Now().Unix()
	}

	opts := TweetRecentSearchOptions{
		StartTime: time.Unix(startTime, 0),
		MaxResult: 10,
	}

	ticker := time.NewTicker(time.Duration(info.Ticker) * time.Minute)
	defer ticker.Stop()

	fieldOpts := TweetFieldOptions{
		TweetFields: []TweetField{
			TweetFieldCreatedAt, TweetFieldLanguage,
			TweetFieldAuthorID, TweetFieldText, TweetFieldGeo,
			TweetFieldID, TweetFieldSource,
		},
		UserFields: []UserField{
			UserFieldID, UserFieldUserName, UserFieldName,
		},

		PlaceFields: []PlaceField{
			PlaceFieldCountryCode, PlaceFieldFullName,
		},
	}

	for {
		ctx := context.Background()

		bulkClient := t.ElasticSearch.Bulk()

		result, lastTweetsTime, err := client.SearchTweets(ctx, key, opts, fieldOpts)

		if err != nil {
			log.Println("Failed to get tweets", "err", err)
		} else {

			for _, doc := range result {
				bulkClient = bulkClient.Add(doc)
			}

			response, err := bulkClient.Do(ctx)

			if err != nil {
				log.Println("Failed to send bulk write to elasticSearch", "key", key, "err", err)
			} else if response.Errors {
				log.Println("Bulk request completed with errors", "key", key)
			} else {
				log.Println("Bulk request succeeded", "key", key)
			}

			opts.StartTime = time.Unix(lastTweetsTime, 0)
		}

		<-ticker.C
	}
}
