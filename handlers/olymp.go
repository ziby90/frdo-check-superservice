package handlers

import (
	"fmt"
	"persons/config"
	"persons/digest"
	"persons/service"
)

func (result *Result) GetListOlympics() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var olympics []digest.Olympics
	db := conn.Order(`created desc`).Find(&olympics)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Олимпиады не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		for _, olympic := range olympics {

			c := map[string]interface{}{
				`id`:      olympic.Id,
				`name`:    olympic.Name + fmt.Sprintf(` (%v)`, olympic.OlympYear),
				`created`: olympic.Created,
			}
			responses = append(responses, c)
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `Олимпиады не найдены.`
		result.Message = &message
		result.Items = []digest.Olympics{}
		return
	}
}
