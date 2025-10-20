package Utility

import (
	"fmt"
)

type ErrorHandler struct {
	errors []string
}

func (handler *ErrorHandler) AddError(fmtMessage string, args ...interface{}) {
	err := fmt.Sprintf(fmtMessage, args...)
	handler.errors = append(handler.errors, err)
}

func (handler *ErrorHandler) AddErrors(errors []string, lastErrors ...string) {
	for _, err := range errors {
		handler.AddError(err)
	}

	for _, err := range lastErrors {
		handler.AddError(err)
	}
}

func (handler *ErrorHandler) GetErrors() []string {
	return handler.errors
}
