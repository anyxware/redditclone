package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"redditclone/internal/model"
	"redditclone/internal/model/customerr"
	"redditclone/pkg/token"
)

func writeToken(w http.ResponseWriter, t string) error {
	resp := []byte(fmt.Sprintf("{\"token\": \"%s\"}", t))
	if _, err := w.Write(resp); err != nil {
		return err
	}

	return nil
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var input map[string]string
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.handleError(w, customerr.RequestNotParsed{Message: err.Error()})
		return
	}

	if errs := h.validator.ValidateBody("Credential", input); len(errs) != 0 {
		h.handleValidationErrors(w, errs)
		return
	}

	cred := model.Credential{Username: input["username"], Password: input["password"]}

	usr, err := h.service.RegisterUser(cred)
	if err != nil {
		h.handleError(w, err)
		return
	}

	authUser := token.AuthUser{ID: usr.ID, Username: usr.Username}
	t, err := h.signer.CreateToken(authUser)
	if err != nil {
		h.handleError(w, err)
		return
	}
	serialized, err := json.Marshal(authUser)
	if err != nil {
		h.handleError(w, err)
		return
	}
	if err = h.sessions.AddCookie(t, serialized); err != nil {
		h.handleError(w, err)
		return
	}

	if err = writeToken(w, t); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var input map[string]string
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.handleError(w, customerr.RequestNotParsed{Message: err.Error()})
		return
	}

	if errs := h.validator.ValidateBody("Credential", input); len(errs) != 0 {
		h.handleValidationErrors(w, errs)
		return
	}

	cred := model.Credential{Username: input["username"], Password: input["password"]}

	usr, err := h.service.LoginUser(cred)
	if err != nil {
		h.handleError(w, err)
		return
	}

	authUser := token.AuthUser{ID: usr.ID, Username: usr.Username}
	t, err := h.signer.CreateToken(authUser)
	if err != nil {
		h.handleError(w, err)
		return
	}
	serialized, err := json.Marshal(authUser)
	if err != nil {
		h.handleError(w, err)
		return
	}
	if err = h.sessions.AddCookie(t, serialized); err != nil {
		h.handleError(w, err)
		return
	}

	if err = writeToken(w, t); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
}
