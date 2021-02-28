package handler

import (
	"frdo-check-superservice/pkg/repository"
	"frdo-check-superservice/pkg/repository/rabbit"
	"frdo-check-superservice/pkg/repository/sql"
	"frdo-check-superservice/pkg/service"
	"testing"
)

func TestMain(m *testing.M) {

}

func TestHandler_Start(t *testing.T) {

}

func getNewService(t *testing.T) *service.Service {
	ch, errCh := rabbit.InitRabbit()
	if errCh != nil {
		t.Errorf(`failed init rabbit connection: %s`, errCh.Error())
	}
	dbs, errDb := sql.InitDB()
	if errDb != nil {
		t.Errorf(`failed init db connection: %s`, errDb.Error())
	}
	repos := repository.NewRepository(dbs, ch)
	s := service.NewService(repos)
	return s
}
