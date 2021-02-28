package service

import (
	"frdo-check-superservice/model"
	"frdo-check-superservice/pkg/repository"
)

type QueueService struct {
	repo repository.QueueWorker
}

func NewQueueService(repo repository.QueueWorker) *QueueService {
	return &QueueService{repo: repo}
}

func (q *QueueService) GetNewMessage() (model.SMEVMessage, error) {
	var result model.SMEVMessage
	msg, err := q.repo.GetNewMessage()
	if err != nil {
		return result, err
	}
	err = q.repo.SetStatusMessage(false)
	if err != nil {
		return result, err
	}
	result, err = q.repo.ParseMessage(msg.Body)
	if err != nil {
		return result, err
	}
	return result, err
}

func (q *QueueService) PostRequest(request model.SMEVMessage, xmlBody string) error {
	result, err := q.repo.SetMessage(request, xmlBody)
	if err != nil {
		return err
	}
	err = q.repo.PublishMessage(result)
	return err
}
