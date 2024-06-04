package component

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"tutorial.io/tutorial-order/internal"
)

func NewRedisClient(
	lc fx.Lifecycle,
	config *shared.Config,
) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword, // allow empty password
		DB:       config.RedisDb,
	})
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return redisClient.Close()
		},
	})
	return redisClient
}
