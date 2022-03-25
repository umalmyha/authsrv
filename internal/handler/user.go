package handler

import (
	"net/http"

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

func (h *UserHandler) AssignRole(w http.ResponseWriter, r *http.Request) error {
	assingment := struct {
		Username string `json:"username"`
		RoleName string `json:"role"`
	}{}

	if err := web.JsonReqBody(r, &assingment); err != nil {
		return err
	}
	return h.userSrv.AssignRole(r.Context(), assingment.Username, assingment.RoleName)
}

func (h *UserHandler) UnassignRole(w http.ResponseWriter, r *http.Request) error {
	assingment := struct {
		Username string `json:"username"`
		RoleName string `json:"role"`
	}{}

	if err := web.JsonReqBody(r, &assingment); err != nil {
		return err
	}
	return h.userSrv.UnassignRole(r.Context(), assingment.Username, assingment.RoleName)
}
