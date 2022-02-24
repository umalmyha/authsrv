package refresh

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RefreshTokenDao struct {
	client *redis.Client
}

func NewRefreshTokenDao(client *redis.Client) *RefreshTokenDao {
	return &RefreshTokenDao{
		client: client,
	}
}

func (dao *RefreshTokenDao) CreateMultiForUser(ctx context.Context, userId string, tokens []RefreshTokenDto) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(tokens); err != nil {
		return errors.New(fmt.Sprintf("failed to serialize tokens to gob format: %s", err.Error()))
	}
	return dao.client.Set(ctx, userId, buf, 0).Err()
}
