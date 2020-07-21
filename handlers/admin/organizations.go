package handlers_admin

import (
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
)

var OrganizationSearchArray = []string{
	`short_title`,
	`ogrn`,
	`kpp`,
}

func (result *Result) GetListOrganization() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var organizations []digest.Organization
	db := conn.Order(result.Sort.Field + ` ` + result.Sort.Order)
	for _, search := range result.Search {
		db = db.Where(`UPPER(`+search[0]+`) LIKE ?`, `%`+strings.ToUpper(search[1])+`%`)
	}
	dbCount := db.Model(&organizations).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Find(&organizations)
	var orgs []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Компании не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		for _, organization := range organizations {
			//db = conn.Model(&campaign).Related(&campaign.CampaignType, `IdCampaignType`)
			//db = conn.Model(&campaign).Related(&campaign.CampaignStatus, `IdCampaignStatus`)
			c := map[string]interface{}{
				`id`:          organization.Id,
				`short_title`: organization.ShortTitle,
				`created`:     organization.Created,
				`is_oovo`:     organization.IsOOVO,
				`ogrn`:        organization.Ogrn,
				`kpp`:         organization.Kpp,
				`actual`:      organization.Actual,
			}
			orgs = append(orgs, c)
		}
		result.Done = true
		result.Items = orgs
		return
	} else {
		result.Done = true
		message := `Организации не найдены.`
		result.Message = &message
		result.Items = []digest.Organization{}
		return
	}
}
func (result *ResultInfo) GetListOrganizationShort() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var organizations []digest.Organization
	db := conn.Where(`actual is true`).Find(&organizations)
	var orgs []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Организации не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		for _, organization := range organizations {
			c := map[string]interface{}{
				`id`:          organization.Id,
				`short_title`: organization.ShortTitle,
				`created`:     organization.Created,
				`is_oovo`:     organization.IsOOVO,
				`ogrn`:        organization.Ogrn,
				`kpp`:         organization.Kpp,
			}
			orgs = append(orgs, c)
		}
		result.Done = true
		result.Items = orgs
		return
	} else {
		result.Done = true
		message := `Организации не найдены.`
		result.Message = &message
		result.Items = []digest.Organization{}
		return
	}
}
