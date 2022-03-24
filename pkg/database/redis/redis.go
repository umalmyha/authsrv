package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// TODO: Think of pool logic
func Connect(opts *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("Error on connection to redis: %s", err.Error())
	}
	return client, nil
}
