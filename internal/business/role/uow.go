package role

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/umalmyha/authsrv/pkg/ddd/uow"
	"github.com/umalmyha/authsrv/pkg/helpers"
)

type unitOfWork struct {
	*uow.SqlxUnitOfWork
	roles          *uow.ChangeSet[RoleDto]
	assignedScopes *uow.ChangeSet[ScopeAssignmentDto]
}

func NewUnitOfWork(db *sqlx.DB) *unitOfWork {
	return &unitOfWork{
		SqlxUnitOfWork: uow.NewSqlxUnitOfWork(db),
		roles:          uow.NewChangeSet[RoleDto](),
		assignedScopes: uow.NewChangeSet[ScopeAssignmentDto](),
	}
}

func (uow *unitOfWork) RegisterClean(role *Role) error {
	uow.roles.Attach(role.ToDto())
	uow.assignedScopes.AttachRange(role.ScopesDto()...)
	return nil
}

func (uow *unitOfWork) RegisterNew(role *Role) error {
	if err := uow.roles.Add(role.ToDto()); err != nil {
		return errors.Wrap(err, "failed to add role DTO to changeset")
	}

	if err := uow.assignedScopes.AddRange(role.ScopesDto()...); err != nil {
		return errors.Wrap(err, "failed to add scope assignments DTOs to changeset")
	}

	return nil
}

func (uow *unitOfWork) RegisterDeleted(role *Role) error {
	if err := uow.roles.Remove(role.ToDto()); err != nil {
		return errors.Wrap(err, "failed to delete role DTO in changeset")
	}

	if err := uow.assignedScopes.RemoveRange(role.ScopesDto()...); err != nil {
		return errors.Wrap(err, "failed to delete scope assignments DTOs in changeset")
	}

	return nil
}

func (uow *unitOfWork) RegisterAmended(role *Role) error {
	roleDto := role.ToDto()
	if err := uow.roles.Update(role.ToDto()); err != nil {
		return errors.Wrap(err, "failed to update role DTO in changeset")
	}

	scopes := role.ScopesDto()
	created, _, deleted := uow.assignedScopes.DeltaWithMatched(scopes, func(scopeDto ScopeAssignmentDto) bool {
		return scopeDto.RoleId == roleDto.Id
	})

	if err := uow.assignedScopes.AddRange(created...); err != nil {
		return errors.Wrap(err, "failed to add scope assignments DTOs to changeset")
	}

	if err := uow.assignedScopes.RemoveRange(deleted...); err != nil {
		return errors.Wrap(err, "failed to delete scope assignments DTOs in changeset")
	}

	return nil
}

func (uow *unitOfWork) Flush(ctx context.Context) error {
	tx, err := uow.Tx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to open transaction")
	}
	defer tx.Rollback()

	rolesDao := NewRoleDao(tx)
	scopesDao := NewScopeAssignmentDao(tx)

	if rmScopes := uow.assignedScopes.Deleted(); len(rmScopes) > 0 {
		scopesGroup := helpers.GroupBy(rmScopes, func(scope ScopeAssignmentDto, _ int, _ []ScopeAssignmentDto) (string, string) {
			return scope.RoleId, scope.ScopeId
		})

		for roleId, scopeIds := range scopesGroup {
			if err := scopesDao.DeleteByRoleIdAndScopeIdsIn(ctx, roleId, scopeIds); err != nil {
				return errors.Wrap(err, "failed to process scopes assignments deletion")
			}
		}
	}

	if deletedRoles := uow.roles.Deleted(); len(deletedRoles) > 0 {
		mapper := func(delRole RoleDto, _ int, _ []RoleDto) string {
			return delRole.Id
		}
		if err := rolesDao.DeleteWhereIdsIn(ctx, helpers.Map(deletedRoles, mapper)); err != nil {
			return errors.Wrap(err, "failed to process roles deletion")
		}
	}

	if createdRoles := uow.roles.Created(); len(createdRoles) > 0 {
		if err := rolesDao.CreateMulti(ctx, createdRoles); err != nil {
			return errors.Wrap(err, "failed to process roles creation")
		}
	}

	if createdScopes := uow.assignedScopes.Created(); len(createdScopes) > 0 {
		if err := scopesDao.CreateMulti(ctx, createdScopes); err != nil {
			return errors.Wrap(err, "failed to process scopes assignments creation")
		}
	}

	if updatedRoles := uow.roles.Updated(); len(updatedRoles) > 0 {
		for _, updRole := range updatedRoles {
			if err := rolesDao.Update(ctx, updRole); err != nil {
				return errors.Wrap(err, "failed to process roles update")
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}
	return uow.Dispose()
}

func (uow *unitOfWork) Dispose() error {
	uow.roles.Cleanup()
	uow.assignedScopes.Cleanup()
	return nil
}
