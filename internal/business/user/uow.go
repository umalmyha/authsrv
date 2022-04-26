package user

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/umalmyha/authsrv/internal/business/refresh"
	dbredis "github.com/umalmyha/authsrv/pkg/database/redis"
	"github.com/umalmyha/authsrv/pkg/ddd/uow"
	"github.com/umalmyha/authsrv/pkg/helpers"
)

type unitOfWork struct {
	*uow.SqlxUnitOfWork
	rdb           *redis.Client
	users         *uow.ChangeSet[UserDto]
	assignedRoles *uow.ChangeSet[RoleAssignmentDto]
	tokens        *uow.ChangeSet[refresh.RefreshTokenDto]
}

func NewUnitOfWork(db *sqlx.DB, rdb *redis.Client) *unitOfWork {
	return &unitOfWork{
		SqlxUnitOfWork: uow.NewSqlxUnitOfWork(db),
		rdb:            rdb,
		users:          uow.NewChangeSet[UserDto](),
		assignedRoles:  uow.NewChangeSet[RoleAssignmentDto](),
		tokens:         uow.NewChangeSet[refresh.RefreshTokenDto](),
	}
}

func (uow *unitOfWork) RegisterClean(user *User) error {
	uow.users.Attach(user.ToDto())
	uow.assignedRoles.AttachRange(user.RolesDto()...)
	uow.tokens.AttachRange(user.TokensDto()...)
	return nil
}

func (uow *unitOfWork) RegisterNew(user *User) error {
	if err := uow.users.Add(user.ToDto()); err != nil {
		return errors.Wrap(err, "failed to add user DTO to changeset")
	}

	if err := uow.assignedRoles.AddRange(user.RolesDto()...); err != nil {
		return errors.Wrap(err, "failed to add roles assignments DTOs to changeset")
	}

	if err := uow.tokens.AddRange(user.TokensDto()...); err != nil {
		return errors.Wrap(err, "failed to add tokens DTOs to changeset")
	}

	return nil
}

func (uow *unitOfWork) RegisterDeleted(user *User) error {
	if err := uow.users.Remove(user.ToDto()); err != nil {
		return errors.Wrap(err, "failed to delete user DTO in changeset")
	}

	if err := uow.assignedRoles.RemoveRange(user.RolesDto()...); err != nil {
		return errors.Wrap(err, "failed to delete roles assignments DTOs in changeset")
	}

	if err := uow.tokens.RemoveRange(user.TokensDto()...); err != nil {
		return errors.Wrap(err, "failed to delete tokens DTOs in changeset")
	}

	return nil
}

func (uow *unitOfWork) RegisterAmended(user *User) error {
	userDto := user.ToDto()
	if err := uow.users.Update(userDto); err != nil {
		return errors.Wrap(err, "failed to update user DTO in changeset")
	}

	createdRoles, _, deletedRoles := uow.assignedRoles.DeltaWithMatched(user.RolesDto(), func(role RoleAssignmentDto) bool {
		return role.UserId == userDto.Id
	})

	if err := uow.assignedRoles.AddRange(createdRoles...); err != nil {
		return errors.Wrap(err, "failed to add roles assignments DTOs to changeset")
	}

	if err := uow.assignedRoles.RemoveRange(deletedRoles...); err != nil {
		return errors.Wrap(err, "failed to remove roles assignments DTOs in changeset")
	}

	createdTokens, _, deletedTokens := uow.tokens.DeltaWithMatched(user.TokensDto(), func(token refresh.RefreshTokenDto) bool {
		return token.UserId == userDto.Id
	})

	if err := uow.tokens.AddRange(createdTokens...); err != nil {
		return errors.Wrap(err, "failed to add tokens DTOs to changeset")
	}

	if err := uow.tokens.RemoveRange(deletedTokens...); err != nil {
		return errors.Wrap(err, "failed to delete tokens DTOs in changeset")
	}

	return nil
}

func (uow *unitOfWork) Flush(ctx context.Context) error {
	tx, err := uow.Tx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to open transaction")
	}
	defer tx.Rollback()

	userDao := NewUserDao(tx)
	roleDao := NewRoleAssignmentDao(tx)
	tokenDao := refresh.NewRefreshTokenDao(uow.rdb)

	userTokens, err := uow.usersEncodedTokens()
	if err != nil {
		return errors.Wrap(err, "failed to encode user tokens")
	}

	if rmRoles := uow.assignedRoles.Deleted(); len(rmRoles) > 0 {
		rolesGroups := helpers.GroupBy(rmRoles, func(role RoleAssignmentDto, _ int, _ []RoleAssignmentDto) (string, string) {
			return role.UserId, role.RoleId
		})

		for userId, rolesIds := range rolesGroups {
			if err := roleDao.DeleteByUserIdAndRoleIdsIn(ctx, userId, rolesIds); err != nil {
				return errors.Wrap(err, "failed to process roles assignments deletion")
			}
		}
	}

	if rmUsers := uow.users.Deleted(); len(rmUsers) > 0 {
		mapper := func(user UserDto, _ int, _ []UserDto) string {
			return user.Id
		}
		if err := userDao.DeleteWhereIdsIn(ctx, helpers.Map(rmUsers, mapper)); err != nil {
			return errors.Wrap(err, "failed to process users deletion")
		}
	}

	if createdUsers := uow.users.Created(); len(createdUsers) > 0 {
		if err := userDao.CreateMulti(ctx, createdUsers); err != nil {
			return errors.Wrap(err, "failed to process users creation")
		}
	}

	if createdRoles := uow.assignedRoles.Created(); len(createdRoles) > 0 {
		if err := roleDao.CreateMulti(ctx, createdRoles); err != nil {
			return errors.Wrap(err, "failed to process roles assignments creation")
		}
	}

	if updatedUsers := uow.users.Updated(); len(updatedUsers) > 0 {
		for _, updUser := range updatedUsers {
			if err := userDao.Update(ctx, updUser); err != nil {
				return errors.Wrap(err, "failed to process users update")
			}
		}
	}

	if len(userTokens) > 0 {
		userIds := helpers.Keys(userTokens)
		pipeline := tokenDao.WithinTxWithAttempts(ctx, userIds, func(pipe redis.Pipeliner) error {
			for userId, encodedToken := range userTokens {
				if err := pipe.Set(ctx, userId, encodedToken, 0).Err(); err != nil {
					return errors.Wrap(err, "failed to update refresh token")
				}
			}
			return nil
		})

		if err := pipeline(3); err != nil {
			return errors.Wrap(err, "failed to update refresh tokens in transaction pipeline")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}
	return uow.Dispose()
}

func (uow *unitOfWork) Dispose() error {
	uow.users.Cleanup()
	uow.assignedRoles.Cleanup()
	uow.tokens.Cleanup()
	return nil
}

func (uow *unitOfWork) usersEncodedTokens() (map[string]string, error) {
	userTokensEncoded := make(map[string]string)
	if tokens := uow.tokens.All(); len(tokens) > 0 {
		userTokens := helpers.GroupBy(tokens, func(token refresh.RefreshTokenDto, _ int, _ []refresh.RefreshTokenDto) (string, refresh.RefreshTokenDto) {
			return token.UserId, token
		})

		for userId, tokens := range userTokens {
			encodedToken, err := dbredis.EncodeGob(tokens)
			if err != nil {
				return nil, err
			}
			userTokensEncoded[userId] = string(encodedToken)
		}
	}
	return userTokensEncoded, nil
}
