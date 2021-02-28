package xmlpars

import (
	"encoding/xml"
	"frdo-check-superservice/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetStringFromBase64(t *testing.T) {
	type test struct {
		Name   string
		In     string
		Expect string
	}
	tests := []test{
		test{
			Name:   "OK",
			In:     "PG5zMTpGUkRPUmVxdWVzdCB4bWxuczpuczE9InVybjovL29icm5hZHpvci1nb3N1c2x1Z2ktcnUvZnJkby8xLjAuMSIgeG1sbnM9InVybjovL3gtYXJ0ZWZhY3RzLXNtZXYtZ292LXJ1L3NlcnZpY2VzL21lc3NhZ2UtZXhjaGFuZ2UvdHlwZXMvMS4xIiB4bWxuczpucz0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiIHhtbG5zOm5zMj0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiPjxuczE6Q2hlY2tEb2N1bWVudFJlcXVlc3Q+PG5zMTpQZXJzb24+PG5zMTpSRkNpdGl6ZW4+dHJ1ZTwvbnMxOlJGQ2l0aXplbj48bnMxOlN1cmVuYW1lPtCf0LXRgtGA0L7QsjwvbnMxOlN1cmVuYW1lPjxuczE6TmFtZT7Qn9C10YLRgDwvbnMxOk5hbWU+PG5zMTpQYXRyb255bWljPtCf0LXRgtGA0L7QstC40Yc8L25zMTpQYXRyb255bWljPjxuczE6QmlydGhkYXRlPjE5OTgtMDItMDM8L25zMTpCaXJ0aGRhdGU+PG5zMTpHZW5kZXI+MTwvbnMxOkdlbmRlcj48L25zMTpQZXJzb24+PG5zMTpFZHVEb2N1bWVudD48bnMxOkVkdUxldmVsPjE8L25zMTpFZHVMZXZlbD48bnMxOlNlcmllcz4xMjM8L25zMTpTZXJpZXM+PG5zMTpOdW1iZXI+NDQ0IDU1NTwvbnMxOk51bWJlcj48bnMxOk9HUk4+MTEyMzQ1Njc4OTU1NTwvbnMxOk9HUk4+PG5zMTpRdWFsaWZpY2F0aW9uPtCY0L3QttC10L3QtdGAPC9uczE6UXVhbGlmaWNhdGlvbj48bnMxOkRvY3VtZW50SXNzdWVEYXRlPjIwMjAtMDYtMDE8L25zMTpEb2N1bWVudElzc3VlRGF0ZT48L25zMTpFZHVEb2N1bWVudD48L25zMTpDaGVja0RvY3VtZW50UmVxdWVzdD48L25zMTpGUkRPUmVxdWVzdD4=",
			Expect: "<ns1:FRDORequest xmlns:ns1=\"urn://obrnadzor-gosuslugi-ru/frdo/1.0.1\" xmlns=\"urn://x-artefacts-smev-gov-ru/services/message-exchange/types/1.1\" xmlns:ns=\"urn://x-artefacts-smev-gov-ru/services/message-exchange/types/basic/1.1\" xmlns:ns2=\"urn://x-artefacts-smev-gov-ru/services/message-exchange/types/basic/1.1\"><ns1:CheckDocumentRequest><ns1:Person><ns1:RFCitizen>true</ns1:RFCitizen><ns1:Surename>Петров</ns1:Surename><ns1:Name>Петр</ns1:Name><ns1:Patronymic>Петрович</ns1:Patronymic><ns1:Birthdate>1998-02-03</ns1:Birthdate><ns1:Gender>1</ns1:Gender></ns1:Person><ns1:EduDocument><ns1:EduLevel>1</ns1:EduLevel><ns1:Series>123</ns1:Series><ns1:Number>444 555</ns1:Number><ns1:OGRN>1123456789555</ns1:OGRN><ns1:Qualification>Инженер</ns1:Qualification><ns1:DocumentIssueDate>2020-06-01</ns1:DocumentIssueDate></ns1:EduDocument></ns1:CheckDocumentRequest></ns1:FRDORequest>",
		},
		//test{
		//	Name: "ERROR",
		//	In:   "pGUkRPUmVxdWVzdCB4bWxuczpuczE9InVybjovL29icm5hZHpvci1nb3N1c2x1Z2ktcnUvZnJkby8xLjAuMSIgeG1sbnM9InVybjovL3gtYXJ0ZWZhY3RzLXNtZXYtZ292LXJ1L3NlcnZpY2VzL21lc3NhZ2UtZXhjaGFuZ2UvdHlwZXMvMS4xIiB4bWxuczpucz0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiIHhtbG5zOm5zMj0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiPjxuczE6Q2hlY2tEb2N1bWVudFJlcXVlc3Q+PG5zMTpQZXJzb24+PG5zMTpSRkNpdGl6ZW4+dHJ1ZTwvbnMxOlJGQ2l0aXplbj48bnMxOlN1cmVuYW1lPtCf0LXRgtGA0L7QsjwvbnMxOlN1cmVuYW1lPjxuczE6TmFtZT7Qn9C10YLRgDwvbnMxOk5hbWU+PG5zMTpQYXRyb255bWljPtCf0LXRgtGA0L7QstC40Yc8L25zMTpQYXRyb255bWljPjxuczE6QmlydGhkYXRlPjE5OTgtMDItMDM8L25zMTpCaXJ0aGRhdGU+PG5zMTpHZW5kZXI+MTwvbnMxOkdlbmRlcj48L25zMTpQZXJzb24+PG5zMTpFZHVEb2N1bWVudD48bnMxOkVkdUxldmVsPjE8L25zMTpFZHVMZXZlbD48bnMxOlNlcmllcz4xMjM8L25zMTpTZXJpZXM+PG5zMTpOdW1iZXI+NDQ0IDU1NTwvbnMxOk51bWJlcj48bnMxOk9HUk4+MTEyMzQ1Njc4OTU1NTwvbnMxOk9HUk4+PG5zMTpRdWFsaWZpY2F0aW9uPtCY0L3QttC10L3QtdGAPC9uczE6UXVhbGlmaWNhdGlvbj48bnMxOkRvY3VtZW50SXNzdWVEYXRlPjIwMjAtMDYtMDE8L25zMTpEb2N1bWVudElzc3VlRGF0ZT48L25zMTpFZHVEb2N1bWVudD48L25zMTpDaGVja0RvY3VtZW50UmVxdWVzdD48L25zMTpGUkRPUmVxdWVzdD4=",
		//},
	}
	for _, testingBody := range tests {
		res, err := GetStringFromBase64(testingBody.In)
		logrus.Printf(`result : %s`, res)
		if testingBody.Name == `OK` {
			if err != nil {
				t.Errorf(`error test return success: %s`, err.Error())
			}
			if testingBody.Expect != `` {
				assert.Equal(t, testingBody.Expect, res, `not a valid expected params`)
			}

		}
		if testingBody.Name == `ERROR` {
			if err == nil {
				t.Errorf(`error test return error nil`)
			}
		}
		logrus.Printf(`succes test - %s`, testingBody.Name)
	}

}
func TestSetStringToBase64(t *testing.T) {
	type test struct {
		Name   string
		In     string
		Expect string
	}
	tests := []test{
		test{
			Name:   "OK",
			Expect: "PG5zMTpGUkRPUmVxdWVzdCB4bWxuczpuczE9InVybjovL29icm5hZHpvci1nb3N1c2x1Z2ktcnUvZnJkby8xLjAuMSIgeG1sbnM9InVybjovL3gtYXJ0ZWZhY3RzLXNtZXYtZ292LXJ1L3NlcnZpY2VzL21lc3NhZ2UtZXhjaGFuZ2UvdHlwZXMvMS4xIiB4bWxuczpucz0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiIHhtbG5zOm5zMj0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiPjxuczE6Q2hlY2tEb2N1bWVudFJlcXVlc3Q+PG5zMTpQZXJzb24+PG5zMTpSRkNpdGl6ZW4+dHJ1ZTwvbnMxOlJGQ2l0aXplbj48bnMxOlN1cmVuYW1lPtCf0LXRgtGA0L7QsjwvbnMxOlN1cmVuYW1lPjxuczE6TmFtZT7Qn9C10YLRgDwvbnMxOk5hbWU+PG5zMTpQYXRyb255bWljPtCf0LXRgtGA0L7QstC40Yc8L25zMTpQYXRyb255bWljPjxuczE6QmlydGhkYXRlPjE5OTgtMDItMDM8L25zMTpCaXJ0aGRhdGU+PG5zMTpHZW5kZXI+MTwvbnMxOkdlbmRlcj48L25zMTpQZXJzb24+PG5zMTpFZHVEb2N1bWVudD48bnMxOkVkdUxldmVsPjE8L25zMTpFZHVMZXZlbD48bnMxOlNlcmllcz4xMjM8L25zMTpTZXJpZXM+PG5zMTpOdW1iZXI+NDQ0IDU1NTwvbnMxOk51bWJlcj48bnMxOk9HUk4+MTEyMzQ1Njc4OTU1NTwvbnMxOk9HUk4+PG5zMTpRdWFsaWZpY2F0aW9uPtCY0L3QttC10L3QtdGAPC9uczE6UXVhbGlmaWNhdGlvbj48bnMxOkRvY3VtZW50SXNzdWVEYXRlPjIwMjAtMDYtMDE8L25zMTpEb2N1bWVudElzc3VlRGF0ZT48L25zMTpFZHVEb2N1bWVudD48L25zMTpDaGVja0RvY3VtZW50UmVxdWVzdD48L25zMTpGUkRPUmVxdWVzdD4=",
			In:     "<ns1:FRDORequest xmlns:ns1=\"urn://obrnadzor-gosuslugi-ru/frdo/1.0.1\" xmlns=\"urn://x-artefacts-smev-gov-ru/services/message-exchange/types/1.1\" xmlns:ns=\"urn://x-artefacts-smev-gov-ru/services/message-exchange/types/basic/1.1\" xmlns:ns2=\"urn://x-artefacts-smev-gov-ru/services/message-exchange/types/basic/1.1\"><ns1:CheckDocumentRequest><ns1:Person><ns1:RFCitizen>true</ns1:RFCitizen><ns1:Surename>Петров</ns1:Surename><ns1:Name>Петр</ns1:Name><ns1:Patronymic>Петрович</ns1:Patronymic><ns1:Birthdate>1998-02-03</ns1:Birthdate><ns1:Gender>1</ns1:Gender></ns1:Person><ns1:EduDocument><ns1:EduLevel>1</ns1:EduLevel><ns1:Series>123</ns1:Series><ns1:Number>444 555</ns1:Number><ns1:OGRN>1123456789555</ns1:OGRN><ns1:Qualification>Инженер</ns1:Qualification><ns1:DocumentIssueDate>2020-06-01</ns1:DocumentIssueDate></ns1:EduDocument></ns1:CheckDocumentRequest></ns1:FRDORequest>",
		},
	}
	for _, testingBody := range tests {
		res := SetStringToBase64(testingBody.In)
		logrus.Printf(`result : %s`, res)
		if testingBody.Expect != `` {
			assert.Equal(t, testingBody.Expect, res, `not a valid expected params`)
		}
		logrus.Printf(`succes test - %s`, testingBody.Name)
	}

}
func TestParsingXML_ParseXmlToStruct(t *testing.T) {
	type test struct {
		Name   string
		In     string
		Expect string
	}
	tests := []test{
		test{
			Name: "OK",
			In:   "<?xml version=\"1.0\" encoding=\"utf-8\"?><ns1:FRDORequest xmlns:ns1=\"urn://obrnadzor-gosuslugi-ru/frdo/1.0.1\" xmlns=\"urn://x-artefacts-smev-gov-ru/services/message-exchange/types/1.1\" xmlns:ns=\"urn://x-artefacts-smev-gov-ru/services/message-exchange/types/basic/1.1\" xmlns:ns2=\"urn://x-artefacts-smev-gov-ru/services/message-exchange/types/basic/1.1\"><ns1:CheckDocumentRequest><ns1:Person><ns1:RFCitizen>true</ns1:RFCitizen><ns1:Surename>Петров</ns1:Surename><ns1:Name>Петр</ns1:Name><ns1:Patronymic>Петрович</ns1:Patronymic><ns1:Birthdate>1998-02-03</ns1:Birthdate><ns1:Gender>1</ns1:Gender></ns1:Person><ns1:EduDocument><ns1:EduLevel>1</ns1:EduLevel><ns1:Series>123</ns1:Series><ns1:Number>444 555</ns1:Number><ns1:OGRN>1123456789555</ns1:OGRN><ns1:Qualification>Инженер</ns1:Qualification><ns1:DocumentIssueDate>2020-06-01</ns1:DocumentIssueDate></ns1:EduDocument></ns1:CheckDocumentRequest></ns1:FRDORequest>",
		},
		test{
			Name: "ERROR",
			In:   "<<<ns1:FRDORequest xmlns:ns1=\"urn://obrnadzor-gosuslugi-ru/frdo/1.0.1\" xmlns=\"urn://x-artefacts-smev-gov-ru/services/message-exchange/types/1.1\" xmlns:ns=\"urn://x-artefacts-smev-gov-ru/services/message-exchange/types/basic/1.1\" xmlns:ns2=\"urn://x-artefacts-smev-gov-ru/services/message-exchange/types/basic/1.1\"><ns1:CheckDocumentRequest><ns1:Person><ns1:RFCitizen>true</ns1:RFCitizen><ns1:Surename>Петров</ns1:Surename><ns1:Name>Петр</ns1:Name><ns1:Patronymic>Петрович</ns1:Patronymic><ns1:Birthdate>1998-02-03</ns1:Birthdate><ns1:Gender>1</ns1:Gender></ns1:Person><ns1:EduDocument><ns1:EduLevel>1</ns1:EduLevel><ns1:Series>123</ns1:Series><ns1:Number>444 555</ns1:Number><ns1:OGRN>1123456789555</ns1:OGRN><ns1:Qualification>Инженер</ns1:Qualification><ns1:DocumentIssueDate>2020-06-01</ns1:DocumentIssueDate></ns1:EduDocument></ns1:CheckDocumentRequest></ns1:FRDORequest>",
		},
	}
	x := NewParsingXML()
	for _, testingBody := range tests {
		res, err := x.ParseXmlToStruct(testingBody.In)
		logrus.Printf(`result : %v`, res)
		if testingBody.Name == `OK` {
			if err != nil {
				t.Errorf(`error test return success: %s`, err.Error())
			}
			if testingBody.Expect != `` {
				assert.Equal(t, testingBody.Expect, res, `not a valid expected params`)
			}

		}
		if testingBody.Name == `ERROR` {
			if err == nil {
				t.Errorf(`error test return error nil`)
			}
		}
		logrus.Printf(`succes test - %s`, testingBody.Name)
	}

}
func TestParsingXML_SetResponseToXml(t *testing.T) {
	type test struct {
		Name   string
		In     model.FRDOResponse
		Expect string
	}
	tests := []test{
		test{
			Name: "OK",
			In: model.FRDOResponse{
				XMLName: xml.Name{},
				Ns1:     "urn://obrnadzor-gosuslugi-ru/frdo/1.0.1",
				CheckDocumentResponse: struct {
					CheckDocumentCode uint   `xml:"ns1:CheckDocumentCode"`
					CheckDocumentText string `xml:"ns1:CheckDocumentText"`
				}{3232, "awf"},
			},
		},
		test{
			Name: "Error",
			In: model.FRDOResponse{
				XMLName: xml.Name{},
				Ns1:     "urn://obrnadzor-gosuslugi-ru/frdo/1.0.1",
				CheckDocumentResponse: struct {
					CheckDocumentCode uint   `xml:"ns1:CheckDocumentCode"`
					CheckDocumentText string `xml:"ns1:CheckDocumentText"`
				}{},
			},
		},
	}
	x := NewParsingXML()
	for _, testingBody := range tests {
		res, err := x.SetResponseToXml(testingBody.In)
		logrus.Printf(`result : %v`, res)
		if testingBody.Name == `OK` {
			if err != nil {
				t.Errorf(`error test return success: %s`, err.Error())
			}
			if testingBody.Expect != `` {
				assert.Equal(t, testingBody.Expect, res, `not a valid expected params`)
			}

		}
		if testingBody.Name == `ERROR` {
			if err == nil {
				t.Errorf(`error test return error nil`)
			}
		}
		logrus.Printf(`succes test - %s`, testingBody.Name)
	}

}
