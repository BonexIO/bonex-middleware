package redis

import (
	"fmt"
	"github.com/wedancedalot/decimal"
	"github.com/garyburd/redigo/redis"
)

const (
	VolumePrefix = "uwv"
)

func (this *redisDAO) makeVolumeRedisKey(address string) string {
	return fmt.Sprintf("%s_%s", VolumePrefix, address)
}

func (this *redisDAO) GetAccountVolume(ipAddress string) (decimal.Decimal, error) {
	redisClient := this.redis.Get()
	defer redisClient.Close()

	result, err := redis.String(redisClient.Do("GET", this.makeVolumeRedisKey(ipAddress)))
	if err != nil {
		if err == redis.ErrNil {
			return decimal.Zero, nil
		}
		return decimal.Zero, err
	}

	v, err := decimal.NewFromString(result)
	if err != nil {
		return decimal.Zero, err
	}

	return v, nil
}

func (this *redisDAO) SetAccountVolume(ipAddress string, volume decimal.Decimal, ttl int64) error {
	redisClient := this.redis.Get()
	defer redisClient.Close()

	return redisClient.Send("SETEX", this.makeVolumeRedisKey(ipAddress), ttl, volume.String())
}