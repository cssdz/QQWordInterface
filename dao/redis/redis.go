package redis

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var rdb *redis.Client

func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("redis.host"),
			viper.GetInt("redis.port")),
		Password: "",                       // no password set
		DB:       viper.GetInt("redis.db"), // use default DB
		PoolSize: 100,
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		zap.L().Error("connect Redis failed", zap.Error(err))
	}
	return err
}

func Close() {
	_ = rdb.Close()
}
