package redis

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"bonex-middleware/config"
	"bonex-middleware/dao"
)

type redisDAO struct {
	redis *redis.Pool
}

func NewRedis(c *config.Config) (dao.RedisDAO, error) {
	// Redis
	r := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   5,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port))
		},
	}
	return &redisDAO{redis: r}, nil
}
