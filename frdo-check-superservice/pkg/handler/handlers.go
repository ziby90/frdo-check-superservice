package handler

import (
	"frdo-check-superservice/model"
	"frdo-check-superservice/pkg/entities"
	"frdo-check-superservice/pkg/service"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	services *service.Service
	entities *entities.Entities
}

func NewHandler(services *service.Service, entities *entities.Entities) *Handler {
	return &Handler{services: services, entities: entities}
}

func (h *Handler) NewJob(smevMessage model.SMEVMessage) {
	var err error
	if smevMessage.Type == `request` {
		err = h.Job(smevMessage)
	}
	logrus.Printf(`finish parsing job, messageId := %s`, smevMessage.MessageId)
	if err != nil {
		msg := err.Error()
		logrus.Errorf(`error for worker smevmessage: %s`, err.Error())
		smevMessage.Error = &msg
		err = h.services.UpdateQueueLog(smevMessage, `IN`)
		if err != nil {
			logrus.Errorf(`error for loging error in db: %s`, err.Error())
		}
	}
}

func (h *Handler) Job(smevMessage model.SMEVMessage) error {

	// parse xml smev message to struct
	request, err := h.services.ParseXmlToStruct(smevMessage.Xml)
	if err != nil {
		return err
	}
	var xmlBody string
	var code uint

	check, err := h.entities.NewKindActionType(smevMessage.KindActionId)
	if err != nil {
		return err
	}
	xmlBody, code, err = check.Search(request)
	if err != nil {
		return err
	}
	smevMessage.StatusResponseId = &code

	// create rabbit message
	// publish rabbit message
	err = h.PostRequest(smevMessage, xmlBody)
	if err != nil {
		return err
	}
	smevMessage.Xml = xmlBody
	// add queue log in db
	_, err = h.AddNewQueueLog(smevMessage, `OUT`)
	if err != nil {
		return err
	}
	//logrus.Printf(`new log id OUT :%d`, logOutId)
	return nil
}

func (h *Handler) AddNewQueueLog(smevMessage model.SMEVMessage, typeRequest string) (uint, error) {
	logInId, err := h.services.AddNewQueueLog(smevMessage, typeRequest)
	return logInId, err
}

func (h *Handler) GetNewMessage() *model.SMEVMessage {
	smevMessage, err := h.services.GetNewMessage()
	if err != nil {
		logrus.Errorf(`error for get message" %s"`, err.Error())
		return nil
	}
	// insert in db logs smev message
	if smevMessage.MessageId == `` {
		logrus.Errorf(`error for get message, empty message_id`)
		return nil
	}
	logId, err := h.AddNewQueueLog(smevMessage, `IN`)
	if err != nil {
		logrus.Errorf(err.Error())
		return nil
	}
	smevMessage.Id = logId
	//logrus.Printf(`new log id IN :%d`, logInId)
	return &smevMessage
}

func (h *Handler) ParseXmlToStruct(smevMessage model.SMEVMessage) (model.FRDORequest, error) {
	request, err := h.services.ParseXmlToStruct(smevMessage.Xml)
	return request, err
}

func (h *Handler) PostRequest(smevMessage model.SMEVMessage, xmlBody string) error {
	err := h.services.PostRequest(smevMessage, xmlBody)
	return err
}
