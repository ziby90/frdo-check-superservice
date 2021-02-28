package rabbit

import (
	"10.10.11.220/ursgis/rabbitMQ.git"
	"encoding/json"
	"errors"
	"frdo-check-superservice/model"
)

type QueueWorker struct {
	ch *rabbitMQ.RabbitChan
}

func NewQueueWorker(ch *rabbitMQ.RabbitChan) *QueueWorker {
	return &QueueWorker{ch: ch}
}
func (q *QueueWorker) SetStatusMessage(requeue bool) error {
	errStatus := q.ch.SetMessageStatus(requeue)
	if errStatus != nil {
		return errors.New(errStatus.Error())
	}
	return nil
}
func (q *QueueWorker) GetNewMessage() (rabbitMQ.Message, error) {
	msg, errMsg := q.ch.GetMessage()
	if errMsg != nil {
		return msg, errors.New(errMsg.Error())
	}
	return msg, nil
}

func (q *QueueWorker) ParseMessage(body []byte) (model.SMEVMessage, error) {
	var smevMessage model.SMEVMessage
	err := json.Unmarshal(body, &smevMessage)
	if err != nil {
		//logrus.Printf(`error for unmarshalling json to smev_message struct: %s`, err.Error())
		return smevMessage, err
	}
	return smevMessage, err
}

func (q *QueueWorker) SetMessage(request model.SMEVMessage, xmlBody string) ([]byte, error) {
	var response model.SMEVMessage
	response = request
	response.Xml = xmlBody
	response.Type = `response`
	body, err := json.Marshal(response)
	if err != nil {
		return body, err
	}
	return body, err
}

func (q *QueueWorker) PublishMessage(msg []byte) error {
	q.ch.SetMessageMeta("FRDO", "", "", "")
	err := q.ch.Publish(`OUT`, msg)
	//logrus.Printf(`error: %v`, err)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
