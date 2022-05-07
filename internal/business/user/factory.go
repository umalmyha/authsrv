package user

import (
	"container/list"
	"fmt"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
	"github.com/umalmyha/authsrv/internal/business/refresh"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/pkg/errors"
	"github.com/umalmyha/authsrv/pkg/helpers"
)

type isExistingUsernameFn func(string) (bool, error)

func FromNewUserDto(dto NewUserDto, cfg valueobj.PasswordConfig, existFn isExistingUsernameFn) (*User, error) {
	validation := errors.NewValidation()

	if dto.Username == "" {
		validation.Add(
			errors.NewBusinessErr("username is mandatory", "username", errors.ViolationSeverityErr, errors.CodeValidationFailed),
		)
	} else {
		if exist, err := existFn(dto.Username); err != nil {
			return nil, pkgerrors.Wrap(err, "failed to check user existence")
		} else if exist {
			validation.Add(
				errors.NewBusinessErr(
					fmt.Sprintf("user with username '%s' already exists", dto.Username),
					"username",
					errors.ViolationSeverityErr,
					errors.CodeValidationFailed,
				),
			)
		}
	}

	username, err := valueobj.NewSolidString(dto.Username)
	if err != nil {
		validation.Add(
			errors.NewBusinessErr(err.Error(), "username", errors.ViolationSeverityErr, errors.CodeValidationFailed),
		)
	}

	email, err := valueobj.NewNilEmailFromPtr(dto.Email)
	if err != nil {
		validation.Add(
			errors.NewBusinessErr(
				fmt.Sprintf("Wrong email provided '%s'. Please, use format myemail@example.com", *dto.Email),
				"email",
				errors.ViolationSeverityErr,
				errors.CodeValidationFailed,
			),
		)
	}

	if dto.Password != dto.ConfirmPassword {
		validation.Add(
			errors.NewBusinessErr("passwords don't match", "confirmPassword", errors.ViolationSeverityErr, errors.CodeValidationFailed),
		)
	}

	password, err := valueobj.GeneratePassword(dto.Password, cfg)
	if err != nil {
		validation.Add(
			errors.NewBusinessErr(
				fmt.Sprintf(err.Error()),
				"password",
				errors.ViolationSeverityErr,
				errors.CodeValidationFailed,
			),
		)
	}

	if validation.HasError() {
		return nil, pkgerrors.Wrap(validation.RaiseValidationErr(errors.ViolationSeverityErr), "validation failed on user creation")
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
		return nil, pkgerrors.Wrap(err, "failed to build username")
	}

	email, err := valueobj.NewNilEmailFromPtr(user.Email)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to build user email")
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
