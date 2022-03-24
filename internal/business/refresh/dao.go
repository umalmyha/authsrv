package refresh

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	dbredis "github.com/umalmyha/authsrv/pkg/database/redis"
)

type RefreshTokenDao struct {
	*dbredis.Store
}

func NewRefreshTokenDao(rdb *redis.Client) *RefreshTokenDao {
	return &RefreshTokenDao{
		Store: dbredis.NewStore(rdb),
	}
}

func (dao *RefreshTokenDao) FindAllForUser(ctx context.Context, userId string) ([]RefreshTokenDto, error) {
	tokens := make([]RefreshTokenDto, 0)

	tokenStr, err := dao.Client().Get(ctx, userId).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return tokens, nil
		}
		return nil, err
	}

	if err := dbredis.DecodeGob([]byte(tokenStr), &tokens); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to deserialize tokens from gob format: %s", err.Error()))
	}
	return tokens, nil
}
