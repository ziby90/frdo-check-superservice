package error_handler

import (
	"testing"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
)

type i interface{}

type strErr string

func (e strErr) Error() string {
	return string(e)

}

var standartSSErr = ErrorType.New(ErrorType{Type: 12}, "Тестовая ошибка 1")
var standartErr strErr

func TestParseErrorToJsonB(t *testing.T) {
	assert.New(t)
	JsonB := ParseErrorToJsonB(standartSSErr)
	var typeJ postgres.Jsonb
	assert.IsType(t, typeJ, JsonB)
	assert.EqualValues(t, `{"code":12,"originalError":"Тестовая ошибка 1","context":{}}`, string(JsonB.RawMessage), "Возвращено значение не по умолчанию")
	standartErr = "некая Тестовая ошибка 2"
	JsonB = ParseErrorToJsonB(standartErr)
	assert.EqualValues(t, `{"code":0,"originalError":"некая Тестовая ошибка 2","context":null}`, string(JsonB.RawMessage), "Возвращено значение не по умолчанию")

}

func TestGetType(t *testing.T) {
	t.Parallel()
	assert.New(t)
	typeErr := GetType(standartSSErr)
	var typeJ ErrorType
	assert.IsType(t, typeJ, typeErr, "Возвращен не тип ошибки суперсервиса")
	typeErr = GetType(standartErr)
	assert.EqualValues(t, 0, typeErr, "Возвращено значение не по умолчанию")
}
