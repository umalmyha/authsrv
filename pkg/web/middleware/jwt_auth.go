package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	webErrs "github.com/umalmyha/authsrv/pkg/web/errors"
	"github.com/umalmyha/authsrv/pkg/web/request"
)

type JwtParserFn func(string) (JwtAuthProvider, error)

type JwtAuthProvider interface {
	Subject() string
	Roles() []string
	Scopes() []string
}

type ctxRolesKey string
type ctxScopesKey string
type ctxUsernameKey string

const CtxRoles ctxRolesKey = "user-roles"
const CtxScopes ctxScopesKey = "user-scopes"
const CtxUsername ctxUsernameKey = "username"

func JwtAuthentication(parserFn JwtParserFn) MiddlewareFn {
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

			jwtAuth, err := parserFn(parts[1])
			if err != nil {
				return errors.Wrapf(webErrs.HttpUnauthorizedErr, "error occurred on parsing token - %v", err)
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, CtxUsername, jwtAuth.Subject())
			ctx = context.WithValue(ctx, CtxRoles, jwtAuth.Roles())
			ctx = context.WithValue(ctx, CtxScopes, jwtAuth.Scopes())

			return nextFn(w, r.WithContext(ctx))
		}
	}
}

func HasRoles(roles ...string) MiddlewareFn {
	return func(nextFn HttpHandlerFn) HttpHandlerFn {
		return func(w http.ResponseWriter, r *http.Request) error {
			if len(roles) > 0 {
				ctx := r.Context()
				ctxRoles, ok := ctx.Value(CtxRoles).([]string)
				if !ok {
					return errors.Wrap(webErrs.HttpInternalServerErr, "roles are missing in context, is jwt authentication middleware were applied?")
				}

				m := findMissingPrivileges(roles, ctxRoles)
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
				ctxScopes, ok := ctx.Value(CtxScopes).([]string)
				if !ok {
					return errors.Wrap(webErrs.HttpForbiddenErr, "scopes are missing in context, is jwt authentication middleware were applied?")
				}

				m := findMissingPrivileges(scopes, ctxScopes)
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
