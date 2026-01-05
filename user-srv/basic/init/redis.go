package init

import (
	"5/work/Newyear/user-srv/basic/config"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func InitRedis() {
	redisConf := config.AppConf.Redis
	config.Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
		Password: redisConf.Password, // no password set
		DB:       0,                  // use default DB
	})

	err := config.Rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("redis init success")
}
