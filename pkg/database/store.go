package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TxExecFunc func(*sqlx.Tx) error

type SqlContextExecFunc func(sqlx.ExtContext) error

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) DB() *sqlx.DB {
	return s.db
}

func (s *Store) WithinTx(ctx context.Context, execFn TxExecFunc) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err := execFn(tx); err != nil {
		return tx.Rollback()
	}
	return tx.Commit()
}
