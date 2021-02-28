package service

import (
	"frdo-check-superservice/pkg/repository"
	mock_repository "frdo-check-superservice/pkg/repository/mocks"
	"frdo-check-superservice/pkg/repository/rabbit"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestQueueService_GetNewMessageUnit(t *testing.T) {
	//s, m := getServiceMockQueue(t)
	//m.EXPECT().GetNewMessage().Return(rabbitMQ.Message{}, nil)

}

func getNewServiceForTestQueue(t *testing.T) *Service {
	ch, err := rabbit.InitRabbit()
	if err != nil {
		t.Errorf("error for init rabbit: %s", err.Error())
	}

	errConsume := ch.SetConsumerSettings("IN")
	if errConsume != nil {
		t.Errorf("error for set consumer: %s", errConsume)
	}

	logrus.Println(`connect to rabbit for reading`)
	repos := repository.NewRepository(nil, ch)
	s := NewService(repos)
	return s
}

func getServiceMockQueue(t *testing.T) (*Service, *mock_repository.MockQueueWorker) {
	ctrl := gomock.NewController(t)
	m := mock_repository.NewMockQueueWorker(ctrl)
	r := repository.Repository{
		SearchInDb:  mock_repository.NewMockSearchInDb(ctrl),
		QueueWorker: m,
		ParsingXML:  mock_repository.NewMockParsingXML(ctrl),
	}
	s := NewService(&r)
	return s, m
}
