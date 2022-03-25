package handler

import (
	"net/http"

	"github.com/umalmyha/authsrv/internal/business/scope"
	"github.com/umalmyha/authsrv/internal/service"
	"github.com/umalmyha/authsrv/pkg/web"
)

type ScopeHandler struct {
	scopeSrv *service.ScopeService
}

func NewScopeHandler(scopeSrv *service.ScopeService) *ScopeHandler {
	return &ScopeHandler{
		scopeSrv: scopeSrv,
	}
}

func (h *ScopeHandler) CreateScope(w http.ResponseWriter, r *http.Request) error {
	var ns scope.NewScopeDto
	if err := web.JsonReqBody(r, &ns); err != nil {
		return err
	}
	return h.scopeSrv.CreateScope(r.Context(), ns)
}
