package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/umalmyha/authsrv/internal/business/scope"
)

type ScopeService struct {
	db *sqlx.DB
}

func NewScopeService(db *sqlx.DB) *ScopeService {
	return &ScopeService{
		db: db,
	}
}

func (srv *ScopeService) CreateScope(ctx context.Context, ns scope.NewScopeDto) error {
	existFn := func(name string) (bool, error) {
		if _, err := scope.NewScopeDao(srv.db).FindByName(ctx, name); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}

	sc, err := scope.FromNewScopeDto(ns, existFn)
	if err != nil {
		return err
	}

	return scope.NewRepository(srv.db).Create(ctx, sc)
}
