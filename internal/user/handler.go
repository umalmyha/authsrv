package user

import (
	"net/http"

	"github.com/umalmyha/authsrv/internal/user/dto"
	"github.com/umalmyha/authsrv/pkg/web"
)

type handler struct {
	userSrv *service
}

func Handler(userSrv *service) *handler {
	return &handler{
		userSrv: userSrv,
	}
}

func (h *handler) Signup(w http.ResponseWriter, r *http.Request) error {
	var nu dto.NewUser
	if err := web.JsonReqBody(r, &nu); err != nil {
		return err
	}

	user, err := h.userSrv.CreateUser(r.Context(), nu)
	if err != nil {
		return err
	}

	return web.RespondJson(w, http.StatusCreated, user)
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) error {
	users, err := h.userSrv.AllUsers(r.Context())
	if err != nil {
		return err
	}

	return web.RespondJson(w, http.StatusOK, users)
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) error {
	user, err := h.userSrv.GetUser(r.Context(), web.UrlParam(r, "id"))
	if err != nil {
		return err
	}
	return web.RespondJson(w, http.StatusOK, user)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) error {
	var uu dto.UpdateUser
	if err := web.JsonReqBody(r, &uu); err != nil {
		return err
	}

	updatedUser, err := h.userSrv.UpdateUser(r.Context(), web.UrlParam(r, "id"), uu)
	if err != nil {
		return err
	}

	return web.RespondJson(w, http.StatusOK, updatedUser)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) error {
	if err := h.userSrv.DeleteUser(r.Context(), web.UrlParam(r, "id")); err != nil {
		return err
	}

	web.RespondStatus(w, http.StatusNoContent)
	return nil
}
