package validator

import (
	"context"
	"fmt"

	"github.com/umalmyha/authsrv/internal/user/dto"
	"github.com/umalmyha/authsrv/internal/user/store"
	"github.com/umalmyha/authsrv/pkg/validate"
	"github.com/umalmyha/authsrv/pkg/web"
)

type newUserValidator struct {
	userStore *store.Store
}

func ForNewUser(ustore *store.Store) *newUserValidator {
	return &newUserValidator{
		userStore: ustore,
	}
}

func (v *newUserValidator) Validate(ctx context.Context, nu dto.NewUser) *web.ValidationResult {
	validation := web.NewValidationResult()

	if nu.Username == "" {
		validation.Add(
			web.NewValidationMessage(
				"username is mandatory",
				"username",
				web.SeverityError,
			),
		)
	} else {
		var existingUser store.User
		if err := v.userStore.ByUsername(ctx, &existingUser, nu.Username)(v.userStore.DB()); err == nil {
			validation.Add(
				web.NewValidationMessage(
					fmt.Sprintf("user with username %s already exists", nu.Username),
					"username",
					web.SeverityError,
				),
			)
		}
	}

	if nu.Email != nil {
		if ok, _ := validate.Email(*nu.Email); !ok {
			validation.Add(
				web.NewValidationMessage(
					fmt.Sprintf("Wrong email provided '%s'. Please, use format myemail@example.com", *nu.Email),
					"email",
					web.SeverityError,
				),
			)
		}
	}

	if nu.Password == "" {
		validation.Add(
			web.NewValidationMessage(
				"password can't be initial",
				"password",
				web.SeverityError,
			),
		)
	}

	if nu.Password != nu.ConfirmPassword {
		validation.Add(
			web.NewValidationMessage(
				"passwords don't match",
				"confirmPassword",
				web.SeverityError,
			),
		)
	}

	return validation
}
