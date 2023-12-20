package redis

import (
	"fmt"
	"web_app/settings"

	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var (
	client *redis.Client
	Nil = redis.Nil
)

// 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			//viper.GetString("redis.host"),
			//viper.GetInt("redis.port"),
			cfg.Host,
			cfg.Port,
		),
		//Password: viper.GetString("redis.password"), // no password set
		//DB:       viper.GetInt("redis.db"),          // use default DB
		//PoolSize: viper.GetInt("redis.pool_size"),
		Password: cfg.Password, // no password set
		DB:       cfg.DB,          // use default DB
		PoolSize: cfg.PoolSize,
	})

	_, err = client.Ping().Result()
	return
}

func Close() {
	_ = client.Close()
}
