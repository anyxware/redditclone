package httperr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type BadRequest struct {
	Message string `json:"message"`
}

func (e BadRequest) Error() string {
	return fmt.Sprintf("bad request: %s", e.Message)
}

type NotFound struct {
	Message string `json:"message"`
}

func (e NotFound) Error() string {
	return fmt.Sprintf("source not found: %s", e.Message)
}

type Unauthorized struct {
	Message string `json:"message"`
}

func (e Unauthorized) Error() string {
	return fmt.Sprintf("unauthorized: %s", e.Message)
}

type Forbidden struct {
	Message string `json:"message"`
}

func (e Forbidden) Error() string {
	return fmt.Sprintf("forbidden: %s", e.Message)
}

type UnprocessableEntityItem struct {
	Location string `json:"location"`
	Param    string `json:"param"`
	Value    string `json:"value"`
	Message  string `json:"msg"`
}

type UnprocessableEntity struct {
	Errors []UnprocessableEntityItem `json:"errors"`
}

type InternalError struct {
	Message string `json:"message"`
}

func (e InternalError) Error() string {
	return fmt.Sprintf("internal error: %s", e.Message)
}

func (e UnprocessableEntity) Error() string {
	errs := make([]string, 0)

	for _, item := range e.Errors {
		errs = append(errs, fmt.Sprintf("location: %s, param: %s, value: %s, message: %s", item.Location, item.Param, item.Value, item.Message))
	}

	return strings.Join([]string{"unprocessable entity", strings.Join(errs, ", ")}, ": ")
}

func HandleError(w http.ResponseWriter, inputErr error) {
	var (
		resp       []byte
		err        error
		statusCode int
		unknown    bool = false
	)

	switch inputErr.(type) {
	case BadRequest:
		resp, err = json.Marshal(inputErr.(BadRequest))
		statusCode = http.StatusBadRequest
	case NotFound:
		resp, err = json.Marshal(inputErr.(NotFound))
		statusCode = http.StatusNotFound
	case Unauthorized:
		resp, err = json.Marshal(inputErr.(Unauthorized))
		statusCode = http.StatusUnauthorized
	case Forbidden:
		resp, err = json.Marshal(inputErr.(Forbidden))
		statusCode = http.StatusForbidden
	case UnprocessableEntity:
		resp, err = json.Marshal(inputErr.(UnprocessableEntity))
		statusCode = http.StatusUnprocessableEntity
	case InternalError:
		resp, err = json.Marshal(inputErr.(InternalError))
		statusCode = http.StatusInternalServerError
	default:
		statusCode = http.StatusInternalServerError
		unknown = true
	}

	if unknown || err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	http.Error(w, string(resp), statusCode)
}
