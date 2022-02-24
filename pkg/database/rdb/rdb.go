package rdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type database string

const DatabasePostgres database = "postgres"

func Connect(ctx context.Context, cfg *config) (*sqlx.DB, error) {
	db, err := connect(ctx, cfg)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.maxOpenConns)
	db.SetConnMaxIdleTime(cfg.connMaxIdleTime)
	db.SetConnMaxLifetime(cfg.connMaxLifetime)
	db.SetMaxIdleConns(cfg.maxIdleConns)

	return db, nil
}

func connect(ctx context.Context, cfg *config) (*sqlx.DB, error) {
	switch cfg.database {
	case DatabasePostgres:
		return connectToPostgesql(ctx, cfg)
	default:
		return nil, errors.New(fmt.Sprintf("unsupported database - %s", cfg.database))
	}
}
