package role

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/umalmyha/authsrv/pkg/helpers"
	"github.com/umalmyha/authsrv/pkg/uow"
)

type unitOfWork struct {
	db             *sqlx.DB
	roles          *uow.ChangeSet[RoleDto]
	assignedScopes *uow.ChangeSet[AssignedScopeDto]
}

func NewUnitOfWork(db *sqlx.DB) *unitOfWork {
	return &unitOfWork{
		db:             db,
		roles:          uow.NewChangeSet[RoleDto](),
		assignedScopes: uow.NewChangeSet[AssignedScopeDto](),
	}
}

func (uow *unitOfWork) RegisterClean(role *Role) error {
	uow.roles.Attach(role.ToDto())
	uow.assignedScopes.AttachRange(role.ScopesDto()...)
	return nil
}

func (uow *unitOfWork) RegisterNew(role *Role) error {
	if err := uow.roles.Add(role.ToDto()); err != nil {
		return err
	}

	if err := uow.assignedScopes.AddRange(role.ScopesDto()...); err != nil {
		return err
	}

	return nil
}

func (uow *unitOfWork) RegisterDeleted(role *Role) error {
	if err := uow.roles.Remove(role.ToDto()); err != nil {
		return err
	}

	if err := uow.assignedScopes.RemoveRange(role.ScopesDto()...); err != nil {
		return err
	}

	return nil
}

func (uow *unitOfWork) RegisterAmended(role *Role) error {
	roleDto := role.ToDto()
	if err := uow.roles.Update(role.ToDto()); err != nil {
		return err
	}

	scopes := role.ScopesDto()
	created, _, deleted := uow.assignedScopes.DeltaWithMatched(scopes, func(scopeDto AssignedScopeDto) bool {
		return scopeDto.RoleId == roleDto.Id
	})

	if err := uow.assignedScopes.AddRange(created...); err != nil {
		return err
	}

	if err := uow.assignedScopes.RemoveRange(deleted...); err != nil {
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

	rolesDao := NewRoleDao(tx)
	scopesDao := NewAssignedScopeDao(tx)

	if rmScopes := uow.assignedScopes.Deleted(); len(rmScopes) > 0 {
		scopesGroup := helpers.GroupBy(rmScopes, func(scope AssignedScopeDto, _ int, _ []AssignedScopeDto) (string, string) {
			return scope.RoleId, scope.ScopeId
		})

		for roleId, scopeIds := range scopesGroup {
			if err := scopesDao.DeleteByRoleIdAndScopeIdsIn(ctx, roleId, scopeIds); err != nil {
				return err
			}
		}
	}

	if deletedRoles := uow.roles.Deleted(); len(deletedRoles) > 0 {
		mapper := func(delRole RoleDto, _ int, _ []RoleDto) string {
			return delRole.Id
		}
		if err := rolesDao.DeleteWhereIdsIn(ctx, helpers.Map(deletedRoles, mapper)); err != nil {
			return err
		}
	}

	if createdRoles := uow.roles.Created(); len(createdRoles) > 0 {
		if err := rolesDao.CreateMulti(ctx, createdRoles); err != nil {
			return err
		}
	}

	if createdScopes := uow.assignedScopes.Created(); len(createdScopes) > 0 {
		if err := scopesDao.CreateMulti(ctx, createdScopes); err != nil {
			return err
		}
	}

	if updatedRoles := uow.roles.Updated(); len(updatedRoles) > 0 {
		for _, updRole := range updatedRoles {
			if err := rolesDao.Update(ctx, updRole); err != nil {
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
	uow.roles.Cleanup()
	uow.assignedScopes.Cleanup()
	return nil
}
