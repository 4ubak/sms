package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/xakpro/sms/internal/interfaces"
	"time"
)

type RedisSt struct {
	lg interfaces.Logger
	r  *redis.Client
}

func NewRedisSt(lg interfaces.Logger, url, psw string, db int) *RedisSt {
	return &RedisSt{
		lg: lg,
		r: redis.NewClient(&redis.Options{
			Addr:     url,
			Password: psw,
			DB:       db,
		}),
	}
}

func (c *RedisSt) Get(key string) ([]byte, bool, error) {
	data, err := c.r.Get(key).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		c.lg.Errorw("Redis: fail to 'get'", err)
		return nil, false, err
	}
	return data, true, nil
}

func (c *RedisSt) Set(key string, value []byte, expiration time.Duration) error {
	err := c.r.Set(key, value, expiration).Err()
	if err != nil {
		c.lg.Errorw("Redis: fail to 'set'", err)
	}
	return err
}
