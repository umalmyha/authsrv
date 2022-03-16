package role

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
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
		return err
	}

	if _, err := dao.ec.ExecContext(ctx, q, params...); err != nil {
		return err
	}

	return nil
}

func (dao *RoleDao) DeleteWhereIdsIn(ctx context.Context, ids []string) error {
	inRange, params, err := rdb.WhereIn(ids)
	if err != nil {
		return err
	}

	q := fmt.Sprintf("DELETE FROM ROLES WHERE ID IN %s", inRange)
	if _, err := dao.ec.ExecContext(ctx, q, params...); err != nil {
		return err
	}

	return nil
}

func (dao *RoleDao) Update(ctx context.Context, role RoleDto) error {
	q := "UPDATE ROLES SET DESCRIPTION = $1 WHERE ID = $2"
	if _, err := dao.ec.ExecContext(ctx, q, role.Description, role.Id); err != nil {
		return err
	}
	return nil
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
		return err
	}

	if _, err := dao.ec.ExecContext(ctx, q, params...); err != nil {
		return err
	}

	return nil
}

func (dao *ScopeAssignmentDao) DeleteByRoleIdAndScopeIdsIn(ctx context.Context, roleId string, scopeIds []string) error {
	inRange, params, err := rdb.WhereIn(scopeIds)
	if err != nil {
		return err
	}

	params = append(params, roleId)
	q := fmt.Sprintf("DELETE FROM ROLES_SCOPES WHERE SCOPE_ID IN %s AND ROLE_ID = $%d", inRange, len(params))
	if _, err := dao.ec.ExecContext(ctx, q, params...); err != nil {
		return err
	}

	return nil
}
