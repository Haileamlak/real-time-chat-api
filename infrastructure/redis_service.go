package infrastructure

import (
	"github.com/go-redis/redis/v8"
)

type RedisService interface {
	GetClient() *redis.Client
	Close() error

}

type redisService struct {
	client *redis.Client
}

func NewRedisService(addr string) RedisService {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &redisService{client: rdb}
}

func (r *redisService) GetClient() *redis.Client {
	return r.client
}

func (r *redisService) Close() error {
	return r.client.Close()
}
