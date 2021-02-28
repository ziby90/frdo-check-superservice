package sql

import (
	"frdo-check-superservice/model"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

var dbs map[string]*sqlx.DB

func TestMain(m *testing.M) {
	connects, err := InitDB()
	dbs = connects
	logrus.Println(dbs, err)
	m.Run()
}

func TestSearch_SearchPersonIntegrate(t *testing.T) {
	search := NewSearch(dbs)
	type test struct {
		Name string
		In   struct {
			Person     model.Person
			NameModule string
		}
		Expect []uint
	}
	tests := []test{
		test{
			Name: "OK",
			In: struct {
				Person     model.Person
				NameModule string
			}{
				Person: model.Person{
					Surename: "testing_frdo_check_surname",
					Name:     "TESTING_FRDO_CHECK_NAME",
				},
				NameModule: `DPO`,
			},
			Expect: []uint{999999999},
		},
		test{
			Name: "OK",
			In: struct {
				Person     model.Person
				NameModule string
			}{
				Person: model.Person{
					Surename: "testing_frdo_check_surname222",
					Name:     "TESTING_FRDO_CHECK_NAME222",
				},
				NameModule: `DPO`,
			},
			Expect: []uint(nil),
		},
	}
	for _, testingBody := range tests {
		var result []uint
		result, err := search.SearchPerson(testingBody.In.Person, testingBody.In.NameModule)
		if testingBody.Name == `OK` {
			if err != nil {
				t.Errorf(`error test return success: %s`, err.Error())
			}
			assert.Equal(t, testingBody.Expect, result, "not correctly return result")
		}
		if testingBody.Name == `ERROR` {
			if err == nil {
				t.Errorf(`error test return error nil`)
			}
		}
		logrus.Debug(result)
	}
}
func TestSearch_SearchDocumentIntegrate(t *testing.T) {
	search := NewSearch(dbs)
	type test struct {
		Name string
		In   struct {
			Document    model.EduDocument
			RecipientId uint
		}
		Expect []uint
	}
	series := `11111111111`
	tests := []test{
		test{
			Name: "OK",
			In: struct {
				Document    model.EduDocument
				RecipientId uint
			}{
				Document: model.EduDocument{
					Series:   &series,
					EduLevel: 6,
					Number:   `11111111111`,
				},
				RecipientId: 999999999,
			},
			Expect: []uint{999999999},
		},
		test{
			Name: "OK",
			In: struct {
				Document    model.EduDocument
				RecipientId uint
			}{
				Document: model.EduDocument{
					Series:   &series,
					EduLevel: 6,
					Number:   `111111111112222`,
				},
				RecipientId: 999999999,
			},
			Expect: []uint(nil),
		},
		test{
			Name: "ERROR",
			In: struct {
				Document    model.EduDocument
				RecipientId uint
			}{
				Document: model.EduDocument{
					Series:   &series,
					EduLevel: 7,
					Number:   `111111111112222`,
				},
				RecipientId: 999999999,
			},
			Expect: []uint(nil),
		},
	}
	for _, testingBody := range tests {
		var result []uint
		result, err := search.SearchDocument(testingBody.In.RecipientId, testingBody.In.Document)
		if testingBody.Name == `OK` {
			if err != nil {
				t.Errorf(`error test return success: %s`, err.Error())
			}
			assert.Equal(t, testingBody.Expect, result, "not correctly return result")
		}
		if testingBody.Name == `ERROR` {
			if err == nil {
				t.Errorf(`error test return error nil`)
			}
		}
		logrus.Debug(result)
	}
}
func TestSearch_SearchAllDocument(t *testing.T) {
	search := NewSearch(dbs)
	type test struct {
		Name string
		In   struct {
			RecipientId uint
			NameModule  string
		}
		Expect []model.EduDocument
	}
	series := `11111111111`
	ogrn := `301231231200013`
	fullTitle := `Индивидуальный предприниматель Фамилия013 Имя013 Отчество013`
	shortTitle := `ИП Фамилия013 И. О.`
	issueDate := `2011-02-01T00:00:00Z`
	dpoTestDocument := model.EduDocument{
		EduLevel:               6,
		Series:                 &series,
		Number:                 `11111111111`,
		OGRN:                   &ogrn,
		FullOrgName:            &fullTitle,
		ShortOrgName:           &shortTitle,
		RegisterDocumentNumber: nil,
		DocumentBlankNumber:    nil,
		DocumentIssueDate:      &issueDate,
	}

	ogrnVpo := `1027739247117`
	fullTitleVPO := `Автономная некоммерческая организация высшего профессионального образования Академический Международный Институт`
	shortTitleVPO := `АМИ`
	issueDateVPO := `2019-06-24T00:00:00Z`
	vpoTestDocument := model.EduDocument{
		EduLevel:               1,
		Series:                 &series,
		Number:                 `11111111111`,
		OGRN:                   &ogrnVpo,
		FullOrgName:            &fullTitleVPO,
		ShortOrgName:           &shortTitleVPO,
		RegisterDocumentNumber: nil,
		DocumentBlankNumber:    nil,
		DocumentIssueDate:      &issueDateVPO,
	}

	tests := []test{
		test{
			Name: "OK",
			In: struct {
				RecipientId uint
				NameModule  string
			}{
				RecipientId: 999999999,
				NameModule:  `DPO`,
			},
			Expect: []model.EduDocument{
				dpoTestDocument,
			},
		},
		test{
			Name: "OK",
			In: struct {
				RecipientId uint
				NameModule  string
			}{
				RecipientId: 999999999,
				NameModule:  `VPO`,
			},
			Expect: []model.EduDocument{
				vpoTestDocument,
			},
		},
	}
	for _, testingBody := range tests {
		var result []model.EduDocument
		result, err := search.SearchAllDocument(testingBody.In.RecipientId, testingBody.In.NameModule)
		if testingBody.Name == `OK` {
			if err != nil {
				t.Errorf(`error test return success: %s`, err.Error())
			}
			assert.Equal(t, testingBody.Expect, result, "not correctly return result")
		}
		if testingBody.Name == `ERROR` {
			if err == nil {
				t.Errorf(`error test return error nil`)
			}
		}
		logrus.Println(result)
	}
}
func TestSearch_AddNewQueue(t *testing.T) {
	type test struct {
		Name      string
		In        model.SMEVMessage
		TypeQueue string
		Expect    string
	}
	var statusResponseId uint = 3
	tests := []test{
		test{
			Name: "OK",
			In: model.SMEVMessage{
				PackageId:     0,
				SmevId:        0,
				MessageId:     "eeba6c3a-6b90-11eb-96bb-005056ae3841",
				KindActionId:  10,
				Type:          "request",
				Xml:           "PG5zMTpGUkRPUmVxdWVzdCB4bWxuczpuczE9InVybjovL29icm5hZHpvci1nb3N1c2x1Z2ktcnUvZnJkby8xLjAuMSIgeG1sbnM9InVybjovL3gtYXJ0ZWZhY3RzLXNtZXYtZ292LXJ1L3NlcnZpY2VzL21lc3NhZ2UtZXhjaGFuZ2UvdHlwZXMvMS4xIiB4bWxuczpucz0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiIHhtbG5zOm5zMj0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiPjxuczE6Q2hlY2tEb2N1bWVudFJlcXVlc3Q+PG5zMTpQZXJzb24+PG5zMTpSRkNpdGl6ZW4+dHJ1ZTwvbnMxOlJGQ2l0aXplbj48bnMxOlN1cmVuYW1lPtCf0LXRgtGA0L7QsjwvbnMxOlN1cmVuYW1lPjxuczE6TmFtZT7Qn9C10YLRgDwvbnMxOk5hbWU+PG5zMTpQYXRyb255bWljPtCf0LXRgtGA0L7QstC40Yc8L25zMTpQYXRyb255bWljPjxuczE6QmlydGhkYXRlPjE5OTgtMDItMDM8L25zMTpCaXJ0aGRhdGU+PG5zMTpHZW5kZXI+MTwvbnMxOkdlbmRlcj48L25zMTpQZXJzb24+PG5zMTpFZHVEb2N1bWVudD48bnMxOkVkdUxldmVsPjE8L25zMTpFZHVMZXZlbD48bnMxOlNlcmllcz4xMjM8L25zMTpTZXJpZXM+PG5zMTpOdW1iZXI+NDQ0IDU1NTwvbnMxOk51bWJlcj48bnMxOk9HUk4+MTEyMzQ1Njc4OTU1NTwvbnMxOk9HUk4+PG5zMTpRdWFsaWZpY2F0aW9uPtCY0L3QttC10L3QtdGAPC9uczE6UXVhbGlmaWNhdGlvbj48bnMxOkRvY3VtZW50SXNzdWVEYXRlPjIwMjAtMDYtMDE8L25zMTpEb2N1bWVudElzc3VlRGF0ZT48L25zMTpFZHVEb2N1bWVudD48L25zMTpDaGVja0RvY3VtZW50UmVxdWVzdD48L25zMTpGUkRPUmVxdWVzdD4=",
				Error:         nil,
				Status:        4,
				Ataches:       nil,
				IsTestMessage: false,
				IsReject:      false,
			},
			TypeQueue: `in`,
		},
		test{
			Name: "OK",
			In: model.SMEVMessage{
				PackageId:        0,
				SmevId:           0,
				MessageId:        "eeba6c3a-6b90-11eb-96bb-005056ae3841",
				KindActionId:     10,
				Type:             "request",
				Xml:              "PG5zMTpGUkRPUmVxdWVzdCB4bWxuczpuczE9InVybjovL29icm5hZHpvci1nb3N1c2x1Z2ktcnUvZnJkby8xLjAuMSIgeG1sbnM9InVybjovL3gtYXJ0ZWZhY3RzLXNtZXYtZ292LXJ1L3NlcnZpY2VzL21lc3NhZ2UtZXhjaGFuZ2UvdHlwZXMvMS4xIiB4bWxuczpucz0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiIHhtbG5zOm5zMj0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiPjxuczE6Q2hlY2tEb2N1bWVudFJlcXVlc3Q+PG5zMTpQZXJzb24+PG5zMTpSRkNpdGl6ZW4+dHJ1ZTwvbnMxOlJGQ2l0aXplbj48bnMxOlN1cmVuYW1lPtCf0LXRgtGA0L7QsjwvbnMxOlN1cmVuYW1lPjxuczE6TmFtZT7Qn9C10YLRgDwvbnMxOk5hbWU+PG5zMTpQYXRyb255bWljPtCf0LXRgtGA0L7QstC40Yc8L25zMTpQYXRyb255bWljPjxuczE6QmlydGhkYXRlPjE5OTgtMDItMDM8L25zMTpCaXJ0aGRhdGU+PG5zMTpHZW5kZXI+MTwvbnMxOkdlbmRlcj48L25zMTpQZXJzb24+PG5zMTpFZHVEb2N1bWVudD48bnMxOkVkdUxldmVsPjE8L25zMTpFZHVMZXZlbD48bnMxOlNlcmllcz4xMjM8L25zMTpTZXJpZXM+PG5zMTpOdW1iZXI+NDQ0IDU1NTwvbnMxOk51bWJlcj48bnMxOk9HUk4+MTEyMzQ1Njc4OTU1NTwvbnMxOk9HUk4+PG5zMTpRdWFsaWZpY2F0aW9uPtCY0L3QttC10L3QtdGAPC9uczE6UXVhbGlmaWNhdGlvbj48bnMxOkRvY3VtZW50SXNzdWVEYXRlPjIwMjAtMDYtMDE8L25zMTpEb2N1bWVudElzc3VlRGF0ZT48L25zMTpFZHVEb2N1bWVudD48L25zMTpDaGVja0RvY3VtZW50UmVxdWVzdD48L25zMTpGUkRPUmVxdWVzdD4=",
				Error:            nil,
				Status:           4,
				Ataches:          nil,
				IsTestMessage:    false,
				IsReject:         false,
				StatusResponseId: &statusResponseId,
			},
			TypeQueue: `out`,
		},
	}
	s := NewSearch(dbs)
	for _, testingBody := range tests {
		newId, err := s.AddNewQueue(testingBody.In, testingBody.TypeQueue)
		if testingBody.Name == `OK` {
			if err != nil {
				t.Errorf(`error test return success: %s`, err.Error())
			}
		}
		if testingBody.Name == `ERROR` {
			if err == nil {
				t.Errorf(`error test return error nil`)
			}
		}
		logrus.Debug(newId)
	}
}
func TestSearch_UpdateQueueLog(t *testing.T) {
	type test struct {
		Name      string
		In        model.SMEVMessage
		TypeQueue string
		Expect    string
	}
	errorMsg := `oh emoe`
	tests := []test{
		test{
			Name: "OK",
			In: model.SMEVMessage{
				Id:            1,
				PackageId:     0,
				SmevId:        0,
				MessageId:     "eeba6c3a-6b90-11eb-96bb-005056ae3841",
				KindActionId:  10,
				Type:          "request",
				Xml:           "PG5zMTpGUkRPUmVxdWVzdCB4bWxuczpuczE9InVybjovL29icm5hZHpvci1nb3N1c2x1Z2ktcnUvZnJkby8xLjAuMSIgeG1sbnM9InVybjovL3gtYXJ0ZWZhY3RzLXNtZXYtZ292LXJ1L3NlcnZpY2VzL21lc3NhZ2UtZXhjaGFuZ2UvdHlwZXMvMS4xIiB4bWxuczpucz0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiIHhtbG5zOm5zMj0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiPjxuczE6Q2hlY2tEb2N1bWVudFJlcXVlc3Q+PG5zMTpQZXJzb24+PG5zMTpSRkNpdGl6ZW4+dHJ1ZTwvbnMxOlJGQ2l0aXplbj48bnMxOlN1cmVuYW1lPtCf0LXRgtGA0L7QsjwvbnMxOlN1cmVuYW1lPjxuczE6TmFtZT7Qn9C10YLRgDwvbnMxOk5hbWU+PG5zMTpQYXRyb255bWljPtCf0LXRgtGA0L7QstC40Yc8L25zMTpQYXRyb255bWljPjxuczE6QmlydGhkYXRlPjE5OTgtMDItMDM8L25zMTpCaXJ0aGRhdGU+PG5zMTpHZW5kZXI+MTwvbnMxOkdlbmRlcj48L25zMTpQZXJzb24+PG5zMTpFZHVEb2N1bWVudD48bnMxOkVkdUxldmVsPjE8L25zMTpFZHVMZXZlbD48bnMxOlNlcmllcz4xMjM8L25zMTpTZXJpZXM+PG5zMTpOdW1iZXI+NDQ0IDU1NTwvbnMxOk51bWJlcj48bnMxOk9HUk4+MTEyMzQ1Njc4OTU1NTwvbnMxOk9HUk4+PG5zMTpRdWFsaWZpY2F0aW9uPtCY0L3QttC10L3QtdGAPC9uczE6UXVhbGlmaWNhdGlvbj48bnMxOkRvY3VtZW50SXNzdWVEYXRlPjIwMjAtMDYtMDE8L25zMTpEb2N1bWVudElzc3VlRGF0ZT48L25zMTpFZHVEb2N1bWVudD48L25zMTpDaGVja0RvY3VtZW50UmVxdWVzdD48L25zMTpGUkRPUmVxdWVzdD4=",
				Error:         &errorMsg,
				Status:        4,
				Ataches:       nil,
				IsTestMessage: false,
				IsReject:      false,
			},
			TypeQueue: `in`,
		},
	}
	s := NewSearch(dbs)
	for _, testingBody := range tests {
		err := s.UpdateQueueLog(testingBody.In, testingBody.TypeQueue)
		if testingBody.Name == `OK` {
			if err != nil {
				t.Errorf(`error test return success: %s`, err.Error())
			}
		}
		if testingBody.Name == `ERROR` {
			if err == nil {
				t.Errorf(`error test return error nil`)
			}
		}
	}
}
func TestGetStringPerson(t *testing.T) {
	type test struct {
		Name   string
		In     model.Person
		Expect struct {
			Query string
			Args  []interface{}
		}
	}
	tests := []test{
		test{
			Name: "OK",
			In: model.Person{
				Surename: "testing_frdo_check_surname",
				Name:     "TESTING_FRDO_CHECK_NAME",
			},
			Expect: struct {
				Query string
				Args  []interface{}
			}{"Select id FROM persons.recipients WHERE upper(surname)=$1 AND upper(name)=$2", []interface{}{"TESTING_FRDO_CHECK_SURNAME", "TESTING_FRDO_CHECK_NAME"}},
		},
	}
	birthday := `2020-11-16`
	patronymic := `testing_frdo_check_patronymic`
	tests = append(tests, test{
		Name: "OK",
		In: model.Person{
			Surename:   "testing_frdo_check_surname",
			Name:       "TESTING_FRDO_CHECK_NAME",
			Birthday:   &birthday,
			Patronymic: &patronymic,
		},
		Expect: struct {
			Query string
			Args  []interface{}
		}{"Select id FROM persons.recipients WHERE upper(surname)=$1 AND upper(name)=$2 AND birthday=$3 AND upper(patronymic)=$4", []interface{}{"TESTING_FRDO_CHECK_SURNAME", "TESTING_FRDO_CHECK_NAME", "2020-11-16", `TESTING_FRDO_CHECK_PATRONYMIC`}},
	})

	for _, testingBody := range tests {
		query, args := GetStringPerson(testingBody.In, "DPO")
		logrus.Println(query, args)
		if testingBody.Name == `OK` {
			assert.Equal(t, testingBody.Expect.Query, query, `not equal query`)
			assert.Equal(t, testingBody.Expect.Args, args, `not equal args`)

		}
		if testingBody.Name == `ERROR` {
			//if err == nil {
			//	t.Errorf(`error test return error nil`)
			//}
		}

	}
}
func TestGetStringDocument(t *testing.T) {
	type test struct {
		Name string
		In   struct {
			Document    model.EduDocument
			RecipientId uint
			NameModule  string
		}
		Expect struct {
			Query string
			Args  []interface{}
		}
	}
	series := `11111111111`
	tests := []test{
		test{
			Name: "OK",
			In: struct {
				Document    model.EduDocument
				RecipientId uint
				NameModule  string
			}{
				Document: model.EduDocument{
					Series: &series,
					Number: `11111111111`,
				},
				RecipientId: 999999999,
				NameModule:  `DPO`,
			},
			Expect: struct {
				Query string
				Args  []interface{}
			}{"Select a.id FROM persons.documents a JOIN admin.organizations org ON org.id=a.organization_id WHERE recipient_id=$1 AND a.number=$2 AND a.series=$3", []interface{}{uint(999999999), "11111111111", "11111111111"}},
		},
	}

	for _, testingBody := range tests {
		logrus.Println(testingBody.In.RecipientId)
		query, args := GetStringDocument(testingBody.In.Document, testingBody.In.RecipientId, testingBody.In.NameModule)
		logrus.Println(query, args)
		if testingBody.Name == `OK` {
			assert.Equal(t, testingBody.Expect.Query, query, `not equal query`)
			assert.Equal(t, testingBody.Expect.Args, args, `not equal args`)
		}
		if testingBody.Name == `ERROR` {
			//if err == nil {
			//	t.Errorf(`error test return error nil`)
			//}
		}

	}
}

//func TestSearch_SearchPersonUnitSuccess(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer db.Close()
//	mock.ExpectQuery(`Select id FROM persons.recipients`).WithArgs("ТЕСТОВОЙ", "ТЕСТ").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(4))
//	sqlxDB := sqlx.NewDb(db, "sqlmock")
//	dbs := make(map[string]*sqlx.DB)
//
//	dbs["DPO"] = sqlxDB
//	search := NewSearch(dbs)
//	// now we execute our method
//	person := model.Person{
//		Surename: "Тестовой",
//		Name:     "Тест",
//	}
//
//	res, err := search.SearchPerson(person)
//	if err != nil {
//		t.Errorf("error was not expected while updating stats: %s", err)
//	}
//	t.Log(res)
//	// we make sure that all expectations were met
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//	}
//}
//func TestSearch_SearchPersonUnitError(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer db.Close()
//	mock.ExpectQuery(`Select id FROM persons.recipients`).WithArgs("ТЕСТОВОЙ", "ТЕСТ")
//	sqlxDB := sqlx.NewDb(db, "sqlmock")
//	dbs := make(map[string]*sqlx.DB)
//
//	dbs["DPO"] = sqlxDB
//	search := NewSearch(dbs)
//	// now we execute our method
//	person := model.Person{
//		Surename: "Тестовой",
//		Name:     "Тест",
//	}
//	if _, err := search.SearchPerson(person); err == nil {
//		t.Errorf("error was not expected while updating stats: %s", err)
//	}
//	// we make sure that all expectations were met
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//	}
//}
