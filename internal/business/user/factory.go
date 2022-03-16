package user

import (
	"container/list"
	"fmt"

	"github.com/google/uuid"
	"github.com/umalmyha/authsrv/internal/business/refresh"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/internal/errors"
	"github.com/umalmyha/authsrv/pkg/helpers"
)

type isExistingUsernameFn func(string) (bool, error)

func FromNewUserDto(dto NewUserDto, cfg valueobj.PasswordConfig, existFn isExistingUsernameFn) (*User, error) {
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

	password, err := valueobj.GeneratePassword(dto.Password, cfg)
	if err != nil {
		validation.Add(
			errors.NewInvariantViolationError(
				fmt.Sprintf(err.Error()),
				"password",
			),
		)
	}

	if validation.Failed() {
		return nil, validation.Error()
	}

	return &User{
		id:          uuid.NewString(),
		username:    username,
		email:       email,
		password:    password,
		isSuperuser: dto.IsSuperuser,
		firstName:   valueobj.NewNilStringFromPtr(dto.FirstName),
		lastName:    valueobj.NewNilStringFromPtr(dto.LastName),
		middleName:  valueobj.NewNilStringFromPtr(dto.MiddleName),
		roles:       list.New(),
		tokens:      list.New(),
		auth:        valueobj.NewUserAuth(nil, nil),
	}, nil
}

func fromDbDtos(user UserDto, roleIds []valueobj.RoleId, tokens []*refresh.RefreshToken, auth valueobj.UserAuth) (*User, error) {
	username, err := valueobj.NewSolidString(user.Username)
	if err != nil {
		return nil, err
	}

	email, err := valueobj.NewNilEmailFromPtr(user.Email)
	if err != nil {
		return nil, err
	}

	return &User{
		id:          user.Id,
		username:    username,
		email:       email,
		password:    valueobj.PasswordFromHash(user.Password),
		isSuperuser: user.IsSuperuser,
		firstName:   valueobj.NewNilStringFromPtr(user.FirstName),
		lastName:    valueobj.NewNilStringFromPtr(user.LastName),
		middleName:  valueobj.NewNilStringFromPtr(user.MiddleName),
		roles:       helpers.ToList(roleIds),
		tokens:      helpers.ToList(tokens),
		auth:        auth,
	}, nil
}
