package infra

import (
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

// NewCache returns a new Redis client based on the provided Redis URL.
func NewCache(redisURL string) (rdb *redis.Client) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse cache url")
	}

	rdb = redis.NewClient(opt)
	_, err = rdb.Ping().Result()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to call cache")
	}

	return
}
