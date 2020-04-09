package handlers

import (
	"github.com/jinzhu/gorm"
	"persons/config"
	"persons/service"
	"strings"
)

type RowsCls struct {
	Id   uint
	Name string
}

var ListClsTableName = []string{
	`benefits`,
	`campaign_statuses`,
	`campaign_types`,
	`directions`,
	`document_sys_category`,
	`directions`,
	`education_levels`,
	`education_forms`,
}

func (result *ResultCls) GetClsResponse(clsName string) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var db *gorm.DB
	var r []RowsCls
	if service.SearchStringInSliceString(clsName, ListClsTableName) >= 0 {
		db = conn.Table(`cls.` + clsName)
	} else {
		message := `Неизвестный справочник.`
		result.Message = &message
		return
	}
	if result.Search != `` {
		db = db.Where(`UPPER(name) like ?`, `%`+strings.ToUpper(result.Search)+`%`)
	}
	db.Select("name,id").Scan(&r)
	result.Items = r

	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Формы образования не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}

	result.Done = true
	return

}
