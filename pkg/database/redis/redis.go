package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func Connect(opts *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, errors.New(fmt.Sprintf("Error on connection to redis: %s", err.Error()))
	}
	return client, nil
}
