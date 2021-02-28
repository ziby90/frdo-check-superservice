package service

import (
	"frdo-check-superservice/model"
	"frdo-check-superservice/pkg/repository"
)

// go:generate mockgen -source=service.go -destination=mocks/mock.go
type SearchInDb interface {
	Search(person model.Person, document model.EduDocument) ([]model.ResultDb, error)
	GetDocuments(person model.Person) ([]model.EduDocument, error)
	AddNewQueueLog(smevMsg model.SMEVMessage, typeQueue string) (uint, error)
	UpdateQueueLog(smevMsg model.SMEVMessage, typeQueue string) error
}

type QueueWorker interface {
	GetNewMessage() (model.SMEVMessage, error)
	PostRequest(request model.SMEVMessage, xmlBody string) error
}

type ParsingXML interface {
	ParseXmlToStruct(body string) (model.FRDORequest, error)
	SetResponseAllDocs(code uint, docs []model.EduDocument) (string, error)
	SetResponseSearch(code uint, message string) (string, error)
}

type Service struct {
	SearchInDb
	QueueWorker
	ParsingXML
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		SearchInDb:  NewSearchService(repos.SearchInDb),
		QueueWorker: NewQueueService(repos.QueueWorker),
		ParsingXML:  NewXmlParsService(repos.ParsingXML),
	}
}
