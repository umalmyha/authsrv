package role

import (
	"container/list"
	"fmt"

	"github.com/google/uuid"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/internal/errors"
)

type isExsitingRoleNameFn func(string) (*Role, error)

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
	} else if exist != nil {
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
