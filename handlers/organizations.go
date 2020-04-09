package handlers

import (
	"database/sql"
	"net/http"
	"persons/config"
	"persons/digest"
	"strconv"
)

func (result *Result) GetListOrganization() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var organizations []digest.Organization
	db := conn.Order(result.Sort.Field + ` ` + result.Sort.Order)
	//if result.Search != `` {
	//	db = db.Where(`UPPER(name) LIKE ?`, `%`+strings.ToUpper(result.Search)+`%`)
	//}
	dbCount := db.Model(&organizations).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Find(&organizations)
	var orgs []interface{}
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
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

func (result *ResultInfo) GetInfoOrganization(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var organization digest.Organization
	db := conn.Find(&organization, ID)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Организация не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		//db = conn.Model(&campaign).Related(&campaign.CampaignType, `IdCampaignType`)
		//db = conn.Model(&campaign).Related(&campaign.CampaignStatus, `IdCampaignStatus`)

		////var educForms []EducationFormRespons
		//var campEducForms []digest.CampaignEducForm
		//db = conn.Where(`id_campaign=?`, campaign.Id).Find(&campEducForms)
		//var campEducLevels []digest.CampaignEducLevel
		//db = conn.Where(`id_campaign=?`, campaign.Id).Find(&campEducLevels)

		c := map[string]interface{}{
			`id`:          organization.Id,
			`short_title`: organization.ShortTitle,
			`created`:     organization.Created,
		}
		result.Done = true
		result.Items = c
		return
	} else {
		result.Done = true
		message := `Организация не найдены.`
		result.Message = &message
		result.Items = make(map[string]string)
		return
	}
}

func CheckOrgCookie(user digest.User, r *http.Request) uint {
	currentOrgId := uint(0)
	cookieOrg, err := r.Cookie(`current-org`)
	if err == nil {
		u64, err := strconv.ParseUint(cookieOrg.Value, 10, 32)
		if err == nil {
			currentOrgId = uint(u64)
		}
	}
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var row *sql.Row
	if user.Role.Code == `administrator` {
		row = conn.Table(`admin.organizations`).Where(`id=?`, currentOrgId).Select(`id`).Row()
	} else {
		row = conn.Table(`admin.organizations_users`).Where(`id_user=? AND id_organization=?`, user.Id, currentOrgId).Select(`id_organization`).Row()
	}

	if err == nil {
		var organizationId uint
		err = row.Scan(&organizationId)
		if organizationId > 0 {
			return organizationId
		}
	}
	return 0

}
