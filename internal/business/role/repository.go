package role

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "failed to read role aggregate by id")
	}

	if !role.IsPresent() {
		return nil, fmt.Errorf("role with id %s doesn't exist", id)
	}

	assignedScopes := repo.uow.assignedScopes.Filter(func(dto ScopeAssignmentDto) bool {
		return dto.RoleId == role.Id
	})

	r, err := fromDbDtos(role, assignedScopes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build role aggregate from db DTOs")
	}

	return r, repo.uow.RegisterClean(r)
}

func (repo *Repository) FindByName(ctx context.Context, name string) (*Role, error) {
	notPresentFn := func() (RoleDto, error) {
		return NewRoleDao(repo.uow.ExtContext()).FindByName(ctx, name)
	}

	role, err := repo.uow.roles.Find(func(dto RoleDto) bool { return dto.Name == name }).IfNotPresent(notPresentFn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read role aggregate by name")
	}

	if !role.IsPresent() {
		return nil, fmt.Errorf("role %s doesn't exist", name)
	}

	assignedScopes := repo.uow.assignedScopes.Filter(func(dto ScopeAssignmentDto) bool {
		return dto.RoleId == role.Id
	})

	r, err := fromDbDtos(role, assignedScopes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build role aggregate from db DTOs")
	}

	return r, repo.uow.RegisterClean(r)
}
