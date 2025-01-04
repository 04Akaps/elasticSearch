package strategy

import (
	"errors"
	"github.com/04Akaps/elasticSearch.git/cache"
	"github.com/go-redis/redis/v7"
	"golang.org/x/sync/singleflight"
	"log"
	"reflect"
	"time"
)

/*
	API에서 Process 최적화를 위한 singleflight 함수
	동시에 하나의 인스턴스에 request가 몰리게 된다면, sfGroup 변수에서
	같은 키에 대해서는 요청을 wait하게 되고, 하나의 요청이라도 처리가 되면
	같은 키에 대해서 같은 결과를 바로 return 해준다.

	목적 : DB에 대해 동시 다발적인 요청을 하나의 요청으로 처리하기 위해 사용
*/

var sfGroup singleflight.Group

func FetchWithSingleFlight[T any](
	key string,
	stdTTL time.Duration,
	manager cache.CacheManager,
	callBack func() (T, error),
) (res T, remainTTL time.Duration, cacheStored bool, err error) {

	remainTTL, err = manager.GetValueWithTTL(key, &res)

	if err == nil {
		return res, remainTTL, cacheStored, nil
	}

	if !errors.Is(err, redis.Nil) {
		log.Println("Failed to get data from redis", "cerr", err)
		// TODO cerr 커스터마이징
		return res, remainTTL, cacheStored, err
	}

	val, err, _ := sfGroup.Do(key, func() (interface{}, error) {
		res, err = callBack()

		if err != nil {
			return nil, err
		}

		if setErr := manager.SetNX(key, &res, stdTTL); setErr != nil {
			log.Println("Failed to set cache data in flight", "cerr", err)
		} else {
			cacheStored = true
		}

		return res, nil
	})

	if err != nil {
		return res, remainTTL, cacheStored, err
	}

	result, ok := val.(T)

	if !ok {
		log.Println("Failed to cast result to type", reflect.TypeOf(result).String())
		// TODO 에러 커스터마이징
		return res, remainTTL, cacheStored, err
	}

	return result, remainTTL, cacheStored, nil
}
