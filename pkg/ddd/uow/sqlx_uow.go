package uow

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type SqlxUnitOfWork struct {
	tx *sqlx.Tx
	db *sqlx.DB
}

func NewSqlxUnitOfWorkWithinTx(db *sqlx.DB, tx *sqlx.Tx) *SqlxUnitOfWork {
	uow := NewSqlxUnitOfWork(db)
	uow.tx = tx
	return uow
}

func NewSqlxUnitOfWork(db *sqlx.DB) *SqlxUnitOfWork {
	return &SqlxUnitOfWork{
		db: db,
	}
}

func (uow *SqlxUnitOfWork) ExtContext() sqlx.ExtContext {
	if uow.tx != nil {
		return uow.tx
	}
	return uow.db
}

func (uow *SqlxUnitOfWork) Tx(ctx context.Context) (*sqlx.Tx, error) {
	if uow.tx != nil {
		return uow.tx, nil
	}
	return uow.db.BeginTxx(ctx, nil)
}
