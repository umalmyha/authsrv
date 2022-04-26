package scope

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
	if err := NewScopeDao(r.db).Create(ctx, dto); err != nil {
		return errors.Wrap(err, "failed to create scope")
	}
	return nil
}
