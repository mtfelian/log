package error

import (
	"fmt"
	. "github.com/mtfelian/utils"
)

// CodeSuccess is a success error code
const (
	CodeSuccess uint = iota
)

// StandardError is a standard error to return with Gin
type StandardError struct {
	FCode    uint   `json:"code"`
	FMessage *string `json:"error,omitempty"`
}

// Error is an interface for error
type Error interface {
	error
	Code() uint
	Message() string
}

// Error implements builtin error interface
func (err StandardError) Error() string {
	return fmt.Sprintf("%d: %s", err.Code(), err.Message())
}

// NewError returns new standard error with code and message from builtin error
func NewError(code uint, err error) Error {
	return StandardError{code, PString(err.Error())}
}

// NewErrorf return new standard error with code, message msg and optional printf args
func NewErrorf(code uint, msg string, args ...interface{}) Error {
	return StandardError{code, PString(fmt.Sprintf(msg, args...))}
}

// MayError makes StandardError from builtin error
func MayError(code uint, err error) Error {
	if err == nil {
		return nil
	}
	return NewError(code, err)
}

// Code returns an error code
func (err StandardError) Code() uint {
	return err.FCode
}

// Message returns an error message
func (err StandardError) Message() string {
	if err.FMessage == nil {
		return "<no message>"
	}
	return *err.FMessage
}

// String реализует интерфейс stringer
func (err StandardError) String() string {
	return fmt.Sprintf("%d: %s", err.Code(), err.Message())
}