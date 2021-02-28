package model

import (
	"encoding/xml"
)

type ResultDb struct {
	IdRecipient uint
	IdDocuments []uint
	NameModule  string
}

type FRDORequest struct {
	XMLName              xml.Name `xml:"FRDORequest"`
	CheckDocumentRequest struct {
		Person      Person      `xml:"Person"`
		EduDocument EduDocument `xml:"EduDocument"`
	} `xml:"CheckDocumentRequest"`
	GetDocumentsRequest struct {
		Person Person `xml:"Person"`
	} `xml:"GetDocumentsRequest"`
}

type FRDOResponse struct {
	XMLName               xml.Name `xml:"ns1:FRDOResponse"`
	Ns1                   string   `xml:"xmlns:ns1,attr"`
	CheckDocumentResponse struct {
		CheckDocumentCode uint   `xml:"ns1:CheckDocumentCode"`
		CheckDocumentText string `xml:"ns1:CheckDocumentText"`
	} `xml:"ns1:CheckDocumentResponse,omitempty"`
	GetDocumentsResponse struct {
		Result      uint          `xml:"ns1:CheckDocumentCode,omitempty"`
		EduDocument []EduDocument `xml:"ns1:Document,omitempty"`
	} `xml:"ns1:GetDocumentsResponse,omitempty"`
}

type CheckDocumentResponse struct {
	CheckDocumentCode uint   `xml:"ns1:CheckDocumentCode"`
	CheckDocumentText string `xml:"ns1:CheckDocumentText"`
}

type GetDocumentsResponse struct {
	Result    uint `xml:"ns1:CheckDocumentCode,omitempty"`
	Documents struct {
		Document []EduDocument `xml:"ns1:Document,omitempty"`
	} `xml:"Documents,omitempty"`
}

type Person struct {
	RFCitizen  *string
	Surename   string
	Name       string
	Patronymic *string
	Birthday   *string
	BirthPlace *string
	Snils      *string
	Gender     *uint
}

type EduDocument struct {
	EduLevel               uint
	DocumentType           *string
	Series                 *string
	Number                 string
	OGRN                   *string
	KPP                    *string
	CodeEdu                *string
	NameEdu                *string
	Qualification          *string
	SpecialtyCode          *string
	Specialty              *string
	City                   *string
	FullOrgName            *string
	ShortOrgName           *string
	RegisterDocumentNumber *string
	DocumentBlankNumber    *string
	DocumentIssueDate      *string
	EduCountYear           *int
	YearStart              *int
	YearFinish             *int
}

type SMEVMessage struct {
	Id               uint     `xml:"-"`
	PackageId        int64    // Id модуля
	SmevId           int64    // Id запроса из модуля
	MessageId        string   // UUID СМЭВ сообщения
	KindActionId     int      // Id action
	Type             string   // тип сообщения "request", "response", "confirm"
	Xml              string   // base64 xml с бизнес-данными.
	Error            *string  // описание ошибки
	Status           int      // статус
	Ataches          []string // список путей файлов вложений
	IsTestMessage    bool     // признак отправки тестового сообщения
	IsReject         bool     // признак отправки в качестве ответа RequestRejected
	StatusResponseId *uint    `xml:"-"`
}
