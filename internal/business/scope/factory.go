package scope

import (
	"fmt"

	"github.com/google/uuid"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/internal/errors"
)

type isExsitingScopeNameFn func(string) (bool, error)

func FromNewScopeDto(dto NewScopeDto, existFn isExsitingScopeNameFn) (*Scope, error) {
	validation := errors.NewValidationResult()
	if dto.Name == "" {
		validation.Add(
			errors.NewInvariantViolationError("scope name can not be empty", "name"),
		)
	}

	scopeName, err := valueobj.NewSolidString(dto.Name)
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
			errors.NewInvariantViolationError(fmt.Sprintf("scope with name %s already exists", dto.Name), "name"),
		)
	}

	if validation.Failed() {
		return nil, validation.Error()
	}

	return &Scope{
		id:          uuid.NewString(),
		name:        scopeName,
		description: valueobj.NewNilStringFromPtr(dto.Description),
	}, nil
}
