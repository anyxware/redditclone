package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"redditclone/internal/model/customerr"
	"redditclone/pkg/token"
	"strings"
)

type loggedWriter struct {
	statusCode     int
	data           []byte
	responseWriter http.ResponseWriter
}

func (w *loggedWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.responseWriter.WriteHeader(statusCode)
}

func (w *loggedWriter) Write(data []byte) (int, error) {
	w.data = data
	return w.responseWriter.Write(data)
}

func (w *loggedWriter) Header() http.Header {
	return w.responseWriter.Header()
}

func (h *Handler) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggedWriter := &loggedWriter{responseWriter: w}
		next.ServeHTTP(loggedWriter, r)
		logrus.Infoln(r.Method, r.URL.Path, loggedWriter.statusCode)
		logrus.Debugln(string(loggedWriter.data))
	})
}

func (h *Handler) getToken(r *http.Request) (string, error) {
	header := r.Header.Get(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		return "", errors.New("invalid auth header")
	}

	return headerParts[1], nil
}

func (h *Handler) getAuthUserFromJWT(t string) (token.AuthUser, error) {
	return h.signer.ParseToken(t)
}

func (h *Handler) getAuthUserFromCookie(t string) (token.AuthUser, error) {
	serializedUser, err := h.sessions.GetCookie(t)
	if err != nil {
		return token.AuthUser{}, err
	}
	var authUser token.AuthUser
	err = json.Unmarshal(serializedUser, &authUser)
	return authUser, err
}

func (h *Handler) authorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := h.getToken(r)
		if err != nil {
			h.handleError(w, customerr.Unauthorized{Message: err.Error()})
			return
		}

		authUser, err := h.getAuthUserFromCookie(t)
		if err != nil {
			h.handleError(w, customerr.Unauthorized{Message: err.Error()})
		}

		usr, err := h.service.GetUserByID(authUser.ID)
		if err != nil {
			if _, ok := err.(customerr.UserNotFoundByID); ok {
				h.handleError(w, customerr.Unauthorized{Message: err.Error()})
			}
			h.handleError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), "user", usr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) recoverPanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				h.handleError(w, errors.New("panic detected"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
