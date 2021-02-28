package service

import (
	"frdo-check-superservice/model"
	"frdo-check-superservice/pkg/repository"
	"frdo-check-superservice/pkg/repository/xmlpars"
)

type XmlParsService struct {
	repo repository.ParsingXML
}

func NewXmlParsService(repo repository.ParsingXML) *XmlParsService {
	return &XmlParsService{repo: repo}
}

func (x *XmlParsService) ParseXmlToStruct(body string) (model.FRDORequest, error) {
	//logrus.Println(`body:`, body)
	str, err := xmlpars.GetStringFromBase64(body)
	//logrus.Printf(`result : %v`, str)
	if err != nil {
		return model.FRDORequest{}, err
	}
	result, err := x.repo.ParseXmlToStruct(str)
	return result, err
}

func (x *XmlParsService) SetResponseSearch(code uint, message string) (string, error) {
	var result string
	response := model.CheckDocumentResponse{
		CheckDocumentCode: code,
		CheckDocumentText: message,
	}
	body, err := x.repo.SetResponseSearchToXml(response)
	if err != nil {
		return result, err
	}
	result = xmlpars.SetStringToBase64(body)
	return result, nil
}

func (x *XmlParsService) SetResponseAllDocs(code uint, docs []model.EduDocument) (string, error) {
	var result string
	response := model.GetDocumentsResponse{
		Result: code,
		Documents: struct {
			Document []model.EduDocument `xml:"ns1:Document,omitempty"`
		}{
			docs,
		},
	}
	body, err := x.repo.SetResponseAllDocsToXml(response)
	if err != nil {
		return result, err
	}
	result = xmlpars.SetStringToBase64(body)
	return result, nil
}
