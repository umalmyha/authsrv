package service

import (
	"github.com/umalmyha/authsrv/internal/business/scope"
)

type ScopeService struct {
	scopeRepo *scope.Repository
}

func NewScopeService(scopeRepo *scope.Repository) *ScopeService {
	return &ScopeService{
		scopeRepo: scopeRepo,
	}
}
