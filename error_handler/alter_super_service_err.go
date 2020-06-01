package error_handler

import (
	"fmt"
	"github.com/pkg/errors"
)

type AlterSuperservisErrorContainer struct {
	Errors []AlterSuperservisError
}

func (ASEC AlterSuperservisErrorContainer) Error() string {
	return "это контейнер для ошибок"
}

//
func NewAlterSuperservisErrorContainer() AlterSuperservisErrorContainer {
	return AlterSuperservisErrorContainer{}
}

func (ASEC *AlterSuperservisErrorContainer) AddErrror(err error) error {

	if customErr, ok := err.(AlterSuperservisError); ok {
		ASEC.Errors = append(ASEC.Errors, customErr)
		return ASEC
	}
	fmt.Println("ADD")
	return ASEC

}
func (ASEC *AlterSuperservisErrorContainer) HaveErrors() bool {
	if len(ASEC.Errors) > 0 {
		return true
	}
	return false
}

type AlterSuperservisError struct {
	ErrorType     AlterErrorType
	OriginalError error
	ContextInfo   []AlterSuperservisError
}

type AlterErrorType struct {
	Code            uint
	UIDEpgu         int64
	UID             string
	ToUserCode      uint
	Object          string
	CommentTemplate string
	Choicer         EntityIDChoicer
}

func (error AlterSuperservisError) Error() string {
	return error.OriginalError.Error()
}

func (AEtype AlterErrorType) New(msg ...interface{}) error {

	return AlterSuperservisError{ErrorType: AEtype, OriginalError: errors.New(fmt.Sprintf(AEtype.CommentTemplate, msg...))}
}

func (AEtype AlterErrorType) NewWithUID(Uid string, msg ...interface{}) error {
	AEtype.UID = Uid
	return AlterSuperservisError{ErrorType: AEtype, OriginalError: errors.New(fmt.Sprintf(AEtype.CommentTemplate, msg...))}
}

func (AEtype AlterErrorType) Make(msg ...interface{}) AlterSuperservisError {

	return AlterSuperservisError{ErrorType: AEtype, OriginalError: errors.New(fmt.Sprintf(AEtype.CommentTemplate, msg...))}
}

func (AEtype AlterErrorType) MakeWithUID(Uid string, msg ...interface{}) AlterSuperservisError {
	AEtype.UID = Uid
	return AlterSuperservisError{ErrorType: AEtype, OriginalError: errors.New(fmt.Sprintf(AEtype.CommentTemplate, msg...))}
}

type EntityIDChoicer interface {
	GetIDAndType() (UIDEpgu int64, UID, Type, ID string)
}

//Создание ошибки для сущьностей с uid и uid_epgu
func (AEtype AlterErrorType) MakeWithUIDChoice(Choicer EntityIDChoicer, msg ...interface{}) AlterSuperservisError {

	var Type interface{}
	var ID interface{}
	AEtype.UIDEpgu, AEtype.UID, Type, ID = Choicer.GetIDAndType()
	msgs := make([]interface{}, cap(msg)+2)
	msgs[0] = Type
	msgs[1] = ID
	for i, Item := range msg {
		msgs[i+2] = Item
	}

	return AlterSuperservisError{ErrorType: AEtype, OriginalError: errors.New(fmt.Sprintf(AEtype.CommentTemplate, msgs...))}
}

// NewType когда необходимо поменять тип обьекта в ошибке
func (AEtype AlterErrorType) NewWithObject(Object string, msg ...interface{}) error {
	AEtype.Object = Object
	return AlterSuperservisError{ErrorType: AEtype, OriginalError: errors.New(fmt.Sprintf(AEtype.CommentTemplate, msg...))}
}

var AlterNoType = AlterErrorType{Code: 0, ToUserCode: 500, Object: "", CommentTemplate: "Неизвестная ошибка"}

//Для установки UID или UIDEpgu при генерации ошибки в коде ниже
func AddUIDOrUIDEpgu(err error, Choicer EntityIDChoicer) error {
	if Err, ok := err.(AlterSuperservisError); ok {
		Err.ErrorType.UIDEpgu, Err.ErrorType.UID, _, _ = Choicer.GetIDAndType()
		return Err
	}
	return err
}

// Для установки Идентификатора Вуза
func SetUID(err error, UID string) error {
	if Err, ok := err.(AlterSuperservisError); ok {
		Err.ErrorType.UID = UID
		return Err
	}
	return err
}

// Для установки Идентификатора ЕПГУ
func SetUIDEpgu(err error, UIDEpgu int64) error {
	if Err, ok := err.(AlterSuperservisError); ok {
		Err.ErrorType.UIDEpgu = UIDEpgu
		return Err
	}
	return err
}
