package redis

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"common/config"

	goRedis "github.com/go-redis/redis"
)

var gClient *goRedis.Client
var one sync.Once

func InitRedis() {
	one.Do(func() {
		cli, err := initRedis()
		if err != nil {
			panic(err)
		}
		gClient = cli
	})
}

func initRedis() (*goRedis.Client, error) {
	conf := config.GetRedisConfig()
	if conf == nil {
		return nil, fmt.Errorf("RedisConfig is nil")
	}

	cli := goRedis.NewClient(&goRedis.Options{
		Addr: conf.Address,
		Password: conf.Password,
		DB: conf.DB,

		//连接池容量及闲置连接数量
		PoolSize:     conf.PoolSize, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: conf.MinIdleConns, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量
	})
	if cli == nil {
		return nil, errors.New("init Redis client fail")
	}
	return cli, nil
}

func KeyExists(key string) (bool, error) {
	exits, err := gClient.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return exits == 1, nil
}

func Get(key string) (string, error) {
	return gClient.Get(key).Result()
}

func Set(key, value string, expiration time.Duration) error {
	return gClient.Set(key, value, expiration).Err()
}

func HGet(key, field string) (string, error) {
	return gClient.HGet(key, field).Result()
}

func HSet(key, field, value string) error {
	return gClient.HSet(key, field, value).Err()
}

func HDel(key string, fields ...string) error {
	return gClient.HDel(key, fields...).Err()
}

func HMSet(key string, kv map[string]interface{}) error {
	return gClient.HMSet(key, kv).Err()
}

func HGetAll(key string) (map[string]string, error) {
	return gClient.HGetAll(key).Result()
}

func Del(key string) error {
	return gClient.Del(key).Err()
}