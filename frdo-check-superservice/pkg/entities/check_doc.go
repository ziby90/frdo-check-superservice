package entities

import (
	"frdo-check-superservice/model"
	"frdo-check-superservice/pkg/service"
)

type CheckDoc struct {
	services *service.Service
}

func NewCheckDoc(services *service.Service) *CheckDoc {
	return &CheckDoc{services: services}
}

func (c *CheckDoc) Search(request model.FRDORequest) (string, uint, error) {
	// search in db persons and documents
	resultsSearch, err := c.SearchInDB(request.CheckDocumentRequest.Person, request.CheckDocumentRequest.EduDocument)
	if err != nil {
		return ``, 0, err
	}
	//logrus.Println(resultsSearch)
	code, message := c.GetResponse(resultsSearch)

	// create xml response
	xmlBody, err := c.SetResponse(code, message)
	if err != nil {
		return ``, 0, err
	}
	return xmlBody, code, nil
}

func (c *CheckDoc) SearchInDB(person model.Person, document model.EduDocument) ([]model.ResultDb, error) {
	resultsSearch, err := c.services.Search(person, document)
	return resultsSearch, err
}

func (c *CheckDoc) GetResponse(resultsSearch []model.ResultDb) (uint, string) {
	code, message := service.GetResponse(resultsSearch)
	return code, message
}

func (c *CheckDoc) SetResponse(code uint, message string) (string, error) {
	xmlBody, err := c.services.SetResponseSearch(code, message)
	return xmlBody, err
}
