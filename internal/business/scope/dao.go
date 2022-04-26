package scope

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type ScopeDao struct {
	ec sqlx.ExtContext
}

func NewScopeDao(ec sqlx.ExtContext) *ScopeDao {
	return &ScopeDao{
		ec: ec,
	}
}

func (d *ScopeDao) Create(ctx context.Context, sc ScopeDto) error {
	q := "INSERT INTO SCOPES(ID, NAME, DESCRIPTION) VALUES($1, $2, $3)"
	if _, err := d.ec.ExecContext(ctx, q, sc.Id, sc.Name, sc.Description); err != nil {
		return errors.Wrap(err, "faield to create scope")
	}
	return nil
}

func (d *ScopeDao) FindById(ctx context.Context, id string) (ScopeDto, error) {
	var sc ScopeDto
	q := "SELECT * FROM SCOPES WHERE ID = $1 LIMIT 1"
	if err := sqlx.GetContext(ctx, d.ec, &sc, q, id); err != nil {
		return sc, errors.Wrap(err, "failed to read role by id")
	}
	return sc, nil
}

func (d *ScopeDao) FindByName(ctx context.Context, name string) (ScopeDto, error) {
	var sc ScopeDto
	q := "SELECT * FROM SCOPES WHERE NAME = $1 LIMIT 1"
	if err := sqlx.GetContext(ctx, d.ec, &sc, q, name); err != nil {
		return sc, errors.Wrap(err, "failed to read role by name")
	}
	return sc, nil
}
