package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/umalmyha/authsrv/internal/business/refresh"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/pkg/helpers"
)

type Repository struct {
	uow *unitOfWork
}

func NewRepository(u *unitOfWork) *Repository {
	return &Repository{
		uow: u,
	}
}

func (repo *Repository) Add(user *User) error {
	return repo.uow.RegisterNew(user)
}

func (repo *Repository) Update(user *User) error {
	return repo.uow.RegisterAmended(user)
}

func (repo *Repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	notPresentFn := func() (UserDto, error) {
		var dto UserDto
		dto, err := NewUserDao(repo.uow.ExtContext()).FindByUsername(ctx, username)
		if err != nil {
			return dto, err
		}
		return dto, nil
	}

	user, err := repo.uow.users.Find(func(user UserDto) bool {
		return user.Username == username
	}).IfNotPresent(notPresentFn)
	if err != nil {
		return nil, err
	}

	if !user.IsPresent() {
		return nil, errors.New(fmt.Sprintf("user %s doesn't exist", username))
	}

	userAuth, err := NewUserAuthDao(repo.uow.ExtContext()).FindAllForUser(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	tokens, err := refresh.NewRefreshTokenDao(repo.uow.rdb).FindAllForUser(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	return repo.buildUser(user, userAuth, tokens)
}

func (repo *Repository) buildUser(user UserDto, userAuthDto []UserAuthDto, tokens []refresh.RefreshTokenDto) (*User, error) {
	uniqueScopes := make(map[string]bool)
	roles := make([]string, 0)
	roleIds := make([]valueobj.RoleId, 0)

	for _, auth := range userAuthDto {
		uniqueScopes[auth.ScopeName] = true
		roles = append(roles, auth.RoleName)

		roleId, err := valueobj.NewRoleId(auth.RoleId)
		if err != nil {
			return nil, err
		}
		roleIds = append(roleIds, roleId)
	}

	userAuth := valueobj.NewUserAuth(roles, helpers.Keys(uniqueScopes))

	refreshTokens := helpers.Map(tokens, func(token refresh.RefreshTokenDto, _ int, _ []refresh.RefreshTokenDto) *refresh.RefreshToken {
		return token.ToRefreshToken()
	})

	return fromDbDtos(user, roleIds, refreshTokens, userAuth)
}
