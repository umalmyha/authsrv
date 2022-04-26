package role

import (
	"container/list"
	"fmt"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/pkg/ddd/errors"
	"github.com/umalmyha/authsrv/pkg/helpers"
)

type isExsitingRoleNameFn func(string) (bool, error)

func FromNewRoleDto(dto NewRoleDto, existFn isExsitingRoleNameFn) (*Role, error) {
	validation := errors.NewValidation()
	if dto.Name == "" {
		validation.AddViolation(
			errors.NewInvariantViolation("role name can not be empty", "name", errors.ViolationSeverityErr),
		)
	}

	roleName, err := valueobj.NewSolidString(dto.Name)
	if err != nil {
		validation.AddViolation(
			errors.NewInvariantViolation(err.Error(), "name", errors.ViolationSeverityErr),
		)
	}

	exist, err := existFn(dto.Name)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to check existence of role by name")
	} else if exist {
		validation.AddViolation(
			errors.NewInvariantViolation(
				fmt.Sprintf("role with name %s already exists", dto.Name),
				"name",
				errors.ViolationSeverityErr,
			),
		)
	}

	if validation.HasError() {
		return nil, pkgerrors.Wrap(validation.Err(), "validation failed for role creation")
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
		return nil, pkgerrors.Wrap(err, "failed to build role name from db entry")
	}

	return &Role{
		id:          roleDto.Id,
		name:        name,
		description: valueobj.NewNilStringFromPtr(roleDto.Description),
		scopes:      helpers.ToList(scopesDto),
	}, nil
}
