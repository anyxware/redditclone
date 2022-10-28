package handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"redditclone/internal/model/customerr"
	"redditclone/pkg/httperr"
	"redditclone/pkg/httpvalidator"
)

func (h *Handler) handleValidationErrors(w http.ResponseWriter, errs []httpvalidator.ValidationError) {
	res := httperr.UnprocessableEntity{Errors: make([]httperr.UnprocessableEntityItem, 0)}

	for _, err := range errs {
		res.Errors = append(res.Errors, httperr.UnprocessableEntityItem{
			Location: err.Location,
			Param:    err.Param,
			Value:    err.Value,
			Message:  err.Message,
		})
		logrus.Errorln(err.Error())
	}

	httperr.HandleError(w, res)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case customerr.UserAlreadyExists:
		httperr.HandleError(w, httperr.UnprocessableEntity{
			Errors: []httperr.UnprocessableEntityItem{{
				Location: "body",
				Param:    "username",
				Value:    err.(customerr.UserAlreadyExists).Username,
				Message:  "already exists",
			}},
		})
	case customerr.WrongCredential:
		httperr.HandleError(w, httperr.Unauthorized{Message: "wrong credential"})
	case customerr.Unauthorized:
		httperr.HandleError(w, httperr.Unauthorized{Message: "user unauthorized"})
	case customerr.UserNotFoundByID:
		httperr.HandleError(w, httperr.NotFound{Message: "user not found"})
	case customerr.UserNotFoundByUsername:
		httperr.HandleError(w, httperr.NotFound{Message: "user not found"})
	case customerr.PostNotFoundByID:
		httperr.HandleError(w, httperr.NotFound{Message: "post not found"})
	case customerr.CommentNotFoundByID:
		httperr.HandleError(w, httperr.NotFound{Message: "comment not found"})
	case customerr.NotOwner:
		httperr.HandleError(w, httperr.Forbidden{Message: "user not own this resource"})
	case customerr.RequestNotParsed:
		httperr.HandleError(w, httperr.BadRequest{Message: "bad request"})
	default:
		httperr.HandleError(w, httperr.InternalError{Message: "internal error"})
	}

	logrus.Errorln(err)
}
