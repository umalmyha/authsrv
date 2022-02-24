package user

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/umalmyha/authsrv/pkg/database/rdb"
)

type UserDao struct {
	ec sqlx.ExtContext
}

func NewUserDao(ec sqlx.ExtContext) *UserDao {
	return &UserDao{
		ec: ec,
	}
}

func (dao *UserDao) CreateMulti(ctx context.Context, users []UserDto) error {
	cols := []string{
		"ID",
		"USERNAME",
		"EMAIL",
		"PASSWORD_HASH",
		"IS_SUPERUSER",
		"FIRST_NAME",
		"LAST_NAME",
		"MIDDLE_NAME",
	}

	applier := func(user UserDto) []any {
		return []any{
			user.Id,
			user.Username,
			user.Email,
			user.Password,
			user.IsSuperuser,
			user.FirstName,
			user.LastName,
			user.MiddleName,
		}
	}

	q, params, err := rdb.BulkInsertQuery("USERS", cols, users, applier)
	if err != nil {
		return err
	}

	if _, err := dao.ec.ExecContext(ctx, q, params...); err != nil {
		return err
	}

	return nil
}

func (dao *UserDao) DeleteWhereIdsIn(ctx context.Context, ids []string) error {
	inRange, params, err := rdb.WhereIn(ids)
	if err != nil {
		return err
	}

	q := fmt.Sprintf("DELETE FROM USERS WHERE ID IN %s", inRange)
	if _, err := dao.ec.ExecContext(ctx, q, params...); err != nil {
		return err
	}

	return nil
}

func (dao *UserDao) Update(ctx context.Context, user UserDto) error {
	q := `UPDATE USERS SET
		EMAIL = $1,
		FIRST_NAME = $2,
		LAST_NAME = $3,
		MIDDLE_NAME = $4,
		IS_SUPERUSER = $5,
		PASSWORD_HASH = $6 WHERE ID = $7`

	params := []any{user.Email, user.FirstName, user.LastName, user.MiddleName, user.IsSuperuser, user.Password, user.Id}
	if _, err := dao.ec.ExecContext(ctx, q, params...); err != nil {
		return err
	}

	return nil
}

func (dao *UserDao) FindByUsername(ctx context.Context, username string) (UserDto, error) {
	var user UserDto
	q := "SELECT * FROM USERS WHERE USERNAME = $1 LIMIT 1"
	if err := sqlx.GetContext(ctx, dao.ec, &user, q, username); err != nil {
		return user, err
	}
	return user, nil
}

type AssignedRoleDao struct {
	ec sqlx.ExtContext
}

func NewAssignedRoleDao(ec sqlx.ExtContext) *AssignedRoleDao {
	return &AssignedRoleDao{
		ec: ec,
	}
}

func (dao *AssignedRoleDao) CreateMulti(ctx context.Context, roles []AssignedRoleDto) error {
	applier := func(role AssignedRoleDto) []any {
		return []any{role.UserId, role.UserId}
	}

	q, params, err := rdb.BulkInsertQuery("USER_ROLES", []string{"USER_ID", "ROLE_ID"}, roles, applier)
	if err != nil {
		return err
	}

	if _, err := dao.ec.ExecContext(ctx, q, params...); err != nil {
		return err
	}

	return nil
}

func (dao *AssignedRoleDao) DeleteByUserIdAndRoleIdsIn(ctx context.Context, userId string, roleIds []string) error {
	inRange, params, err := rdb.WhereIn(roleIds)
	if err != nil {
		return err
	}

	params = append(params, userId)
	q := fmt.Sprintf("DELETE FROM USER_ROLES WHERE ROLE_ID IN %s AND USER_ID = $%d", inRange, len(params))
	if _, err := dao.ec.ExecContext(ctx, q, params...); err != nil {
		return err
	}

	return nil
}