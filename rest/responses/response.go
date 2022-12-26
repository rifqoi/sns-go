package responses

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	snsgo "github.com/rifqoi/sns-go"
)

type Response struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, msg string, err error) {
	var status int

	var ierr *snsgo.Error
	if errors.As(err, &ierr) {
		status = http.StatusInternalServerError

	}

	switch ierr.Code() {
	case snsgo.ErrorCodeInvalidArg:
		status = http.StatusBadRequest
	case snsgo.ErrorCodeUnauthorized:
		status = http.StatusUnauthorized
	case snsgo.ErrorCodeNotFound:
		status = http.StatusNotFound
	case snsgo.ErrorCodeUnknown:
		fallthrough
	default:
		status = http.StatusInternalServerError
	}

	resp := Response{
		Error:   err.Error(),
		Message: msg,
	}

	render.Status(r, status)
	render.JSON(w, r, resp)
}

func SuccessResponse(w http.ResponseWriter, r *http.Request, status int, msg string, data any) {
	resp := Response{
		Message: msg,
		Data:    data,
	}

	render.Status(r, status)
	render.JSON(w, r, resp)
}
