package rabbit

import (
	"10.10.11.220/ursgis/rabbitMQ.git"
	"fmt"
	"frdo-check-superservice/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestQueueWorker_GetNewMessageIntegrate(t *testing.T) {
	ch, err := InitRabbit()
	if err != nil {
		fmt.Println(`err`, err)
		t.Errorf("error for init rabbit: %s", err.Error())
	}

	queue := NewQueueWorker(ch)
	msg, errQueue := queue.GetNewMessage()
	if errQueue != nil {
		t.Errorf("error for get message: %s", errQueue.Error())
	}
	logrus.Println(string(msg.Body))
	ch.Close()
}
func TestQueueWorker_PublishMessage(t *testing.T) {
	ch, err := InitRabbit()
	if err != nil {
		t.Errorf("error for init rabbit: %s", err.Error())
	}

	type test struct {
		Name   string
		In     string
		Expect string
	}
	tests := []test{
		test{
			Name: "OK",
			In:   `{"Id":0,"PackageId":12,"SmevId":0,"MessageId":"7f91173e-6ace-11eb-bad5-005056ae3841","KindActionId":10,"Type":"response","Xml":"PG5zMTpGUkRPUmVxdWVzdCB4bWxuczpuczE9InVybjovL29icm5hZHpvci1nb3N1c2x1Z2ktcnUvZnJkby8xLjAuMSIgeG1sbnM9InVybjovL3gtYXJ0ZWZhY3RzLXNtZXYtZ292LXJ1L3NlcnZpY2VzL21lc3NhZ2UtZXhjaGFuZ2UvdHlwZXMvMS4xIiB4bWxuczpucz0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiIHhtbG5zOm5zMj0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiPjxuczE6Q2hlY2tEb2N1bWVudFJlcXVlc3Q+PG5zMTpQZXJzb24+PG5zMTpSRkNpdGl6ZW4+dHJ1ZTwvbnMxOlJGQ2l0aXplbj48bnMxOlN1cmVuYW1lPtCf0LXRgtGA0L7QsjwvbnMxOlN1cmVuYW1lPjxuczE6TmFtZT7Qn9C10YLRgDwvbnMxOk5hbWU+PG5zMTpQYXRyb255bWljPtCf0LXRgtGA0L7QstC40Yc8L25zMTpQYXRyb255bWljPjxuczE6QmlydGhkYXRlPjE5OTgtMDItMDM8L25zMTpCaXJ0aGRhdGU+PG5zMTpHZW5kZXI+MTwvbnMxOkdlbmRlcj48L25zMTpQZXJzb24+PG5zMTpFZHVEb2N1bWVudD48bnMxOkVkdUxldmVsPjE8L25zMTpFZHVMZXZlbD48bnMxOlNlcmllcz4xMjM8L25zMTpTZXJpZXM+PG5zMTpOdW1iZXI+NDQ0IDU1NTwvbnMxOk51bWJlcj48bnMxOk9HUk4+MTEyMzQ1Njc4OTU1NTwvbnMxOk9HUk4+PG5zMTpRdWFsaWZpY2F0aW9uPtCY0L3QttC10L3QtdGAPC9uczE6UXVhbGlmaWNhdGlvbj48bnMxOkRvY3VtZW50SXNzdWVEYXRlPjIwMjAtMDYtMDE8L25zMTpEb2N1bWVudElzc3VlRGF0ZT48L25zMTpFZHVEb2N1bWVudD48L25zMTpDaGVja0RvY3VtZW50UmVxdWVzdD48L25zMTpGUkRPUmVxdWVzdD4=","Error":"","Status":4,"Ataches":null,"IsTestMessage":false,"IsReject":false,"StatusResponseId":1}`,
		},
	}
	q := NewQueueWorker(ch)
	for _, testingBody := range tests {
		msg := []byte(testingBody.In)
		err := q.PublishMessage(msg)
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
	q.ch.Close()
}
func TestQueueWorker_ParseMessage(t *testing.T) {
	type test struct {
		Name   string
		In     string
		Expect string
	}
	tests := []test{
		test{
			Name: "OK",
			In:   `{"PackageId":12,"SmevId3":0,"MessageId":"7f91173e-6ace-11eb-bad5-005056ae3841","KindActionId":10,"Type":"request","Xml":"PG5zMTpGUkRPUmVxdWVzdCB4bWxuczpuczE9InVybjovL29icm5hZHpvci1nb3N1c2x1Z2ktcnUvZnJkby8xLjAuMSIgeG1sbnM9InVybjovL3gtYXJ0ZWZhY3RzLXNtZXYtZ292LXJ1L3NlcnZpY2VzL21lc3NhZ2UtZXhjaGFuZ2UvdHlwZXMvMS4xIiB4bWxuczpucz0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiIHhtbG5zOm5zMj0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiPgoJPG5zMTpDaGVja0RvY3VtZW50UmVxdWVzdD4KCQk8bnMxOlBlcnNvbj4JCQkKCQkJPG5zMTpTdXJlbmFtZT50ZXN0aW5nX2ZyZG9fY2hlY2tfc3VybmFtZTwvbnMxOlN1cmVuYW1lPgoJCQk8bnMxOk5hbWU+VEVTVElOR19GUkRPX0NIRUNLX05BTUU8L25zMTpOYW1lPgoJCQk8bnMxOlBhdHJvbnltaWM+dGVzdGluZ19mcmRvX2NoZWNrX3BhdHJvbnltaWM8L25zMTpQYXRyb255bWljPgoJCQk8bnMxOkJpcnRoZGF0ZT4xOTgwLTAyLTAxPC9uczE6QmlydGhkYXRlPgoJCQk8bnMxOkdlbmRlcj4xPC9uczE6R2VuZGVyPgoJCTwvbnMxOlBlcnNvbj4KCQk8bnMxOkVkdURvY3VtZW50PgoJCQk8bnMxOkVkdUxldmVsPjY8L25zMTpFZHVMZXZlbD4KCQkJPG5zMTpTZXJpZXM+MTExMTExMTExMTE8L25zMTpTZXJpZXM+CgkJCTxuczE6TnVtYmVyPjExMTExMTExMTExPC9uczE6TnVtYmVyPgoJCQk8bnMxOlF1YWxpZmljYXRpb24+0JjQvdC20LXQvdC10YA8L25zMTpRdWFsaWZpY2F0aW9uPgoJCQk8bnMxOkRvY3VtZW50SXNzdWVEYXRlPjIwMTEtMDItMDE8L25zMTpEb2N1bWVudElzc3VlRGF0ZT4KCQk8L25zMTpFZHVEb2N1bWVudD4KCTwvbnMxOkNoZWNrRG9jdW1lbnRSZXF1ZXN0Pgo8L25zMTpGUkRPUmVxdWVzdD4=","Error":"","Status":4,"Ataches":null,"IsTestMessage":false,"IsReject":false,"StatusResponseId":null}`,
		},
		test{
			Name: "ERROR",
			In:   `{"PackageId":"12ewq","SmevId3":0,"MessageId":"7f91173e-6ace-11eb-bad5-005056ae3841","KindActionId":10,"Type":"request","Xml":"PG5zMTpGUkRPUmVxdWVzdCB4bWxuczpuczE9InVybjovL29icm5hZHpvci1nb3N1c2x1Z2ktcnUvZnJkby8xLjAuMSIgeG1sbnM9InVybjovL3gtYXJ0ZWZhY3RzLXNtZXYtZ292LXJ1L3NlcnZpY2VzL21lc3NhZ2UtZXhjaGFuZ2UvdHlwZXMvMS4xIiB4bWxuczpucz0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiIHhtbG5zOm5zMj0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiPjxuczE6Q2hlY2tEb2N1bWVudFJlcXVlc3Q+PG5zMTpQZXJzb24+PG5zMTpSRkNpdGl6ZW4+dHJ1ZTwvbnMxOlJGQ2l0aXplbj48bnMxOlN1cmVuYW1lPtCf0LXRgtGA0L7QsjwvbnMxOlN1cmVuYW1lPjxuczE6TmFtZT7Qn9C10YLRgDwvbnMxOk5hbWU+PG5zMTpQYXRyb255bWljPtCf0LXRgtGA0L7QstC40Yc8L25zMTpQYXRyb255bWljPjxuczE6QmlydGhkYXRlPjE5OTgtMDItMDM8L25zMTpCaXJ0aGRhdGU+PG5zMTpHZW5kZXI+MTwvbnMxOkdlbmRlcj48L25zMTpQZXJzb24+PG5zMTpFZHVEb2N1bWVudD48bnMxOkVkdUxldmVsPjE8L25zMTpFZHVMZXZlbD48bnMxOlNlcmllcz4xMjM8L25zMTpTZXJpZXM+PG5zMTpOdW1iZXI+NDQ0IDU1NTwvbnMxOk51bWJlcj48bnMxOk9HUk4+MTEyMzQ1Njc4OTU1NTwvbnMxOk9HUk4+PG5zMTpRdWFsaWZpY2F0aW9uPtCY0L3QttC10L3QtdGAPC9uczE6UXVhbGlmaWNhdGlvbj48bnMxOkRvY3VtZW50SXNzdWVEYXRlPjIwMjAtMDYtMDE8L25zMTpEb2N1bWVudElzc3VlRGF0ZT48L25zMTpFZHVEb2N1bWVudD48L25zMTpDaGVja0RvY3VtZW50UmVxdWVzdD48L25zMTpGUkRPUmVxdWVzdD4=","Error":"","Status":4,"Ataches":null,"IsTestMessage":false,"IsReject":false,"StatusResponseId":null}`,
		},
	}
	var ch *rabbitMQ.RabbitChan
	q := NewQueueWorker(ch)
	for _, testingBody := range tests {
		msg := []byte(testingBody.In)
		_, err := q.ParseMessage(msg)
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
func TestQueueWorker_SetMessage(t *testing.T) {
	type test struct {
		Name string
		In   struct {
			model model.SMEVMessage
			body  string
		}
		Expect string
	}
	var statusResponseId uint = 1
	tests := []test{
		test{
			Name:   "OK",
			Expect: `{"Id":0,"PackageId":12,"SmevId":0,"MessageId":"7f91173e-6ace-11eb-bad5-005056ae3841","KindActionId":10,"Type":"response","Xml":"PG5zMTpGUkRPUmVxdWVzdCB4bWxuczpuczE9InVybjovL29icm5hZHpvci1nb3N1c2x1Z2ktcnUvZnJkby8xLjAuMSIgeG1sbnM9InVybjovL3gtYXJ0ZWZhY3RzLXNtZXYtZ292LXJ1L3NlcnZpY2VzL21lc3NhZ2UtZXhjaGFuZ2UvdHlwZXMvMS4xIiB4bWxuczpucz0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiIHhtbG5zOm5zMj0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiPjxuczE6Q2hlY2tEb2N1bWVudFJlcXVlc3Q+PG5zMTpQZXJzb24+PG5zMTpSRkNpdGl6ZW4+dHJ1ZTwvbnMxOlJGQ2l0aXplbj48bnMxOlN1cmVuYW1lPtCf0LXRgtGA0L7QsjwvbnMxOlN1cmVuYW1lPjxuczE6TmFtZT7Qn9C10YLRgDwvbnMxOk5hbWU+PG5zMTpQYXRyb255bWljPtCf0LXRgtGA0L7QstC40Yc8L25zMTpQYXRyb255bWljPjxuczE6QmlydGhkYXRlPjE5OTgtMDItMDM8L25zMTpCaXJ0aGRhdGU+PG5zMTpHZW5kZXI+MTwvbnMxOkdlbmRlcj48L25zMTpQZXJzb24+PG5zMTpFZHVEb2N1bWVudD48bnMxOkVkdUxldmVsPjE8L25zMTpFZHVMZXZlbD48bnMxOlNlcmllcz4xMjM8L25zMTpTZXJpZXM+PG5zMTpOdW1iZXI+NDQ0IDU1NTwvbnMxOk51bWJlcj48bnMxOk9HUk4+MTEyMzQ1Njc4OTU1NTwvbnMxOk9HUk4+PG5zMTpRdWFsaWZpY2F0aW9uPtCY0L3QttC10L3QtdGAPC9uczE6UXVhbGlmaWNhdGlvbj48bnMxOkRvY3VtZW50SXNzdWVEYXRlPjIwMjAtMDYtMDE8L25zMTpEb2N1bWVudElzc3VlRGF0ZT48L25zMTpFZHVEb2N1bWVudD48L25zMTpDaGVja0RvY3VtZW50UmVxdWVzdD48L25zMTpGUkRPUmVxdWVzdD4=","Error":null,"Status":4,"Ataches":null,"IsTestMessage":false,"IsReject":false,"StatusResponseId":1}`,
			In: struct {
				model model.SMEVMessage
				body  string
			}{model: model.SMEVMessage{
				PackageId:        12,
				SmevId:           0,
				MessageId:        "7f91173e-6ace-11eb-bad5-005056ae3841",
				KindActionId:     10,
				Error:            nil,
				Status:           4,
				Ataches:          nil,
				IsTestMessage:    false,
				IsReject:         false,
				StatusResponseId: &statusResponseId,
			}, body: "PG5zMTpGUkRPUmVxdWVzdCB4bWxuczpuczE9InVybjovL29icm5hZHpvci1nb3N1c2x1Z2ktcnUvZnJkby8xLjAuMSIgeG1sbnM9InVybjovL3gtYXJ0ZWZhY3RzLXNtZXYtZ292LXJ1L3NlcnZpY2VzL21lc3NhZ2UtZXhjaGFuZ2UvdHlwZXMvMS4xIiB4bWxuczpucz0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiIHhtbG5zOm5zMj0idXJuOi8veC1hcnRlZmFjdHMtc21ldi1nb3YtcnUvc2VydmljZXMvbWVzc2FnZS1leGNoYW5nZS90eXBlcy9iYXNpYy8xLjEiPjxuczE6Q2hlY2tEb2N1bWVudFJlcXVlc3Q+PG5zMTpQZXJzb24+PG5zMTpSRkNpdGl6ZW4+dHJ1ZTwvbnMxOlJGQ2l0aXplbj48bnMxOlN1cmVuYW1lPtCf0LXRgtGA0L7QsjwvbnMxOlN1cmVuYW1lPjxuczE6TmFtZT7Qn9C10YLRgDwvbnMxOk5hbWU+PG5zMTpQYXRyb255bWljPtCf0LXRgtGA0L7QstC40Yc8L25zMTpQYXRyb255bWljPjxuczE6QmlydGhkYXRlPjE5OTgtMDItMDM8L25zMTpCaXJ0aGRhdGU+PG5zMTpHZW5kZXI+MTwvbnMxOkdlbmRlcj48L25zMTpQZXJzb24+PG5zMTpFZHVEb2N1bWVudD48bnMxOkVkdUxldmVsPjE8L25zMTpFZHVMZXZlbD48bnMxOlNlcmllcz4xMjM8L25zMTpTZXJpZXM+PG5zMTpOdW1iZXI+NDQ0IDU1NTwvbnMxOk51bWJlcj48bnMxOk9HUk4+MTEyMzQ1Njc4OTU1NTwvbnMxOk9HUk4+PG5zMTpRdWFsaWZpY2F0aW9uPtCY0L3QttC10L3QtdGAPC9uczE6UXVhbGlmaWNhdGlvbj48bnMxOkRvY3VtZW50SXNzdWVEYXRlPjIwMjAtMDYtMDE8L25zMTpEb2N1bWVudElzc3VlRGF0ZT48L25zMTpFZHVEb2N1bWVudD48L25zMTpDaGVja0RvY3VtZW50UmVxdWVzdD48L25zMTpGUkRPUmVxdWVzdD4="},
		},
	}
	var ch *rabbitMQ.RabbitChan
	q := NewQueueWorker(ch)
	for _, testingBody := range tests {
		res, err := q.SetMessage(testingBody.In.model, testingBody.In.body)
		if testingBody.Name == `OK` {
			if err != nil {
				t.Errorf(`error test return success: %s`, err.Error())
			}
			if testingBody.Expect != `` {
				assert.Equal(t, testingBody.Expect, string(res), `not a valid expected params`)
			}
		}
		if testingBody.Name == `ERROR` {
			if err == nil {
				t.Errorf(`error test return error nil`)
			}
		}
	}

}
