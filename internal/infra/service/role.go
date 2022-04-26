package service

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/umalmyha/authsrv/internal/business/role"
	"github.com/umalmyha/authsrv/internal/business/scope"
)

type RoleService struct {
	db *sqlx.DB
}

func NewRoleService(db *sqlx.DB) *RoleService {
	return &RoleService{
		db: db,
	}
}

func (srv *RoleService) CreateRole(ctx context.Context, nr role.NewRoleDto) error {
	uow := role.NewUnitOfWork(srv.db)
	repo := role.NewRepository(uow)

	existFn := func(name string) (bool, error) {
		if _, err := role.NewRoleDao(srv.db).FindByName(ctx, name); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}

	r, err := role.FromNewRoleDto(nr, existFn)
	if err != nil {
		return errors.Wrap(err, "failed to build role from DTO")
	}

	if err := repo.Add(r); err != nil {
		return errors.Wrap(err, "failed to add role to repository")
	}

	return uow.Flush(ctx)
}

func (srv *RoleService) AssignScope(ctx context.Context, roleName string, scopeName string) error {
	uow := role.NewUnitOfWork(srv.db)
	repo := role.NewRepository(uow)

	r, err := repo.FindByName(ctx, roleName)
	if err != nil {
		return errors.Wrap(err, "failed to find role by name")
	}

	if err := r.AssignScope(scopeName, srv.findScopeByNameFn(ctx)); err != nil {
		return errors.Wrap(err, "failed to assign scope")
	}

	if err := repo.Update(r); err != nil {
		return errors.Wrap(err, "failed to update role in repository")
	}

	return uow.Flush(ctx)
}

func (srv *RoleService) UnassignScope(ctx context.Context, roleName string, scopeName string) error {
	uow := role.NewUnitOfWork(srv.db)
	repo := role.NewRepository(uow)

	r, err := repo.FindByName(ctx, roleName)
	if err != nil {
		return errors.Wrap(err, "failed to find role by name")
	}

	if err := r.UnassignScope(scopeName, srv.findScopeByNameFn(ctx)); err != nil {
		return errors.Wrap(err, "failed to unassign scope")
	}

	if err := repo.Update(r); err != nil {
		return errors.Wrap(err, "failed to update role in repository")
	}

	return uow.Flush(ctx)
}

func (srv *RoleService) findScopeByNameFn(ctx context.Context) role.ScopeFinderByNameFn {
	return func(name string) (scope.ScopeDto, error) {
		var dto scope.ScopeDto
		dto, err := scope.NewScopeDao(srv.db).FindByName(ctx, name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return dto, nil
			}
			return dto, err
		}
		return dto, nil
	}
}
