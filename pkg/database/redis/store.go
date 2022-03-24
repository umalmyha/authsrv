package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type PipelineFn func(redis.Pipeliner) error
type AttemptsPipelineFn func(int) error

type Store struct {
	rdb *redis.Client
}

func NewStore(rdb *redis.Client) *Store {
	return &Store{
		rdb: rdb,
	}
}

func (s *Store) WithinTx(ctx context.Context, watchKeys []string, pipeFns ...PipelineFn) error {
	if len(pipeFns) == 0 {
		return errors.New("no pipeline functions provided")
	}

	return s.rdb.Watch(ctx, func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			for _, pipeFn := range pipeFns {
				if err := pipeFn(pipe); err != nil {
					return err
				}
			}
			return nil
		})
		return err
	}, watchKeys...)
}

func (s *Store) WithinTxWithAttempts(ctx context.Context, watchKeys []string, pipeFns ...PipelineFn) AttemptsPipelineFn {
	return func(attempts int) error {
		if attempts <= 0 {
			return errors.New("number of attempts must be non-zero positive number")
		}

		for i := 0; i < attempts; i++ {
			err := s.WithinTx(ctx, watchKeys, pipeFns...)
			if err != nil {
				if err == redis.TxFailedErr {
					continue
				}
				return err
			}
			return nil
		}

		return fmt.Errorf("failed to finalize the operations within transaction after %d attempts", attempts)
	}
}

func (s *Store) Client() *redis.Client {
	return s.rdb
}
