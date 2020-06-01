package error_handler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm/dialects/postgres"
)

//------------------------------------------------------------------------------------

type ErrorGetter interface {
	AddError(UID, ErrorDecryption string, ErrorCode uint, UIDEpgu int64)
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

// AddAlterErrorContext adds a context to an AlterSuperservisError
func AddAlterErrorContext(err error, AddErr AlterSuperservisError) error {
	if customErr, ok := err.(AlterSuperservisError); ok {
		context := append(customErr.ContextInfo, AddErr)
		return AlterSuperservisError{ErrorType: customErr.ErrorType, OriginalError: customErr.OriginalError, ContextInfo: context}
	} else {
		return AlterSuperservisError{ErrorType: AlterErrorType{}, OriginalError: err, ContextInfo: customErr.ContextInfo}
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
	if customErr, ok := err.(AlterSuperservisError); ok {
		mapContext := make(map[string]string)
		for i, contextInfostr := range customErr.ContextInfo {
			mapContext[strconv.Itoa(i)] = "field :" + contextInfostr.ErrorType.Object + " message : " + contextInfostr.Error()
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
	if customErr, ok := err.(AlterSuperservisError); ok {
		mapContext := make(map[string]string)
		mapContext["0"] = "field: OriginalErrorText" + " message :" + err.Error()
		for i, contextInfostr := range customErr.ContextInfo {
			mapContext[strconv.Itoa(i+1)] = "field :" + contextInfostr.ErrorType.Object + " message : " + contextInfostr.Error()
		}
		return mapContext
	}

	return nil
}

func GetErrsToToken(err error, ErrPayload ErrorGetter) {
	ToUserErrType := GetErrorCodeToUser(err)
	if ToUserErrType >= 6000 {
		ErrPayload.AddError("", "Ошибка сервера", 6000, 0)
	} else {
		if customErr, ok := err.(SuperservisError); ok {
			ErrPayload.AddError("", customErr.Error(), customErr.errorType.ToUserType, 0)
			return
		}
		if customErr, ok := err.(AlterSuperservisError); ok {
			ErrPayload.AddError(customErr.ErrorType.UID, customErr.Error(), customErr.ErrorType.ToUserCode, customErr.ErrorType.UIDEpgu)
			return
		}
		if customErr, ok := err.(AlterSuperservisErrorContainer); ok {
			fmt.Println(len(customErr.Errors))
			for _, ErrItem := range customErr.Errors {
				ErrPayload.AddError(ErrItem.ErrorType.UID, ErrItem.Error(), ErrItem.ErrorType.ToUserCode, ErrItem.ErrorType.UIDEpgu)
			}
		}

	}

	//Сделать логирование в ошибки файл

}

func GetContextToUser(err error) (ToUserCode uint, Context string) {
	ToUserErrType := GetErrorCodeToUser(err)
	if ToUserErrType >= 6000 {
		if val, ok := toUserErrorMessage[int(ToUserErrType)]; ok {
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
	if AlterSuperservisError, ok := err.(*AlterSuperservisError); ok {
		return AlterSuperservisError.ErrorType.Code
	}
	return NoType.Type
}

// GetType returns the error type ToUserType
func GetErrorCode(err error) uint {
	if SuperservisError, ok := err.(SuperservisError); ok {
		return SuperservisError.errorType.Type
	}
	if AlterSuperservisError, ok := err.(AlterSuperservisError); ok {
		return AlterSuperservisError.ErrorType.ToUserCode
	}
	return NoType.Type
}

// GetType returns the error type ToUserType
func GetErrorCodeToUser(err error) uint {
	if SuperservisError, ok := err.(SuperservisError); ok {
		return SuperservisError.errorType.ToUserType
	}
	if AlterSuperservisError, ok := err.(AlterSuperservisError); ok {
		return AlterSuperservisError.ErrorType.ToUserCode
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
