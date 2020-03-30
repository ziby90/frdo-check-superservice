package handlers

import (
	"persons/config"
	"persons/digest"
)

type EducationFormRespons struct {
	Id uint 		`json:"id"`
	Name string 	`json:"text"`
}

type EducationLevelRespons struct {
	Id uint 		`json:"id"`
	Name string 	`json:"text"`
}


func (EducationFormRespons) TableName() string {
	return "cls.education_forms"
}
func (EducationLevelRespons) TableName() string {
	return "cls.education_levels"
}

func GetEducFormResponse(Id uint) EducationFormRespons{
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var educForm digest.EducationForm
	db := conn.Find(&educForm, Id)
	if db.Error!=nil {
		return EducationFormRespons{}
	}
	return EducationFormRespons{
		Id:    	educForm.Id,
		Name: 	educForm.Name,
	}
}

func GetEducLevelResponse(Id uint) EducationLevelRespons{
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var educLevel digest.EducationLevel
	db := conn.Find(&educLevel, Id)
	if db.Error!=nil {
		return EducationLevelRespons{}
	}
	return EducationLevelRespons{
		Id:    	educLevel.Id,
		Name: 	educLevel.Name,
	}
}

func (result *ResultCls) GetEducFormResponse(){
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var educForms []EducationFormRespons
	db := conn.Find(&educForms)
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
	if db.RowsAffected>0 {
		result.Done = true
		result.Items = educForms
		return
	}
}
func (result *ResultCls) GetEducLevelResponse(){
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var educLevels []EducationLevelRespons
	db := conn.Find(&educLevels)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Уровни образования не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected>0 {
		result.Done = true
		result.Items = educLevels
		return
	}
}