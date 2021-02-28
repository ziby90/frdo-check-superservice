package entities

import (
	"frdo-check-superservice/model"
	"frdo-check-superservice/pkg/service"
)

type GetAllDocs struct {
	services *service.Service
}

func NewGetAllDocs(services *service.Service) *GetAllDocs {
	return &GetAllDocs{services: services}
}

func (g *GetAllDocs) Search(request model.FRDORequest) (string, uint, error) {
	// search in db persons and documents
	resultsSearch, err := g.SearchInDB(request.GetDocumentsRequest.Person)
	if err != nil {
		return ``, 0, err
	}

	//logrus.Println(resultsSearch)
	code, _ := g.GetResponse(resultsSearch)

	// create xml response
	xmlBody, err := g.SetResponse(code, resultsSearch)
	if err != nil {
		return ``, 0, err
	}
	return xmlBody, code, nil
}

func (g *GetAllDocs) SearchInDB(person model.Person) ([]model.EduDocument, error) {
	resultsSearch, err := g.services.GetDocuments(person)
	return resultsSearch, err
}

func (g *GetAllDocs) GetResponse(resultsSearch []model.EduDocument) (uint, string) {
	code, message := service.GetResponseAllDocs(resultsSearch)
	return code, message
}

func (g *GetAllDocs) SetResponse(code uint, docs []model.EduDocument) (string, error) {
	xmlBody, err := g.services.SetResponseAllDocs(code, docs)
	return xmlBody, err
}
