package user

import (
	"context"

	"github.com/pkg/errors"

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
		return NewUserDao(repo.uow.ExtContext()).FindByUsername(ctx, username)
	}

	user, err := repo.uow.users.Find(func(user UserDto) bool { return user.Username == username }).IfNotPresent(notPresentFn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read user aggregate")
	}

	if !user.IsPresent() {
		return nil, errors.Errorf("user %s doesn't exist", username)
	}

	userAuth, err := NewUserAuthDao(repo.uow.ExtContext()).FindAllForUser(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	tokens, err := refresh.NewRefreshTokenDao(repo.uow.rdb).FindAllForUser(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	u, err := repo.buildUser(user, userAuth, tokens)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build user from db DTOs")
	}

	return u, repo.uow.RegisterClean(u)
}

func (repo *Repository) buildUser(user UserDto, userAuthDto []UserAuthDto, tokens []refresh.RefreshTokenDto) (*User, error) {
	uniqueScopeNames := make(map[string]bool)
	uniqueRolesWithNames := make(map[string]string)
	roles := make([]string, 0)
	roleIds := make([]valueobj.RoleId, 0)

	for _, auth := range userAuthDto {
		uniqueScopeNames[auth.ScopeName] = true
		uniqueRolesWithNames[auth.RoleId] = auth.RoleName
	}

	userAuth := valueobj.NewUserAuth(roles, helpers.Keys(uniqueScopeNames))

	for roleId, roleName := range uniqueRolesWithNames {
		roles = append(roles, roleName)

		roleIdent, err := valueobj.NewRoleId(roleId)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build role identifier")
		}
		roleIds = append(roleIds, roleIdent)
	}

	refreshTokens := helpers.Map(tokens, func(token refresh.RefreshTokenDto, _ int, _ []refresh.RefreshTokenDto) *refresh.RefreshToken {
		return token.ToRefreshToken()
	})

	return fromDbDtos(user, roleIds, refreshTokens, userAuth)
}
