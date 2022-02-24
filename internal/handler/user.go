package handler

import (
	"net/http"

	"github.com/umalmyha/authsrv/internal/business/user"
	"github.com/umalmyha/authsrv/internal/service"
	"github.com/umalmyha/authsrv/pkg/web"
)

type UserHandler struct {
	userSrv *service.UserService
}

func NewUserHandler(userSrv *service.UserService) *UserHandler {
	return &UserHandler{
		userSrv: userSrv,
	}
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) error {
	var nu user.NewUserDto
	if err := web.JsonReqBody(r, &nu); err != nil {
		return err
	}
	return h.userSrv.CreateUser(r.Context(), nu)
}
