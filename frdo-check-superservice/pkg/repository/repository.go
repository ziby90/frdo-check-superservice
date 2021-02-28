package repository

import (
	"10.10.11.220/ursgis/rabbitMQ.git"
	"frdo-check-superservice/model"
	"frdo-check-superservice/pkg/repository/rabbit"
	"frdo-check-superservice/pkg/repository/sql"
	"frdo-check-superservice/pkg/repository/xmlpars"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type SearchInDb interface {
	SearchPerson(person model.Person, nameModule string) ([]uint, error)
	SearchDocument(recipientId uint, document model.EduDocument) ([]uint, error)
	SearchAllDocument(recipientId uint, nameModule string) ([]model.EduDocument, error)
	AddNewQueue(smevMsg model.SMEVMessage, typeQueue string) (uint, error)
	GetDbConnections() map[string]*sqlx.DB
	GetEducationLevel() map[int]string
	UpdateQueueLog(smevMsg model.SMEVMessage, typeQueue string) error
}

type QueueWorker interface {
	GetNewMessage() (rabbitMQ.Message, error)
	ParseMessage(body []byte) (model.SMEVMessage, error)
	SetMessage(request model.SMEVMessage, xmlBody string) ([]byte, error)
	PublishMessage(msg []byte) error
	SetStatusMessage(requeue bool) error
}

type ParsingXML interface {
	ParseXmlToStruct(body string) (model.FRDORequest, error)
	SetResponseSearchToXml(response model.CheckDocumentResponse) (string, error)
	SetResponseAllDocsToXml(response model.GetDocumentsResponse) (string, error)
}

type Repository struct {
	SearchInDb
	QueueWorker
	ParsingXML
}

func NewRepository(dbs map[string]*sqlx.DB, ch *rabbitMQ.RabbitChan) *Repository {
	return &Repository{
		SearchInDb:  sql.NewSearch(dbs),
		QueueWorker: rabbit.NewQueueWorker(ch),
		ParsingXML:  xmlpars.NewParsingXML(),
	}
}
