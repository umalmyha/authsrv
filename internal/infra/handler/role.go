package handler

import (
	"net/http"

	"github.com/umalmyha/authsrv/internal/business/role"
	"github.com/umalmyha/authsrv/internal/infra/service"
	"github.com/umalmyha/authsrv/pkg/web/request"
)

type RoleHandler struct {
	roleSrv *service.RoleService
}

func NewRoleHandler(roleSrv *service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleSrv: roleSrv,
	}
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) error {
	var nr role.NewRoleDto
	if err := request.JsonReqBody(r, &nr); err != nil {
		return err
	}
	return h.roleSrv.CreateRole(r.Context(), nr)
}

func (h *RoleHandler) AssignScope(w http.ResponseWriter, r *http.Request) error {
	assignment := struct {
		RoleName  string `json:"role"`
		ScopeName string `json:"scope"`
	}{}

	if err := request.JsonReqBody(r, &assignment); err != nil {
		return err
	}
	return h.roleSrv.AssignScope(r.Context(), assignment.RoleName, assignment.ScopeName)
}

func (h *RoleHandler) UnassignScope(w http.ResponseWriter, r *http.Request) error {
	assignment := struct {
		RoleName  string `json:"role"`
		ScopeName string `json:"scope"`
	}{}

	if err := request.JsonReqBody(r, &assignment); err != nil {
		return err
	}
	return h.roleSrv.UnassignScope(r.Context(), assignment.RoleName, assignment.ScopeName)
}
