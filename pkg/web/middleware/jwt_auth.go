package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	webErrs "github.com/umalmyha/authsrv/pkg/web/errors"
	"github.com/umalmyha/authsrv/pkg/web/request"
)

type JwtValidatorFn func(string) (AuthClaimsProvider, error)

type AuthClaimsProvider interface {
	Username() string
	Roles() []string
	Scopes() []string
}

type ctxClaimsKey string
type ctxUsernameKey string

const CtxClaims ctxClaimsKey = "claims"
const CtxUsername ctxUsernameKey = "username"

func JwtAuthentication(validatorFn JwtValidatorFn) MiddlewareFn {
	return func(nextFn HttpHandlerFn) HttpHandlerFn {
		return func(w http.ResponseWriter, r *http.Request) error {
			h := request.GetHeader(r, "Authorization")
			if h == "" {
				h = request.GetHeader(r, "authorization")
			}

			parts := strings.Split(h, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return errors.Wrap(webErrs.HttpUnauthorizedErr, "incorrect authorization header, expected format 'Bearer <token>'")
			}

			jwtAuth, err := validatorFn(parts[1])
			if err != nil {
				return errors.Wrapf(webErrs.HttpUnauthorizedErr, "error occurred on parsing token - %v", err)
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, CtxUsername, jwtAuth.Username())
			ctx = context.WithValue(ctx, CtxClaims, jwtAuth)

			return nextFn(w, r.WithContext(ctx))
		}
	}
}

func HasRoles(roles ...string) MiddlewareFn {
	return func(nextFn HttpHandlerFn) HttpHandlerFn {
		return func(w http.ResponseWriter, r *http.Request) error {
			if len(roles) > 0 {
				ctx := r.Context()
				claims, ok := ctx.Value(CtxClaims).(AuthClaimsProvider)
				if !ok {
					return errors.Wrap(webErrs.HttpInternalServerErr, "claims are missing in context, is jwt authentication middleware was applied?")
				}

				m := findMissingPrivileges(roles, claims.Roles())
				if len(m) > 0 {
					return errors.Wrapf(webErrs.HttpForbiddenErr, "authorization failed, missing roles %v", m)
				}
			}
			return nextFn(w, r)
		}
	}
}

func HasScopes(scopes ...string) MiddlewareFn {
	return func(nextFn HttpHandlerFn) HttpHandlerFn {
		return func(w http.ResponseWriter, r *http.Request) error {
			if len(scopes) > 0 {
				ctx := r.Context()
				claims, ok := ctx.Value(CtxClaims).(AuthClaimsProvider)
				if !ok {
					return errors.Wrap(webErrs.HttpForbiddenErr, "claims are missing in context, is jwt authentication middleware were applied?")
				}

				m := findMissingPrivileges(scopes, claims.Scopes())
				if len(m) > 0 {
					return errors.Wrapf(webErrs.HttpForbiddenErr, "authorization failed, missing scopes %v", m)
				}
			}
			return nextFn(w, r)
		}
	}
}

func findMissingPrivileges(required []string, actual []string) []string {
	m := make([]string, 0)
	for _, r := range required {
		found := false
		for _, a := range actual {
			if r == a {
				found = true
				break
			}
		}

		if found == false {
			m = append(m, r)
		}
	}
	return m
}
