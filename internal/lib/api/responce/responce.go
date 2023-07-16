package responce

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Responce struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Responce {
	return Responce{
		Status: StatusOK,
	}
}

func Error(msg string) Responce {
	return Responce{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationErrors(errs validator.ValidationErrors) Responce {
	var errMsg []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("field %s is required", err.Field()))
		case "url":
			errMsg = append(errMsg, fmt.Sprintf("field %s is not valid URL", err.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}
	return Responce{
		Status: StatusError,
		Error:  strings.Join(errMsg, ", "),
	}
}
