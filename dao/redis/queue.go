package redis

import (
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"bonex-middleware/types"
)

const SendQueueName = "bonex_faucet"

func (this *redisDAO) AddToQueue(qi *types.QueueItem) error {
	redisClient := this.redis.Get()
	defer redisClient.Close()

	data, err := json.Marshal(qi)
	if err != nil {
		return err
	}

	return redisClient.Send("RPUSH", SendQueueName, string(data))
}

func (this *redisDAO) PopFromQueue() (*types.QueueItem, error) {
	redisClient := this.redis.Get()
	defer redisClient.Close()

	result, err := redis.Bytes(redisClient.Do("RPOP", SendQueueName))
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		}
		return nil, err
	}

	var qi types.QueueItem

	err = json.Unmarshal(result, &qi)
	if err != nil {
		return nil, err
	}

	return &qi, nil
}
