package service

import (
	"frdo-check-superservice/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestXmlParsService_SetResponseAllDocs(t *testing.T) {
	s := getNewServiceForTestSearch()
	type test struct {
		Name string
		In   struct {
			Code uint
			Docs []model.EduDocument
		}
		Expect string
	}

	tests := []test{
		test{
			Name:   "OK",
			Expect: `PG5zMTpGUkRPUmVzcG9uc2UgeG1sbnM6bnMxPSJ1cm46Ly9vYnJuYWR6b3ItZ29zdXNsdWdpLXJ1L2ZyZG8vMS4wLjEiPjxuczE6R2V0RG9jdW1lbnRzUmVzcG9uc2U+PG5zMTpSZXN1bHQ+MTwvbnMxOlJlc3VsdD48bnMxOkRvY3VtZW50cz48bnMxOkRvY3VtZW50PjxFZHVMZXZlbD4xPC9FZHVMZXZlbD48TnVtYmVyPjExMTExMTExMTExMTE8L051bWJlcj48L25zMTpEb2N1bWVudD48bnMxOkRvY3VtZW50PjxFZHVMZXZlbD4yPC9FZHVMZXZlbD48TnVtYmVyPjIyMjIyMjIyMjIyMjI8L051bWJlcj48L25zMTpEb2N1bWVudD48L25zMTpEb2N1bWVudHM+PC9uczE6R2V0RG9jdW1lbnRzUmVzcG9uc2U+PC9uczE6RlJET1Jlc3BvbnNlPg==`,
			In: struct {
				Code uint
				Docs []model.EduDocument
			}{
				Code: 1,
				Docs: []model.EduDocument{
					model.EduDocument{
						EduLevel: 1,
						Number:   "1111111111111",
					},
					model.EduDocument{
						EduLevel: 2,
						Number:   "2222222222222",
					},
				},
			},
		},
	}
	for _, testingBody := range tests {
		result, err := s.SetResponseAllDocs(testingBody.In.Code, testingBody.In.Docs)
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
func TestXmlParsService_SetResponseSearch(t *testing.T) {
	s := getNewServiceForTestSearch()
	type test struct {
		Name string
		In   struct {
			Code    uint
			Message string
		}
		Expect string
	}

	tests := []test{
		test{
			Name:   "OK",
			Expect: `PG5zMTpGUkRPUmVzcG9uc2UgeG1sbnM6bnMxPSJ1cm46Ly9vYnJuYWR6b3ItZ29zdXNsdWdpLXJ1L2ZyZG8vMS4wLjEiPjxuczE6Q2hlY2tEb2N1bWVudFJlc3BvbnNlPjxuczE6Q2hlY2tEb2N1bWVudENvZGU+MTwvbnMxOkNoZWNrRG9jdW1lbnRDb2RlPjxuczE6Q2hlY2tEb2N1bWVudFRleHQ+0LTQvtC60YPQvNC10L3RgiDQv9C+INGA0LXQutCy0LjQt9C40YLQsNC8INC90LUg0L3QsNC50LTQtdC9PC9uczE6Q2hlY2tEb2N1bWVudFRleHQ+PC9uczE6Q2hlY2tEb2N1bWVudFJlc3BvbnNlPjwvbnMxOkZSRE9SZXNwb25zZT4=`,
			In: struct {
				Code    uint
				Message string
			}{
				Code:    1,
				Message: `документ по реквизитам не найден`,
			},
		},
	}
	for _, testingBody := range tests {
		result, err := s.SetResponseSearch(testingBody.In.Code, testingBody.In.Message)
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
