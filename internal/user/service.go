package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/umalmyha/authsrv/internal/user/dto"
	"github.com/umalmyha/authsrv/internal/user/store"
	"github.com/umalmyha/authsrv/internal/user/validator"
	"github.com/umalmyha/authsrv/pkg/database"
	"github.com/umalmyha/authsrv/pkg/validate"
	"github.com/umalmyha/authsrv/pkg/web"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	logger    *zap.SugaredLogger
	userStore *store.Store
}

func Service(logger *zap.SugaredLogger, userStore *store.Store) *service {
	return &service{
		logger:    logger,
		userStore: userStore,
	}
}

func (s *service) CreateUser(ctx context.Context, nu dto.NewUser) (*dto.User, error) {
	validation := validator.ForNewUser(store.NewStore(s.userStore.DB())).Validate(ctx, nu)
	if !validation.Valid() {
		s.logger.Info("validation failed for user creation")
		return nil, web.NewValidationError(validation.Messages()...)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("password generation failed: %s", err.Error())
		return nil, err
	}

	user := store.User{
		Id:           uuid.NewString(),
		Username:     nu.Username,
		Email:        database.NewNullString(nu.Email),
		PasswordHash: string(hash),
		IsSuperuser:  nu.IsSuperuser,
		FirstName:    database.NewNullString(nu.FirstName),
		LastName:     database.NewNullString(nu.LastName),
		MiddleName:   database.NewNullString(nu.MiddleName),
	}

	if err = s.userStore.Create(ctx, user)(s.userStore.DB()); err != nil {
		s.logger.Errorf("error occurred on reading user from database: %s", err.Error())
		return nil, err
	}

	userCreated := dto.UserDtoFromStore(user)
	return &userCreated, nil
}

func (s *service) AllUsers(ctx context.Context) ([]dto.User, error) {
	users := make([]store.User, 0)
	err := s.userStore.All(ctx, &users)(s.userStore.DB())
	if err != nil {
		s.logger.Errorf("error occurred on reading users from database: %s", err.Error())
		return nil, err
	}

	result := make([]dto.User, 0)
	for _, user := range users {
		result = append(result, dto.UserDtoFromStore(user))
	}

	return result, nil
}

func (s *service) GetUser(ctx context.Context, id string) (*dto.User, error) {
	if ok, _ := validate.UUID(id); !ok {
		return nil, web.NewValidationError(
			web.NewValidationMessage(
				`provided user id is not valid uuid`,
				"id",
				web.SeverityError,
			),
		)
	}

	var user store.User
	if err := s.userStore.ById(ctx, &user, id)(s.userStore.DB()); err != nil {
		s.logger.Errorf("error occurred on reading user with id %s from database: %s", id, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, web.NewNotFoundError(fmt.Sprintf("user with id %s not found", id))
		}
		return nil, err
	}

	result := dto.UserDtoFromStore(user)
	return &result, nil
}

func (s *service) UpdateUser(ctx context.Context, id string, uu dto.UpdateUser) (*dto.User, error) {
	validation := validator.ForUpdateUser().Validate(id, uu)
	if !validation.Valid() {
		return nil, web.NewValidationError(validation.Messages()...)
	}

	var user store.User
	if err := s.userStore.ById(ctx, &user, id)(s.userStore.DB()); err != nil {
		s.logger.Errorf("error occurred on reading user with id %s from database: %s", id, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, web.NewNotFoundError(fmt.Sprintf("user with id %s not found", id))
		}
		return nil, err
	}

	user.Email = database.NewNullString(uu.Email)
	user.FirstName = database.NewNullString(uu.FirstName)
	user.LastName = database.NewNullString(uu.LastName)
	user.MiddleName = database.NewNullString(uu.MiddleName)

	if err := s.userStore.Update(ctx, user)(s.userStore.DB()); err != nil {
		s.logger.Errorf("error occurred on update of user with id %s: %s", id, err.Error())
		return nil, err
	}

	result := dto.UserDtoFromStore(user)
	return &result, nil
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	_, err := s.GetUser(ctx, id)
	if err != nil {
		return err
	}

	if err := s.userStore.Delete(ctx, id)(s.userStore.DB()); err != nil {
		s.logger.Errorf("error occurred on delete user with id %s: %s", id, err.Error())
		return err
	}

	return nil
}
