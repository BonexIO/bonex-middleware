package types

import (
	"fmt"
	"net/http"
)

type ErrorCode string

// common
const (
	ErrService          ErrorCode = "ERR_SERVICE"
	ErrNotFound         ErrorCode = "ERR_NOT_FOUND"
	ErrAlreadyExists    ErrorCode = "ERR_ALREADY_EXISTS"
	ErrAlreadyEnabled   ErrorCode = "ERR_ALREADY_ENABLED"
	ErrAlreadyDisabled  ErrorCode = "ERR_ALREADY_DISABLED"
	ErrBadRequest       ErrorCode = "ERR_BAD_REQUEST"
	ErrBadParam         ErrorCode = "ERR_BAD_PARAM"
	ErrNotAllowed       ErrorCode = "ERR_NOT_ALLOWED"
	ErrBadJwt           ErrorCode = "ERR_BAD_JWT"
	ErrBalanceNotEnough ErrorCode = "ERR_BALANCE_NOT_ENOUGH"
	ErrBadAuth          ErrorCode = "ERR_BAD_AUTH"
	ErrBadSignature     ErrorCode = "ERR_BAD_SIGNATURE"
	ErrLimitReached     ErrorCode = "ERR_LIMIT_REACHED"
)

type (
	ServiceError interface {
		error
		ErrorCode() ErrorCode
		ToMap(http.ResponseWriter) map[string]interface{}
		GetHttpCode() int
	}

	Error struct {
		error
		Code        ErrorCode
		Value       string
		Description string
	}
)

func (e Error) Error() string {
	return fmt.Sprintf("%s %s", string(e.Code), e.Value)
}

func (e Error) ErrorCode() ErrorCode {
	return e.Code
}

// ToMap converts Error object to map[string]interface{}
func (e Error) ToMap() map[string]interface{} {
	r := map[string]interface{}{
		"error": string(e.Code),
	}

	if string(e.Value) != "" {
		r["value"] = string(e.Value)
	}

	if string(e.Description) != "" {
		r["description"] = string(e.Description)
	}

	return r
}

// GetHttpCode return a Http error code
func (e Error) GetHttpCode() int {
	switch e.Code {
	case ErrService:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

// New creates an Error object
func NewError(code ErrorCode, value ...string) *Error {
	e := &Error{Code: code}
	if len(value) > 0 {
		e.Value = value[0]
	}
	return e
}

// NewWithDesc creates an Error object with description
func NewWithDesc(code ErrorCode, desc string, value ...string) *Error {
	e := &Error{Code: code, Description: desc}
	if len(value) > 0 {
		e.Value = value[0]
	}
	return e
}

// FromError creates a new Error (ErrService) from common golang error
func FromError(err error) *Error {
	if err != nil {
		return &Error{
			Code:  ErrService,
			Value: "",
		}
	}

	return nil
}
