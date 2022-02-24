package user

import (
	"container/list"
	"fmt"

	"github.com/google/uuid"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/internal/errors"
)

type isExistingUsernameFn func(string) (bool, error)

func FromNewUserDto(dto NewUserDto, existFn isExistingUsernameFn) (*User, error) {
	validation := errors.NewValidationResult()

	if dto.Username == "" {
		validation.Add(
			errors.NewInvariantViolationError("username is mandatory", "username"),
		)
	} else {
		if exist, err := existFn(dto.Username); err != nil {
			return nil, err
		} else if exist {
			validation.Add(
				errors.NewInvariantViolationError(fmt.Sprintf("user with username '%s' already exists", dto.Username), "username"),
			)
		}
	}

	username, err := valueobj.NewSolidString(dto.Username)
	if err != nil {
		validation.Add(
			errors.NewInvariantViolationError(err.Error(), "username"),
		)
	}

	email, err := valueobj.NewNilEmailFromPtr(dto.Email)
	if err != nil {
		validation.Add(
			errors.NewInvariantViolationError(
				fmt.Sprintf("Wrong email provided '%s'. Please, use format myemail@example.com", *dto.Email),
				"email",
			),
		)
	}

	if dto.Password != dto.ConfirmPassword {
		validation.Add(
			errors.NewInvariantViolationError("passwords don't match", "confirmPassword"),
		)
	}

	cfg := valueobj.PasswordConfig{Min: 4, Max: 15, HasDigit: true, HasUppercase: true}
	password, err := valueobj.GeneratePassword(dto.Password, cfg)
	if err != nil {
		validation.Add(
			errors.NewInvariantViolationError(
				fmt.Sprintf("password must has length from %d to %d symbols, has one upper case symbol and has one digit", cfg.Min, cfg.Max),
				"confirmPassword",
			),
		)
	}

	if validation.Failed() {
		return nil, validation.Error()
	}

	return &User{
		id:            uuid.NewString(),
		username:      username,
		email:         email,
		password:      password,
		isSuperuser:   dto.IsSuperuser,
		firstName:     valueobj.NewNilStringFromPtr(dto.FirstName),
		lastName:      valueobj.NewNilStringFromPtr(dto.LastName),
		middleName:    valueobj.NewNilStringFromPtr(dto.MiddleName),
		assignedRoles: list.New(),
	}, nil
}
