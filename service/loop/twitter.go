package loop

import (
	"context"
	"github.com/04Akaps/elasticSearch.git/cache"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/repository/elasticSearch"
	"github.com/04Akaps/elasticSearch.git/twitter"
	twitterDependency "github.com/g8rswimmer/go-twitter"
	"log"
	"time"
)

type TweetsLoop struct {
	cfg           config.Config
	ElasticSearch elasticSearch.ElasticSearch
	CacheManager  *cache.CacheManager
}

func RunTwitterLoop(
	cfg config.Config,
	elasticSearch elasticSearch.ElasticSearch,
	cacheManager *cache.CacheManager,
) {
	l := &TweetsLoop{cfg, elasticSearch, cacheManager}

	for key, info := range l.cfg.Twitter {
		twitterClient := twitter.NewTwitterClient(info)

		go l.runTwitterClient(twitterClient, key, info)
	}
}

func (t *TweetsLoop) runTwitterClient(client twitter.Twitter, key string, info config.Twitter) {

	opts := twitterDependency.TweetRecentSearchOptions{
		//StartTime:  time.Now().Sub(time.Hour * 1),
		//EndTime:
		MaxResult: 10,
	}

	ticker := time.NewTicker(time.Duration(info.Ticker) * time.Minute)
	defer ticker.Stop()

	fieldOpts := twitterDependency.TweetFieldOptions{
		TweetFields: []twitterDependency.TweetField{
			twitterDependency.TweetFieldCreatedAt, twitterDependency.TweetFieldLanguage,
			twitterDependency.TweetFieldAuthorID, twitterDependency.TweetFieldText, twitterDependency.TweetFieldGeo,
			twitterDependency.TweetFieldID, twitterDependency.TweetFieldSource,
		},
	}

	for {
		ctx := context.Background()

		bulkClient := t.ElasticSearch.Bulk()

		result, err := client.SearchTweets(ctx, key, opts, fieldOpts)

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

		}

		<-ticker.C
	}

}
