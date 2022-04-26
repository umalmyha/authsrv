package handler

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/umalmyha/authsrv/internal/business/refresh"
	"github.com/umalmyha/authsrv/internal/business/user"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/internal/infra/service"
	"github.com/umalmyha/authsrv/pkg/web/request"
	"github.com/umalmyha/authsrv/pkg/web/response"
)

type AuthHandler struct {
	authSrv *service.AuthService
	rfrCfg  valueobj.RefreshTokenConfig
}

func NewAuthHandler(authSrv *service.AuthService, rfrCfg valueobj.RefreshTokenConfig) *AuthHandler {
	return &AuthHandler{
		authSrv: authSrv,
		rfrCfg:  rfrCfg,
	}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) error {
	var nu user.NewUserDto
	if err := request.JsonReqBody(r, nu); err != nil {
		return err
	}

	return h.authSrv.Signup(r.Context(), nu)
}

func (h *AuthHandler) Signin(w http.ResponseWriter, r *http.Request) error {
	var signin user.SigninDto
	if err := request.JsonReqBody(r, signin); err != nil {
		return err
	}

	username, password, err := request.BasicAuth(r)
	if err != nil {
		return err
	}
	signin.Username = username
	signin.Password = password

	refreshCookie := h.rfrCfg.CookieName()
	if request.GetCookieValue(r, refreshCookie) != "" {
		return errors.New("refresh token cookie is set, logout first or refresh session")
	}

	jwt, rfrToken, err := h.authSrv.Signin(r.Context(), signin)
	if err != nil {
		return err
	}

	return h.respondTokens(w, jwt, rfrToken)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	// TODO: Find user: inject to context from JWT in middleware?
	var logout user.LogoutDto
	if err := request.JsonReqBody(r, &logout); err != nil {
		return err
	}

	rfrTokenCookie := h.rfrCfg.CookieName()

	refreshTokenId := request.GetCookieValue(r, rfrTokenCookie)
	logout.RefreshTokenId = refreshTokenId

	if err := h.authSrv.Logout(r.Context(), logout); err != nil {
		return err
	}

	response.DeleteCookie(r, w, rfrTokenCookie)
	return nil
}

func (h *AuthHandler) RefreshSession(w http.ResponseWriter, r *http.Request) error {
	// TODO: Find user: inject to context from JWT in middleware?
	refreshTokenId := request.GetCookieValue(r, h.rfrCfg.CookieName())
	if refreshTokenId == "" {
		return errors.New("refresh token is not provided")
	}

	var rfr user.RefreshDto
	if err := request.JsonReqBody(r, &rfr); err != nil {
		return err
	}
	rfr.RefreshTokenId = refreshTokenId

	// TODO: Think of allowed errors
	jwt, rfrToken, err := h.authSrv.RefreshSession(r.Context(), rfr)
	if err != nil {
		return err
	}

	return h.respondTokens(w, jwt, rfrToken)
}

func (h *AuthHandler) respondTokens(w http.ResponseWriter, jwt valueobj.Jwt, rfrToken *refresh.RefreshToken) error {
	cookie := &http.Cookie{
		Name:     h.rfrCfg.CookieName(),
		Value:    rfrToken.Id(),
		MaxAge:   rfrToken.UnixExpiresIn(),
		HttpOnly: true,
	}
	response.SetCookie(w, cookie)

	signinData := struct {
		AccessToken string `json:"accessToken"`
		ExpiresAt   int64  `json:"expiresAt"`
		TokenType   string `json:"tokenType"`
	}{
		AccessToken: jwt.String(),
		ExpiresAt:   jwt.ExpiresAt(),
		TokenType:   jwt.TokenType(),
	}

	if err := response.RespondJson(w, http.StatusOK, signinData); err != nil {
		return err
	}

	return nil
}
