package handlers

import (
	"github.com/jinzhu/gorm"
	"persons/config"
	"persons/service"
	"strings"
)

type RowsCls struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type SysRows struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	NameTable string `json:"name_table"`
}

var ListClsTableName = []string{
	`achievement_categories`,
	`appeal_statuses`,
	`benefits`,
	`campaign_statuses`,
	`campaign_types`,
	`composition_themes`,
	`compatriot_categories`,
	`directions`,
	`disability_types`,
	`document_categories`,
	`document_sys_category`,
	`document_types`,
	`education_forms`,
	`education_levels`,
	`education_sources`,
	`entrance_test_types`,
	`gender`,
	`level_budget`,
	`military_categories`,
	`olympic_diploma_types`,
	`olympic_levels`,
	`okcm`,
	`orphan_categories`,
	`parents_lost_categories`,
	`radiation_work_categories`,
	`regions`,
	`subjects`,
	`veteran_categories`,
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

func (result *ResultCls) GetClsSysCategoryResponse() {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var r []SysRows

	db := conn.Table(`cls.document_sys_categories`).Select("id, name, name_table").Scan(&r)
	result.Items = r

	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Справочник не найден.`
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