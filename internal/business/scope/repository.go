package scope

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, scope *Scope) error {
	dto := scope.Dto()
	return NewScopeDao(r.db).Create(ctx, dto)
}
