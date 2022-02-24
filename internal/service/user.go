package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/umalmyha/authsrv/internal/business/user"
	"go.uber.org/zap"
)

type UserService struct {
	logger *zap.SugaredLogger
	db     *sqlx.DB
}

func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (srv *UserService) CreateUser(ctx context.Context, u user.NewUserDto) error {
	uow := user.NewUnitOfWork(srv.db)
	repo := user.NewRepository(uow)

	existUsernameFn := func(username string) (bool, error) {
		if _, err := user.NewUserDao(srv.db).FindByUsername(ctx, u.Username); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}

	user, err := user.FromNewUserDto(u, existUsernameFn)
	if err != nil {
		return err
	}

	if err := repo.Add(user); err != nil {
		return err
	}

	return uow.Flush(ctx)
}
