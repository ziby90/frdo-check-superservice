package service

import (
	"frdo-check-superservice/model"
	"frdo-check-superservice/pkg/repository"
	mock_repository "frdo-check-superservice/pkg/repository/mocks"
	"frdo-check-superservice/pkg/repository/sql"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestSearchService_SearchUnit(t *testing.T) {
	s, m := getServiceMockSearch(t)
	person := model.Person{
		Surename: "testing_frdo_check_surname",
		Name:     "TESTING_FRDO_CHECK_NAME",
	}
	m.EXPECT().SearchPerson(person, `DPO`).Return([]uint{999999999}, nil)
	series := `111111111112222`
	document := model.EduDocument{
		Series:   &series,
		EduLevel: 6,
		Number:   `111111111112222`,
	}
	m.EXPECT().SearchDocument(uint(999999999), document).Return([]uint{999999999}, nil)
	res, err := s.Search(person, document)
	expect := []model.ResultDb{
		model.ResultDb{
			IdRecipient: 999999999,
			IdDocuments: []uint{999999999},
			NameModule:  "DPO",
		},
	}
	if err != nil {
		t.Errorf(`error: %s`, err.Error())
	}
	assert.Equal(t, expect, res, `not equal result`)
}
func TestSearchService_GetDocumentsIntegrate(t *testing.T) {
	s := getNewServiceForTestSearch()
	type test struct {
		Name string
		In   struct {
			Person model.Person
		}
		Expect []model.EduDocument
	}

	tests := []test{
		test{
			Name: "OK",
			In: struct {
				Person model.Person
			}{
				Person: model.Person{
					Surename: "testing_frdo_check_surname",
					Name:     "TESTING_FRDO_CHECK_NAME",
				},
			},
		},
	}
	for _, testingBody := range tests {
		res, err := s.GetDocuments(testingBody.In.Person)
		if testingBody.Name == `OK` {
			if err != nil {
				t.Errorf(`error test return success: %s`, err.Error())
			}
			if len(testingBody.Expect) > 0 {
				assert.Equal(t, testingBody.Expect, res, "not correctly return result")
			}
		}
		if testingBody.Name == `ERROR` {
			if err == nil {
				t.Errorf(`error test return error nil`)
			}
		}
		logrus.Println(res)
	}
}
func TestSearchService_SearchIntegrate(t *testing.T) {
	s := getNewServiceForTestSearch()
	type test struct {
		Name string
		In   struct {
			Person   model.Person
			Document model.EduDocument
		}
		Expect []model.ResultDb
	}
	series := `11111111111`
	tests := []test{
		test{
			Name: "OK",
			In: struct {
				Person   model.Person
				Document model.EduDocument
			}{
				Person: model.Person{
					Surename: "testing_frdo_check_surname",
					Name:     "TESTING_FRDO_CHECK_NAME",
				},
				Document: model.EduDocument{
					Series:   &series,
					Number:   `11111111111`,
					EduLevel: 6,
				},
			},
			Expect: []model.ResultDb{
				model.ResultDb{
					IdRecipient: 999999999,
					IdDocuments: []uint{999999999},
					NameModule:  "DPO",
				},
			},
		},
	}
	for _, testingBody := range tests {
		res, err := s.Search(testingBody.In.Person, testingBody.In.Document)
		if testingBody.Name == `OK` {
			if err != nil {
				t.Errorf(`error test return success: %s`, err.Error())
			}
			if len(testingBody.Expect) > 0 {
				assert.Equal(t, testingBody.Expect, res, "not correctly return result")
			}
		}
		if testingBody.Name == `ERROR` {
			if err == nil {
				t.Errorf(`error test return error nil`)
			}
		}
		logrus.Debug(res)
	}
}

func TestSearchService_AddNewQueueLogUnit(t *testing.T) {
	s, m := getServiceMockSearch(t)
	msg := model.SMEVMessage{}
	m.EXPECT().AddNewQueue(msg, `IN`).Return(uint(1), nil)
	res, err := s.AddNewQueueLog(msg, `IN`)
	expect := uint(1)
	if err != nil {
		t.Errorf(`error: %s`, err.Error())
	}
	assert.Equal(t, expect, res, `not equal result`)
}

func getNewServiceForTestSearch() *Service {
	dbs, err := sql.InitDB()
	logrus.Println(dbs, err)
	repos := repository.NewRepository(dbs, nil)
	s := NewService(repos)
	return s
}
func getServiceMockSearch(t *testing.T) (*Service, *mock_repository.MockSearchInDb) {
	ctrl := gomock.NewController(t)
	m := mock_repository.NewMockSearchInDb(ctrl)
	eduLevels := make(map[int]string)
	eduLevels[6] = `DPO`
	var db *sqlx.DB
	m.EXPECT().GetEducationLevel().Return(eduLevels).AnyTimes()
	m.EXPECT().GetDbConnections().Return(map[string]*sqlx.DB{
		`DPO`: db,
	})

	r := repository.Repository{
		SearchInDb:  m,
		QueueWorker: mock_repository.NewMockQueueWorker(ctrl),
		ParsingXML:  mock_repository.NewMockParsingXML(ctrl),
	}
	s := NewService(&r)
	return s, m
}
