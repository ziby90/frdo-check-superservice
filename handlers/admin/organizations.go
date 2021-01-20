package handlers_admin

import (
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
	"time"
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
				`created_at`:  organization.CreatedAt,
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
func (result *Result) GetListRequestsLinks() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var requests []digest.RequestLinks
	db := conn.Order(result.Sort.Field + ` ` + result.Sort.Order)
	//for _, search := range result.Search {
	//	db = db.Where(`UPPER(`+search[0]+`) LIKE ?`, `%`+strings.ToUpper(search[1])+`%`)
	//}
	dbCount := db.Model(&requests).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Find(&requests)
	var requestLinks []interface{}
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Заявки не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		for _, request := range requests {
			db = conn.Model(&request).Related(&request.Organization, `IdOrganization`)
			db = conn.Model(&request).Related(&request.User, `IdUser`)
			//db = conn.Model(&campaign).Related(&campaign.CampaignStatus, `IdCampaignStatus`)
			var doc interface{}
			file := false
			if request.ConfirmingDoc != nil {
				doc = map[string]interface{}{
					`title`: request.ConfirmingDoc,
				}
				file = true
			} else {
				doc = nil
			}
			c := map[string]interface{}{
				`id`:          request.Id,
				`short_title`: request.Organization.ShortTitle,
				`created_at`:  request.CreatedAt,
				`is_oovo`:     request.Organization.IsOOVO,
				`ogrn`:        request.Organization.Ogrn,
				`kpp`:         request.Organization.Kpp,
				`actual`:      request.Organization.Actual,
				`login`:       request.User.Login,
				`name`:        request.User.Name,
				`surname`:     request.User.Surname,
				`patronymic`:  request.User.Patronymic,
				"file":        file,
				"doc":         doc,
			}
			requestLinks = append(requestLinks, c)
		}
		result.Done = true
		result.Items = requestLinks
		return
	} else {
		result.Done = true
		message := `Заявки не найдены.`
		result.Message = &message
		result.Items = []digest.RequestLinks{}
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
				`created_at`:  organization.CreatedAt,
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
func (result *ResultInfo) GetFileRequestLink(ID uint) {
	result.Done = false
	conn := &config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var doc digest.RequestLinks
	db := conn.Where(`id=?`, ID).Find(&doc)
	if db.Error != nil {
		if db.Error.Error() == "record not found" {
			result.Done = false
			message := "Документ не найден."
			result.Message = &message
			result.Items = []interface{}{}
			return
		}
		message := "Ошибка подключения к БД."
		result.Message = &message
		return
	}
	if doc.ConfirmingDoc != nil && *doc.ConfirmingDoc != `` {
		filename := *doc.ConfirmingDoc
		path := getPath(doc.IdUser, doc.TableName(), doc.CreatedAt) + filename
		result.Items = path
	} else {
		message := "Файл не найден."
		result.Message = &message
		return
	}
	result.Done = true
	return
}
func (result *ResultInfo) AcceptRequestLink(idRequest uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var item digest.RequestLinks
	db := conn.Where(`id_status!=1`).Find(&item, idRequest)
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Заявка не найдена.`
			result.Message = &message
			return
		}
		result.SetErrorResult(`Ошибка подключения к БД. `)
		return
	}
	tx := db.Begin()
	if db.RowsAffected > 0 {
		t := time.Now()
		item.UpdatedAt = &t
		item.IdStatus = 2
		item.IdAuthor = &result.User.Id
		db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&item)
		if db.Error != nil {
			result.SetErrorResult(db.Error.Error())
			tx.Rollback()
			return
		}
		link := digest.OrganizationsUsers{
			IdUser:         item.IdUser,
			IdOrganization: item.IdOrganization,
			IdLink:         item.Id,
		}
		db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&link)
		if db.Error != nil {
			result.SetErrorResult(db.Error.Error())
			tx.Rollback()
			return
		}
		result.Done = true
		result.Items = map[string]interface{}{
			`id_request`:      item.Id,
			`id_link`:         link.Id,
			`id_user`:         item.IdUser,
			`id_organization`: item.IdOrganization,
		}
		tx.Commit()
		return
	} else {
		result.Done = true
		message := `Заявка не найдена.`
		result.Message = &message
		result.Items = []digest.User{}
		return
	}
}
func (result *ResultInfo) DeclineRequestLink(idRequest uint, comment string) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var item digest.RequestLinks
	db := conn.Where(`id_status!=1`).Find(&item, idRequest)
	if db.Error != nil {
		if db.Error.Error() == service.ErrorDbNotFound {
			result.Done = true
			message := `Заявка не найдена.`
			result.Message = &message
			return
		}
		result.SetErrorResult(`Ошибка подключения к БД. `)
		return
	}
	tx := db.Begin()
	if db.RowsAffected > 0 {
		t := time.Now()
		item.UpdatedAt = &t
		item.IdStatus = 3
		item.Comment = &comment
		item.IdAuthor = &result.User.Id
		db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&item)
		if db.Error != nil {
			result.SetErrorResult(db.Error.Error())
			tx.Rollback()
			return
		}
		result.Done = true
		result.Items = map[string]interface{}{
			`id_request`:      item.Id,
			`id_user`:         item.IdUser,
			`id_organization`: item.IdOrganization,
		}
		tx.Commit()
		return
	} else {
		result.Done = true
		message := `Заявка не найдена.`
		result.Message = &message
		result.Items = []digest.User{}
		return
	}
}
