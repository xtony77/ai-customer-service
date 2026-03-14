package domain

import (
	"errors"
	"net/http"
)

type ErrorFormat struct {
	HttpStatus int
	Code       int
	Message    string
}

var (
	ErrorServer            = ErrorFormat{HttpStatus: http.StatusInternalServerError, Code: http.StatusInternalServerError, Message: "Server Error"}
	ErrorBadRequest        = ErrorFormat{HttpStatus: http.StatusBadRequest, Code: http.StatusBadRequest, Message: "bad request"}
	ErrPromptGuardRejected = errors.New("question rejected by prompt guard")
)
