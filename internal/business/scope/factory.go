package scope

import (
	"fmt"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/pkg/ddd/errors"
)

type isExistingScopeNameFn func(string) (bool, error)

func FromNewScopeDto(dto NewScopeDto, existFn isExistingScopeNameFn) (*Scope, error) {
	validation := errors.NewValidation()
	if dto.Name == "" {
		validation.AddViolation(
			errors.NewInvariantViolation("scope name can not be empty", "name", errors.ViolationSeverityErr),
		)
	}

	scopeName, err := valueobj.NewSolidString(dto.Name)
	if err != nil {
		validation.AddViolation(
			errors.NewInvariantViolation(err.Error(), "name", errors.ViolationSeverityErr),
		)
	}

	exist, err := existFn(dto.Name)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "faield to check scope existence")
	} else if exist {
		validation.AddViolation(
			errors.NewInvariantViolation(fmt.Sprintf("scope with name %s already exists", dto.Name), "name", errors.ViolationSeverityErr),
		)
	}

	if validation.HasError() {
		return nil, pkgerrors.Wrap(validation.Err(), "validation failed for scope creation")
	}

	return &Scope{
		id:          uuid.NewString(),
		name:        scopeName,
		description: valueobj.NewNilStringFromPtr(dto.Description),
	}, nil
}
