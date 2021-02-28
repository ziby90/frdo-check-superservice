package xmlpars

import (
	"encoding/base64"
	"encoding/xml"
	"frdo-check-superservice/model"
	"github.com/sirupsen/logrus"
)

func (x *ParsingXML) ParseXmlToStruct(body string) (model.FRDORequest, error) {
	var xmlEntityRequest model.FRDORequest
	err := xml.Unmarshal([]byte(body), &xmlEntityRequest)
	if err != nil {
		logrus.Printf(`error parsing xml to struct: %s`, err.Error())
	}
	return xmlEntityRequest, err
}

func (x *ParsingXML) SetResponseToXml(response model.FRDOResponse) (string, error) {
	result, err := xml.Marshal(response)
	if err != nil {
		logrus.Printf(`error set struct to xml string: %s`, err.Error())
	}
	return string(result), err
}

func (x *ParsingXML) SetResponseSearchToXml(response model.CheckDocumentResponse) (string, error) {
	// create struct to ready answer xml
	r := struct {
		XMLName               xml.Name `xml:"ns1:FRDOResponse"`
		Ns1                   string   `xml:"xmlns:ns1,attr"`
		CheckDocumentResponse struct {
			CheckDocumentCode uint   `xml:"ns1:CheckDocumentCode,omitempty"`
			CheckDocumentText string `xml:"ns1:CheckDocumentText,omitempty"`
		} `xml:"ns1:CheckDocumentResponse,omitempty"`
	}{}
	// set params
	r.Ns1 = "urn://obrnadzor-gosuslugi-ru/frdo/1.0.1"

	r.CheckDocumentResponse = struct {
		CheckDocumentCode uint   `xml:"ns1:CheckDocumentCode,omitempty"`
		CheckDocumentText string `xml:"ns1:CheckDocumentText,omitempty"`
	}{response.CheckDocumentCode, response.CheckDocumentText}
	result, err := xml.Marshal(r)
	if err != nil {
		logrus.Printf(`error set struct to xml string: %s`, err.Error())
	}
	return string(result), err
}

func (x *ParsingXML) SetResponseAllDocsToXml(response model.GetDocumentsResponse) (string, error) {
	// create struct to ready answer xml
	r := struct {
		XMLName              xml.Name `xml:"ns1:FRDOResponse"`
		Ns1                  string   `xml:"xmlns:ns1,attr"`
		GetDocumentsResponse struct {
			Result    uint `xml:"ns1:Result,omitempty"`
			Documents struct {
				EduDocument []model.EduDocument `xml:"ns1:Document,omitempty"`
			} `xml:"ns1:Documents,omitempty"`
		} `xml:"ns1:GetDocumentsResponse,omitempty"`
	}{}
	// set params
	r.Ns1 = "urn://obrnadzor-gosuslugi-ru/frdo/1.0.1"
	docs := struct {
		EduDocument []model.EduDocument `xml:"ns1:Document,omitempty"`
	}{
		response.Documents.Document,
	}
	r.GetDocumentsResponse = struct {
		Result    uint `xml:"ns1:Result,omitempty"`
		Documents struct {
			EduDocument []model.EduDocument `xml:"ns1:Document,omitempty"`
		} `xml:"ns1:Documents,omitempty"`
	}{response.Result, docs}

	result, err := xml.Marshal(r)
	if err != nil {
		logrus.Printf(`error set struct to xml string: %s`, err.Error())
	}
	return string(result), err
}

func GetStringFromBase64(b string) (string, error) {
	r, err := base64.StdEncoding.DecodeString(b)
	return string(r), err
}

func SetStringToBase64(b string) string {
	r := base64.StdEncoding.EncodeToString([]byte(b))
	return string(r)
}
