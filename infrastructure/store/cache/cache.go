package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type RedisStore struct {
	client *redis.Client
}

func New(redisURL string) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}

	return &RedisStore{client: client}
}

func (s *RedisStore) GetClient() *redis.Client {
	return s.client
}

func (s *RedisStore) Shutdown() {
	if err := s.client.Close(); err != nil {
		log.Error().Err(err)
	}
}
