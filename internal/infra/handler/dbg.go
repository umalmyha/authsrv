package handler

import (
	"net/http"

	"github.com/umalmyha/authsrv/pkg/web/response"
)

type DebugHandler struct{}

func NewDebugHandler() *DebugHandler {
	return &DebugHandler{}
}

func (h *DebugHandler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	// TODO: Add some data later on
	response.RespondStatus(w, http.StatusOK)
}
