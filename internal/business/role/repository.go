package role

import (
	"context"
	"fmt"
)

type Repository struct {
	uow *unitOfWork
}

func NewRepository(u *unitOfWork) *Repository {
	return &Repository{
		uow: u,
	}
}

func (repo *Repository) Add(role *Role) error {
	return repo.uow.RegisterNew(role)
}

func (repo *Repository) Update(role *Role) error {
	return repo.uow.RegisterAmended(role)
}

func (repo *Repository) FindById(ctx context.Context, id string) (*Role, error) {
	notPresentFn := func() (RoleDto, error) {
		return NewRoleDao(repo.uow.ExtContext()).FindById(ctx, id)
	}

	role, err := repo.uow.roles.FindByKey(id).IfNotPresent(notPresentFn)
	if err != nil {
		return nil, err
	}

	if !role.IsPresent() {
		return nil, fmt.Errorf("role with id %s doesn't exist", id)
	}

	assignedScopes := repo.uow.assignedScopes.Filter(func(dto ScopeAssignmentDto) bool {
		return dto.RoleId == role.Id
	})

	return fromDbDtos(role, assignedScopes)
}

func (repo *Repository) FindByName(ctx context.Context, name string) (*Role, error) {
	notPresentFn := func() (RoleDto, error) {
		return NewRoleDao(repo.uow.ExtContext()).FindByName(ctx, name)
	}

	role, err := repo.uow.roles.Find(func(dto RoleDto) bool {
		return dto.Name == name
	}).IfNotPresent(notPresentFn)
	if err != nil {
		return nil, err
	}

	if !role.IsPresent() {
		return nil, fmt.Errorf("role %s doesn't exist", name)
	}

	assignedScopes := repo.uow.assignedScopes.Filter(func(dto ScopeAssignmentDto) bool {
		return dto.RoleId == role.Id
	})

	return fromDbDtos(role, assignedScopes)
}
