package role

import (
	"container/list"
	"fmt"

	"github.com/google/uuid"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/internal/errors"
	"github.com/umalmyha/authsrv/pkg/helpers"
)

type isExsitingRoleNameFn func(string) (bool, error)

func FromNewRoleDto(dto NewRoleDto, existFn isExsitingRoleNameFn) (*Role, error) {
	validation := errors.NewValidationResult()
	if dto.Name == "" {
		validation.Add(
			errors.NewInvariantViolationError("role name can not be empty", "name"),
		)
	}

	roleName, err := valueobj.NewSolidString(dto.Name)
	if err != nil {
		validation.Add(
			errors.NewInvariantViolationError(err.Error(), "name"),
		)
	}

	exist, err := existFn(dto.Name)
	if err != nil {
		return nil, err
	} else if exist {
		validation.Add(
			errors.NewInvariantViolationError(fmt.Sprintf("role with name %s already exists", dto.Name), "name"),
		)
	}

	if validation.Failed() {
		return nil, validation.Error()
	}

	return &Role{
		id:          uuid.NewString(),
		name:        roleName,
		description: valueobj.NewNilStringFromPtr(dto.Description),
		scopes:      list.New(),
	}, nil
}

func fromDbDtos(roleDto RoleDto, scopesDto []ScopeAssignmentDto) (*Role, error) {
	name, err := valueobj.NewSolidString(roleDto.Name)
	if err != nil {
		return nil, err
	}

	return &Role{
		id:          roleDto.Id,
		name:        name,
		description: valueobj.NewNilStringFromPtr(roleDto.Description),
		scopes:      helpers.ToList(scopesDto),
	}, nil
}
