package error_handler

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

var NoType = ErrorType{Type: 0, ToUserType: 0}

type ErrorType struct {
	Type       uint
	ToUserType uint
}

var toUserErrorMessage = map[int]string{
	500: "Внутренняя ошибка сервера",
	704: "Элемент не найден",
}

type SuperservisError struct {
	errorType     ErrorType
	originalError error
	contextInfo   []errorContext
}

type superservisJsonBError struct {
	Code          int               `json:"code"`
	OriginalError string            `json:"originalError"`
	Context       map[string]string `json:"context"`
}

type errorContext struct {
	Field   string
	Message string
}

func (error SuperservisError) Error() string {
	return error.originalError.Error()
}

func (Etype ErrorType) New(msg string) error {
	return SuperservisError{errorType: Etype, originalError: errors.New(msg)}
}

func (Etype ErrorType) Newf(msg string, args ...interface{}) error {
	err := fmt.Errorf(msg, args...)
	return SuperservisError{errorType: Etype, originalError: err}
}

func (Etype ErrorType) Wrap(err error, msg string) error {

	return Etype.Wrapf(err, msg)
}

func (Etype ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	newErr := errors.Wrapf(err, msg, args...)
	return SuperservisError{errorType: Etype, originalError: newErr}
}

// New creates a no type error

func NewSuperservisError() *SuperservisError {
	return &SuperservisError{errorType: NoType, originalError: errors.New(`Standart error`)}
}

func NewSuperservisGormError(GormErorrs []error) error {
	NewSErr := SuperservisError{errorType: ErrorType{0, 500}, originalError: errors.New(`Ошибка(и) при выполнении операции в БД`)}
	NewNotFoundErr := SuperservisError{errorType: ErrorType{0, 704}, originalError: errors.New(`"Запись не найдена"`)}
	if len(GormErorrs) > 0 {
		for i, Error := range GormErorrs {
			if Error.Error() == "record not found" {
				NewSErr = AddErrorContext(NewNotFoundErr, strconv.Itoa(i), Error.Error()).(SuperservisError)
				return NewSErr
			}
			NewSErr = AddErrorContext(NewSErr, strconv.Itoa(i), Error.Error()).(SuperservisError)
		}
		return NewSErr
	}
	return nil
}

func New(ErrorCode, ToUserErrorCode uint, msg string) error {
	return SuperservisError{errorType: ErrorType{Type: ErrorCode, ToUserType: ToUserErrorCode}, originalError: errors.New(msg)}
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	return SuperservisError{errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

// Wrap wrans an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// Wrapf wraps an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(SuperservisError); ok {
		return SuperservisError{
			errorType:     customErr.errorType,
			originalError: wrappedError,
			contextInfo:   customErr.contextInfo,
		}
	}

	return SuperservisError{errorType: NoType, originalError: wrappedError}
}
