package role

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/umalmyha/authsrv/pkg/database/rdb"
)

type RoleDao struct {
	ec sqlx.ExtContext
}

func NewRoleDao(ec sqlx.ExtContext) *RoleDao {
	return &RoleDao{
		ec: ec,
	}
}

func (dao *RoleDao) CreateMulti(ctx context.Context, roles []RoleDto) error {
	applier := func(role RoleDto) []any {
		return []any{role.Id, role.Name, role.Description}
	}

	q, params, err := rdb.BulkInsertQuery("ROLES", []string{"ID", "NAME", "DESCRIPTION"}, roles, applier)
	if err != nil {
		return errors.Wrap(err, "failed to build bulk insert SQL query for roles creation")
	}

	if _, err := dao.ec.ExecContext(ctx, q, params...); err != nil {
		return errors.Wrap(err, "failed to crate roles")
	}

	return nil
}

func (dao *RoleDao) DeleteWhereIdsIn(ctx context.Context, ids []string) error {
	inRange, params, err := rdb.WhereIn(ids)
	if err != nil {
		return errors.Wrap(err, "failed to generate SQL where clause for roles deletion")
	}

	q := fmt.Sprintf("DELETE FROM ROLES WHERE ID IN %s", inRange)
	if _, err := dao.ec.ExecContext(ctx, q, params...); err != nil {
		return errors.Wrap(err, "failed to delete roles")
	}

	return nil
}

func (dao *RoleDao) Update(ctx context.Context, role RoleDto) error {
	q := "UPDATE ROLES SET DESCRIPTION = $1 WHERE ID = $2"
	if _, err := dao.ec.ExecContext(ctx, q, role.Description, role.Id); err != nil {
		return errors.Wrap(err, "failed to update roles")
	}
	return nil
}

func (dao *RoleDao) FindByName(ctx context.Context, name string) (RoleDto, error) {
	var r RoleDto
	q := "SELECT ID, NAME, DESCRIPTION FROM ROLES WHERE NAME = $1"
	if err := sqlx.GetContext(ctx, dao.ec, &r, q, name); err != nil {
		return r, errors.Wrap(err, "failed to find role by name")
	}
	return r, nil
}

func (dao *RoleDao) FindById(ctx context.Context, id string) (RoleDto, error) {
	var r RoleDto
	q := "SELECT ID, NAME, DESCRIPTION FROM ROLES WHERE ID = $1"
	if err := sqlx.GetContext(ctx, dao.ec, &r, q, id); err != nil {
		return r, errors.Wrap(err, "failed to find role by id")
	}
	return r, nil
}

type ScopeAssignmentDao struct {
	ec sqlx.ExtContext
}

func NewScopeAssignmentDao(ec sqlx.ExtContext) *ScopeAssignmentDao {
	return &ScopeAssignmentDao{
		ec: ec,
	}
}

func (dao *ScopeAssignmentDao) CreateMulti(ctx context.Context, scopes []ScopeAssignmentDto) error {
	applier := func(scope ScopeAssignmentDto) []any {
		return []any{scope.RoleId, scope.ScopeId}
	}

	q, params, err := rdb.BulkInsertQuery("ROLES_SCOPES", []string{"ROLE_ID", "SCOPE_ID"}, scopes, applier)
	if err != nil {
		return errors.Wrap(err, "failed to build bulk insert SQL query for scope assignments creation")
	}

	if _, err := dao.ec.ExecContext(ctx, q, params...); err != nil {
		return errors.Wrap(err, "failed to create scope assignments")
	}

	return nil
}

func (dao *ScopeAssignmentDao) DeleteByRoleIdAndScopeIdsIn(ctx context.Context, roleId string, scopeIds []string) error {
	inRange, params, err := rdb.WhereIn(scopeIds)
	if err != nil {
		return errors.Wrap(err, "failed to generate SQL where clause for scope assignments deletion")
	}

	params = append(params, roleId)
	q := fmt.Sprintf("DELETE FROM ROLES_SCOPES WHERE SCOPE_ID IN %s AND ROLE_ID = $%d", inRange, len(params))
	if _, err := dao.ec.ExecContext(ctx, q, params...); err != nil {
		return errors.Wrap(err, "failed to delete scope assignments")
	}

	return nil
}
