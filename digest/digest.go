package digest

import (
	"persons/error_handler"
	"time"
)

type Result struct {
	Done   bool
	Errors []Error
}

type Entity interface {
	CheckUid(uid string, p PrimaryDataDigest) error
	GetById(id uint) (p *PrimaryDataDigest, err error)
}
type PrimaryDataDigest struct {
	Id             uint
	Uid            *string
	Actual         bool
	IdOrganization uint
	Created        time.Time
	TableName      string
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
