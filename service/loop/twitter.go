package loop

import (
	"fmt"
	"github.com/04Akaps/elasticSearch.git/cache"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/repository/elasticSearch"
	"github.com/04Akaps/elasticSearch.git/twitter"
	twitterDependency "github.com/g8rswimmer/go-twitter/v2"
	"log"
	"time"
)

type TweetsLoop struct {
	cfg           config.Config
	ElasticSearch elasticSearch.ElasticSearch
	CacheManager  cache.CacheManager
}

func RunTwitterLoop(
	cfg config.Config,
	elasticSearch elasticSearch.ElasticSearch,
	cacheManager cache.CacheManager,
) {
	l := TweetsLoop{cfg, elasticSearch, cacheManager}

	for key, info := range l.cfg.Twitter {
		twitterClient := twitter.NewTwitterClient(info)

		go l.runTwitterClient(twitterClient, key, info)
	}
}

func (t TweetsLoop) runTwitterClient(client twitter.Twitter, key string, info config.Twitter) {

	opts := twitterDependency.TweetSearchOpts{
		MaxResults: info.Counter,
		//StartTime:  time.Now().Sub(time.Hour * 1),
		//EndTime:
	}

	ticker := time.NewTicker(5e9)
	defer ticker.Stop()

	for {

		result, err := client.SearchTweets(key, opts)

		if err != nil {
			log.Println("Failed to get tweets", "err", err)
		} else {
			// TODO insert to elasticSearch
			fmt.Println("test", result)
		}

		<-ticker.C
	}

}
