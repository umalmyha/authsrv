package dbg

import (
	"net/http"

	"github.com/umalmyha/authsrv/pkg/web"
)

type handler struct{}

func Handler() *handler {
	return &handler{}
}

func (h *handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	// TODO: Add some data later on
	web.RespondStatus(w, http.StatusOK)
}
