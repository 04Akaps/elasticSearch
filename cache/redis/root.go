package redis

import (
	"github.com/04Akaps/elasticSearch.git/common/json"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/go-redis/redis/v7"
	"log"
	"math/rand"
	"time"
)

const (
	_defaultCacheBeta = 1
)

type Redis struct {
	client *redis.Client

	//Sync mutex.Mutex

	marshal   func(v interface{}) ([]byte, error)
	unMarshal func(data []byte, v interface{}) error

	_cacheBeta float64

	log log.Logger
}

func NewRedis(cfg config.Config) Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Cache.Redis.DataSource,
		Password: cfg.Cache.Redis.Password,
		DB:       cfg.Cache.Redis.DB,
		Username: cfg.Cache.Redis.UserName,
	})

	log.Println("redis connect started")

	_, err := client.Ping().Result()

	if err != nil {
		log.Panic("Failed to connect to redis client", "err", err)
	}

	r := Redis{
		client: client,
		//Sync:       mutex.NewMutex(client, module, "link:scanner:api", ":"),
		marshal:    json.JsonHandler.Marshal,
		unMarshal:  json.JsonHandler.Unmarshal,
		_cacheBeta: cfg.Cache.Redis.Beta,
	}

	if r._cacheBeta == 0 {
		r._cacheBeta = _defaultCacheBeta
	}

	log.Println("Success to connect redis")

	return r
}

func (r Redis) Set(key string, v interface{}, ttl time.Duration) error {
	value, err := r.marshal(v)

	if err != nil {
		return err
	}

	return r.client.Set(key, value, ttl).Err()
}

func (r Redis) SetNX(key string, v interface{}, ttl time.Duration) error {
	value, err := r.marshal(v)

	if err != nil {
		return err
	}

	return r.client.SetNX(key, value, ttl).Err()
}

func (r Redis) Get(key string, buffer interface{}) error {
	return r.getValue(r.client.Get(key), buffer)
}

func (r Redis) GetValueWithTTL(key string, buffer interface{}) (time.Duration, error) {

	commands := []func(pipe redis.Pipeliner) redis.Cmder{
		func(pipe redis.Pipeliner) redis.Cmder { return pipe.Get(key) },
		func(pipe redis.Pipeliner) redis.Cmder { return pipe.TTL(key) },
	}

	cmds, err := r.executePipeLine(commands)
	if err != nil {
		return 0, err
	}

	getCmd := cmds[0].(*redis.StringCmd)
	err = r.getValue(getCmd, buffer)
	if err != nil {
		return 0, err
	}

	ttlCmd := cmds[1].(*redis.DurationCmd)
	ttl, err := ttlCmd.Result()
	if err != nil {
		return 0, err
	}

	return ttl, nil
}

func (r Redis) executePipeLine(commands []func(pipe redis.Pipeliner) redis.Cmder) ([]redis.Cmder, error) {
	pipe := r.client.Pipeline()

	// 각 명령을 파이프라인에 추가
	cmds := make([]redis.Cmder, len(commands))
	for i, cmdFunc := range commands {
		cmds[i] = cmdFunc(pipe)
	}

	// 파이프라인 실행
	_, err := pipe.Exec()
	if err != nil {
		return nil, err
	}

	return cmds, nil
}

// 남은 TTL을 기반으로 갱신 확률 계산
func (r Redis) PERCompute(std, remainTTL time.Duration) bool {
	remaining := remainTTL.Seconds()
	total := std.Seconds()
	probability := 1 - (remaining / total)

	adjustedProbability := probability * r._cacheBeta

	// 0과 1 사이로 보정
	if adjustedProbability < 0 {
		adjustedProbability = 0
	} else if adjustedProbability > 1 {
		adjustedProbability = 1
	}

	random := rand.Float64()
	return random < adjustedProbability
}

func (r Redis) getValue(getCmd *redis.StringCmd, buffer interface{}) error {
	v, err := getCmd.Bytes()
	if err != nil {
		return err
	}

	err = r.unMarshal(v, buffer)
	if err != nil {
		return err
	}

	return nil
}
