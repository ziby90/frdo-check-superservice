package digest

import (
	"persons/error_handler"
)

type Result struct {
	Done   bool
	Errors []Error
}

type XmlPath struct {
	PathXml          string
	PathXsd          string
	TrueResultAction string
}

var validBase64Err = error_handler.ErrorType{Type: 1, ToUserType: 500}
var validJsonErr = error_handler.ErrorType{Type: 2, ToUserType: 500}
var validDataErr = error_handler.ErrorType{Type: 3, ToUserType: 500}
var validXmlErr = error_handler.ErrorType{Type: 4, ToUserType: 500}

type Error struct {
	Code         string
	Message      string
	Index        interface{}
	JsonValidate []string
}
