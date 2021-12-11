package validator

import (
	"fmt"

	"github.com/umalmyha/authsrv/internal/user/dto"
	"github.com/umalmyha/authsrv/pkg/validate"
	"github.com/umalmyha/authsrv/pkg/web"
)

type updateUserValidator struct{}

func ForUpdateUser() *updateUserValidator {
	return &updateUserValidator{}
}

func (v *updateUserValidator) Validate(id string, uu dto.UpdateUser) *web.ValidationResult {
	validation := web.NewValidationResult()

	if ok, _ := validate.UUID(id); !ok {
		validation.Add(
			web.NewValidationMessage(
				`provided user id is not valid uuid`,
				"id",
				web.SeverityError,
			),
		)
	}

	if uu.Email != nil {
		if ok, _ := validate.Email(*uu.Email); !ok {
			validation.Add(
				web.NewValidationMessage(
					fmt.Sprintf("Wrong email provided '%s'. Please, use format myemail@example.com", *uu.Email),
					"email",
					web.SeverityError,
				),
			)
		}
	}

	return validation
}
