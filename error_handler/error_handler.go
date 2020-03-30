package error_handler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
)

var toUserErrorMessage = map[uint]string{
	500: "Внутренняя ошибка сервера",
}
var NoType = ErrorType{Type: 0, ToUserType: 0}

type ErrorType struct {
	Type       uint
	ToUserType uint
}

type AlterSuperservisError struct {
	ErrorType     ErrorType
	OriginalError error
	ContextInfo   []errorContext
}

// func NewAlterSuperservisError() *AlterSuperservisError {
// 	return &AlterSuperservisError{ErrorType: NoType, OriginalError: errors.New(`AlterStandart error`)}
// }

// func (*AlterSuperservisError) New(ErrorCode, ToUserErrorCode uint, msg string) *AlterSuperservisError {
// 	return &AlterSuperservisError{ErrorType: ErrorType{Type: ErrorCode, ToUserType: ToUserErrorCode}, OriginalError: errors.New(msg)}
// }

// func (AE *AlterSuperservisError) Error1() string {
// 	return AE.OriginalError.Error()
// }

// func (*AlterSuperservisError) AddErrorContext(err error, field, message string) error {
// 	if customErr, ok := err.(AlterSuperservisError); ok {
// 		context := append(customErr.ContextInfo, errorContext{Field: field, Message: message})
// 		return AlterSuperservisError{ErrorType: customErr.ErrorType, OriginalError: customErr.OriginalError, ContextInfo: context}
// 	} else {
// 		return AlterSuperservisError{ErrorType: NoType, OriginalError: "err", ContextInfo: customErr.ContextInfo}
// 	}
// }

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

func NewTokenErrorContext() *TokenErrorContext {
	return &TokenErrorContext{}
}

type TokenErrorContext struct {
	Header           TokenErrorHeader `xml:"-"`
	ErrorCode        uint
	ErrorDescription string
	EntityType       string
	UID              string
	UIDParent        string
}

type TokenErrorHeader struct {
	IDToken  uint
	DataType string
	Cert64   string
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
	NewSErr := SuperservisError{errorType: NoType, originalError: errors.New(`Ошибка(и) при выполнении операции в БД`)}
	if len(GormErorrs) > 0 {
		for i, Error := range GormErorrs {
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

// AddErrorContext adds a context to an error
func AddErrorContext(err error, field, message string) error {
	if customErr, ok := err.(SuperservisError); ok {
		context := append(customErr.contextInfo, errorContext{Field: field, Message: message})
		return SuperservisError{errorType: customErr.errorType, originalError: customErr.originalError, contextInfo: context}
	} else {
		return SuperservisError{errorType: NoType, originalError: err, contextInfo: customErr.contextInfo}
	}
}

// GetErrorContext returns the error context
func GetErrorContext(err error) map[string]string {
	if customErr, ok := err.(SuperservisError); ok {
		mapContext := make(map[string]string)
		for i, contextInfostr := range customErr.contextInfo {
			mapContext[strconv.Itoa(i)] = "field :" + contextInfostr.Field + " message : " + contextInfostr.Message
		}
		return mapContext
	}

	return nil
}

// GetOriginalErrorMesageAndContext returns original message in [0] map key, then other context, if it exist
func GetOriginalErrorMesageAndContext(err error) map[string]string {
	if customErr, ok := err.(SuperservisError); ok {
		mapContext := make(map[string]string)
		mapContext["0"] = "field: OriginalErrorText" + " message :" + err.Error()
		for i, contextInfostr := range customErr.contextInfo {
			mapContext[strconv.Itoa(i+1)] = "field :" + contextInfostr.Field + " message : " + contextInfostr.Message
		}
		return mapContext
	}
	return nil
}

func GetContextToUser(err error) (ToUserCode uint, Context string) {
	ToUserErrType := GetToUserType(err)
	if ToUserErrType >= 500 {
		if val, ok := toUserErrorMessage[ToUserErrType]; ok {
			Context = val
		} else {
			Context = "Неизвестная ошибка"
		}
		return ToUserErrType, Context
	} else {
		Contex := GetOriginalErrorMesageAndContext(err)
		fmt.Println(Contex)
		ContexJson, _ := json.Marshal(Contex)
		//Сделать логирование в ошибки файл

		return ToUserErrType, string(ContexJson)

	}
}

// GetType returns the error type
func GetType(err error) uint {
	if SuperservisError, ok := err.(*SuperservisError); ok {
		return SuperservisError.errorType.Type
	}
	return NoType.Type
}

// GetType returns the error type ToUserType
func GetToUserType(err error) uint {
	if SuperservisError, ok := err.(SuperservisError); ok {
		return SuperservisError.errorType.ToUserType
	}
	return NoType.ToUserType
}

func ParseErrorToJsonB(err error) postgres.Jsonb {
	sserr := superservisJsonBError{Code: int(GetType(err)), Context: GetErrorContext(err), OriginalError: err.Error()}
	JsonOb, err := json.Marshal(sserr)
	if err != nil {
		fmt.Println(err.Error())
	}
	return postgres.Jsonb{RawMessage: JsonOb}
}

func ParseErrorToXML(err error) string {
	JsonOb, err := xml.Marshal(GetErrorContext(err))
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(JsonOb)
}

// func ShowStructure(s interface{}) {
// 	a := reflect.ValueOf(s)
// 	numfield := reflect.ValueOf(s).Elem().NumField()
// 	if a.Kind() != reflect.Ptr {
// 		log.Fatal("wrong type struct")
// 	}
// 	for x := 0; x < numfield; x++ {
// 		fmt.Printf("Name field: `%s`  Type: `%s`\n", reflect.TypeOf(s).Elem().Field(x).Name,
// 			reflect.ValueOf(s).Elem().Field(x).Type())
// 	}
// }
