package refresh

import (
	"context"

	"github.com/pkg/errors"

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
		return nil, errors.Wrap(err, "failed to read all tokens for user")
	}

	if err := dbredis.DecodeGob([]byte(tokenStr), &tokens); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize tokens from gob format")
	}
	return tokens, nil
}
