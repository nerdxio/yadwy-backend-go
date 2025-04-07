package common

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type ErrorCode string

type Error struct {
	msg  string
	code ErrorCode
}

func WrapErrorf(code ErrorCode, format string, a ...interface{}) error {
	return &Error{
		code: code,
		msg:  fmt.Sprintf(format, a...),
	}
}

func NewErrorf(code ErrorCode, format string, a ...interface{}) error {
	return WrapErrorf(code, format, a...)
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Code() ErrorCode {
	return e.code
}

func SendError(w http.ResponseWriter, status int, errCode, message string) {
	err := Encode(w, status, ErrorResponse{
		Error:   errCode,
		Message: message,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// FormatValidationError formats validator errors into a readable message
func FormatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var errMsg string
		for _, e := range validationErrors {
			if errMsg != "" {
				errMsg += "; "
			}
			errMsg += fmt.Sprintf("field %s is %s", e.Field(), e.Tag())
		}
		return errMsg
	}
	return err.Error()
}
