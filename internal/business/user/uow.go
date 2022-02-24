package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/umalmyha/authsrv/pkg/helpers"
	"github.com/umalmyha/authsrv/pkg/uow"
)

type unitOfWork struct {
	db            *sqlx.DB
	users         *uow.ChangeSet[UserDto]
	assignedRoles *uow.ChangeSet[AssignedRoleDto]
}

func NewUnitOfWork(db *sqlx.DB) *unitOfWork {
	return &unitOfWork{
		db:            db,
		users:         uow.NewChangeSet[UserDto](),
		assignedRoles: uow.NewChangeSet[AssignedRoleDto](),
	}
}

func (uow *unitOfWork) RegisterClean(user *User) error {
	uow.users.Attach(user.ToDto())
	uow.assignedRoles.AttachRange(user.RolesDto()...)
	return nil
}

func (uow *unitOfWork) RegisterNew(user *User) error {
	if err := uow.users.Add(user.ToDto()); err != nil {
		return err
	}

	if err := uow.assignedRoles.AddRange(user.RolesDto()...); err != nil {
		return err
	}
	return nil
}

func (uow *unitOfWork) RegisterDeleted(user *User) error {
	if err := uow.users.Remove(user.ToDto()); err != nil {
		return err
	}

	if err := uow.assignedRoles.RemoveRange(user.RolesDto()...); err != nil {
		return err
	}
	return nil
}

func (uow *unitOfWork) RegisterAmended(user *User) error {
	userDto := user.ToDto()
	if err := uow.users.Update(userDto); err != nil {
		return err
	}

	created, _, deleted := uow.assignedRoles.DeltaWithMatched(user.RolesDto(), func(role AssignedRoleDto) bool {
		return role.UserId == userDto.Id
	})

	if err := uow.assignedRoles.AddRange(created...); err != nil {
		return err
	}

	if err := uow.assignedRoles.RemoveRange(deleted...); err != nil {
		return err
	}

	return nil
}

func (uow *unitOfWork) Flush(ctx context.Context) error {
	tx, err := uow.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	userDao := NewUserDao(tx)
	roleDao := NewAssignedRoleDao(tx)

	if rmRoles := uow.assignedRoles.Deleted(); len(rmRoles) > 0 {
		rolesGroup := helpers.GroupBy(rmRoles, func(role AssignedRoleDto, _ int, _ []AssignedRoleDto) (string, string) {
			return role.UserId, role.UserId
		})

		for userId, rolesId := range rolesGroup {
			if err := roleDao.DeleteByUserIdAndRoleIdsIn(ctx, userId, rolesId); err != nil {
				return err
			}
		}
	}

	if rmUsers := uow.users.Deleted(); len(rmUsers) > 0 {
		mapper := func(user UserDto, _ int, _ []UserDto) string {
			return user.Id
		}
		if err := userDao.DeleteWhereIdsIn(ctx, helpers.Map(rmUsers, mapper)); err != nil {
			return err
		}
	}

	if createdUsers := uow.users.Created(); len(createdUsers) > 0 {
		if err := userDao.CreateMulti(ctx, createdUsers); err != nil {
			return err
		}
	}

	if createdRoles := uow.assignedRoles.Created(); len(createdRoles) > 0 {
		if err := roleDao.CreateMulti(ctx, createdRoles); err != nil {
			return err
		}
	}

	if updatedUsers := uow.users.Updated(); len(updatedUsers) > 0 {
		for _, updUser := range updatedUsers {
			if err := userDao.Update(ctx, updUser); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (uow *unitOfWork) Dispose() error {
	uow.users.Cleanup()
	uow.assignedRoles.Cleanup()
	return nil
}
