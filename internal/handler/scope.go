package handler

import (
	"github.com/umalmyha/authsrv/internal/service"
)

type ScopeHandler struct {
	scopeSrv *service.ScopeService
}

func NewScopeHandler(scopeSrv *service.ScopeService) *ScopeHandler {
	return &ScopeHandler{
		scopeSrv: scopeSrv,
	}
}
