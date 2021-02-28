package handler

import (
	"fmt"
	"frdo-check-superservice/model"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	Task = make(chan model.SMEVMessage)
	Stop = make(chan struct{})
)

func (h *Handler) Start() {
	fmt.Println(`start daemon at`, time.Now())
	go h.Tasker()
	for workerNumber := 1; workerNumber < 3; workerNumber++ {
		go h.Worker(workerNumber)
	}
}

func (h *Handler) Tasker() {
LOOP:
	for {
		smevMessage := h.GetNewMessage()
		if smevMessage != nil {
			msg := fmt.Sprintf(`new Message from rabbit, messageId := %s, logId:= %d`, smevMessage.MessageId, smevMessage.Id)
			fmt.Println(msg)
			logrus.Println(msg)
			Task <- *smevMessage
		} else {
			time.Sleep(10 * time.Second)
		}
		select {
		case <-Stop:
			break LOOP
		default:
		}
	}
}

func (h *Handler) Worker(workerNumber int) {
LOOP:
	for {
		select {
		case <-Stop:
			break LOOP
		case smevMessage := <-Task:
			msg := fmt.Sprintf(`worker #%d get new job, messageId := %s, logId:= %d`, workerNumber, smevMessage.MessageId, smevMessage.Id)
			fmt.Println(msg)
			logrus.Println(msg)
			h.NewJob(smevMessage)
			msg = fmt.Sprintf(`worker #%d finish job, messageId := %s, logId:= %d`, workerNumber, smevMessage.MessageId, smevMessage.Id)
			fmt.Println(msg)
			logrus.Println(msg)
		default:
		}
		time.Sleep(10 * time.Second)
	}
}
