package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/umalmyha/authsrv/internal/business/refresh"
	"github.com/umalmyha/authsrv/internal/business/user"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
)

type AuthService struct {
	db         *sqlx.DB
	rdb        *redis.Client
	jwtCfg     valueobj.JwtConfig
	passCfg    valueobj.PasswordConfig
	refreshCfg valueobj.RefreshTokenConfig
}

func NewAuthService(db *sqlx.DB, rdb *redis.Client, jwtCfg valueobj.JwtConfig, rfrCfg valueobj.RefreshTokenConfig, passCfg valueobj.PasswordConfig) *AuthService {
	return &AuthService{
		db:         db,
		rdb:        rdb,
		jwtCfg:     jwtCfg,
		refreshCfg: rfrCfg,
		passCfg:    passCfg,
	}
}

func (srv *AuthService) Signup(ctx context.Context, u user.NewUserDto) error {
	uow := user.NewUnitOfWork(srv.db, srv.rdb)
	repo := user.NewRepository(uow)

	existUsernameFn := func(username string) (bool, error) {
		if _, err := user.NewUserDao(srv.db).FindByUsername(ctx, username); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}

	user, err := user.FromNewUserDto(u, srv.passCfg, existUsernameFn)
	if err != nil {
		return errors.Wrap(err, "failed to create user from DTO")
	}

	if err := repo.Add(user); err != nil {
		return errors.Wrap(err, "failed to create user from DTO")
	}

	return uow.Flush(ctx)
}

func (srv *AuthService) Signin(ctx context.Context, signin user.SigninDto) (valueobj.Jwt, *refresh.RefreshToken, error) {
	var accessToken valueobj.Jwt
	var refreshToken *refresh.RefreshToken

	uow := user.NewUnitOfWork(srv.db, srv.rdb)
	repo := user.NewRepository(uow)

	user, err := repo.FindByUsername(ctx, signin.Username)
	if err != nil {
		return accessToken, refreshToken, errors.Wrap(err, "failed to find user in repository")
	}

	if user == nil {
		return accessToken, refreshToken, errors.Errorf("user %s doesn't exist", signin.Username)
	}

	verified, err := user.VerifyPassword(signin.Password)
	if err != nil {
		return accessToken, refreshToken, errors.Wrap(err, "failed to verify user password")
	}

	if !verified {
		return accessToken, refreshToken, errors.New("password is incorrect")
	}

	issuedAt := time.Now().UTC()

	accessToken, err = user.GenerateJwt(issuedAt, srv.jwtCfg)
	if err != nil {
		return accessToken, refreshToken, errors.Wrap(err, "failed to generate access token")
	}

	refreshToken, err = user.GenerateRefreshToken(signin.Fingerprint, issuedAt, srv.refreshCfg)
	if err != nil {
		return accessToken, refreshToken, errors.Wrap(err, "failed to generate refresh token")
	}

	if err := repo.Update(user); err != nil {
		return accessToken, refreshToken, errors.Wrap(err, "failed to update user in repository")
	}

	if err := uow.Flush(ctx); err != nil {
		return accessToken, refreshToken, errors.Wrap(err, "failed to flush changes")
	}

	return accessToken, refreshToken, nil
}

func (srv *AuthService) Logout(ctx context.Context, logout user.LogoutDto) error {
	uow := user.NewUnitOfWork(srv.db, srv.rdb)
	defer uow.Dispose()

	repo := user.NewRepository(uow)
	user, err := repo.FindByUsername(ctx, logout.Username)
	if err != nil {
		return errors.Wrap(err, "failed to find user in repository")
	}

	if user == nil {
		return errors.Errorf("user %s doesn't exist", logout.Username)
	}

	return user.DiscardRefreshToken(logout)
}

func (srv *AuthService) RefreshSession(ctx context.Context, rfr user.RefreshDto) (valueobj.Jwt, error) {
	uow := user.NewUnitOfWork(srv.db, srv.rdb)
	repo := user.NewRepository(uow)

	usr, err := repo.FindByUsername(ctx, rfr.Username)
	if err != nil {
		return valueobj.Jwt{}, errors.Wrap(err, "failed to find user in repository")
	}

	if usr == nil {
		return valueobj.Jwt{}, errors.Errorf("user %s doesn't exist", rfr.Username)
	}

	now := time.Now().UTC()
	jwt, err := usr.RefreshSession(rfr, now, srv.jwtCfg)
	if err != nil && !errors.Is(err, refresh.RefreshTokenExpiredErr) {
		return jwt, errors.Wrap(err, "failed to refresh session")
	}

	if flushErr := uow.Flush(ctx); flushErr != nil {
		return valueobj.Jwt{}, flushErr
	}

	return jwt, err
}
