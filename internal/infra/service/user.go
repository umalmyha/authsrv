package service

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/umalmyha/authsrv/internal/business/role"
	"github.com/umalmyha/authsrv/internal/business/user"
)

type UserService struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func NewUserService(db *sqlx.DB, rdb *redis.Client) *UserService {
	return &UserService{
		db:  db,
		rdb: rdb,
	}
}

func (srv *UserService) AssignRole(ctx context.Context, username string, roleName string) error {
	uow := user.NewUnitOfWork(srv.db, srv.rdb)
	repo := user.NewRepository(uow)

	user, err := repo.FindByUsername(ctx, username)
	if err != nil {
		return errors.Wrap(err, "failed to find user in repository")
	}

	if user == nil {
		return errors.Errorf("user %s doesn't exist", username)
	}

	if err := user.AssignRole(roleName, srv.findRoleByNameFn(ctx)); err != nil {
		return errors.Wrap(err, "failed to assign role")
	}

	if err := repo.Update(user); err != nil {
		return errors.Wrap(err, "failed to update user in repository")
	}

	return uow.Flush(ctx)
}

func (srv *UserService) UnassignRole(ctx context.Context, username string, roleName string) error {
	uow := user.NewUnitOfWork(srv.db, srv.rdb)
	repo := user.NewRepository(uow)

	user, err := repo.FindByUsername(ctx, username)
	if err != nil {
		return errors.Wrap(err, "failed to find user in repository")
	}

	if user == nil {
		return errors.Errorf("user %s doesn't exist", username)
	}

	if err := user.UnassignRole(roleName, srv.findRoleByNameFn(ctx)); err != nil {
		return errors.Wrap(err, "failed to unassign role")
	}

	if err := repo.Update(user); err != nil {
		return errors.Wrap(err, "failed to update user in repository")
	}

	return uow.Flush(ctx)
}

func (srv *UserService) findRoleByNameFn(ctx context.Context) user.RoleFinderByNameFn {
	return func(name string) (role.RoleDto, error) {
		var dto role.RoleDto
		dto, err := role.NewRoleDao(srv.db).FindByName(ctx, name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return dto, nil
			}
			return dto, err
		}
		return dto, nil
	}
}
