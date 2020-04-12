package handlers

import (
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
)

func (result *Result) GetListAdmissionVolume(IdCampaign uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var admissions []digest.AdmissionVolume
	db := conn.Where(`id_organization=? AND id_campaign=?`, result.User.CurrentOrganization.Id, IdCampaign).Order(result.Sort.Field + ` ` + result.Sort.Order)
	for _, search := range result.Search {
		if service.SearchStringInSliceString(search[0], campaignSearchArray) >= 0 {
			db = db.Where(`UPPER(`+search[0]+`) LIKE ?`, `%`+strings.ToUpper(search[1])+`%`)
		}
	}
	dbCount := db.Model(&admissions).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Find(&admissions)
	var responses []interface{}
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `КЦП не найдены.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		for _, admission := range admissions {
			db = conn.Model(&admission).Related(&admission.Direction, `IdDirection`)
			db = conn.Model(&admission).Related(&admission.EducationLevel, `IdEducationLevel`)
			//db = conn.Model(&admission).Related(&admission.CampaignStatus, `IdCampaignStatus`)

			responses = append(responses, map[string]interface{}{
				`id`:                     admission.Id,
				`id_campaign`:            admission.IdCampaign,
				`id_organization`:        admission.IdOrganization,
				`id_author`:              admission.IdAuthor,
				`id_campaign_educ_level`: admission.IdCampaignEducLevel,
				`id_direction`:           admission.Direction.Id,
				`name_direction`:         admission.Direction.Name,
				`id_education_level`:     admission.EducationLevel.Id,
				`name_education_level`:   admission.EducationLevel.Name,
				`budget_o`:               admission.BudgetO,
				`budget_oz`:              admission.BudgetOz,
				`budget_z`:               admission.BudgetZ,
				`quota_o`:                admission.QuotaO,
				`quota_oz`:               admission.QuotaOz,
				`quota_z`:                admission.QuotaZ,
				`paid_o`:                 admission.PaidO,
				`paid_oz`:                admission.PaidOz,
				`paid_z`:                 admission.PaidZ,
				`target_o`:               admission.TargetO,
				`target_oz`:              admission.TargetOz,
				`target_z`:               admission.TargetZ,
				`created`:                admission.Created,
			})
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `КЦП не найдены.`
		result.Message = &message
		result.Items = []digest.AdmissionVolume{}
		return
	}
}

func (result *ResultInfo) GetInfoAdmission(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var admission digest.AdmissionVolume
	db := conn.Find(&admission, ID)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `Компания не найдена.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД. `
		result.Message = &message
		return
	}
	if db.RowsAffected > 0 {
		db = conn.Model(&admission).Related(&admission.Direction, `IdDirection`)
		db = conn.Model(&admission).Related(&admission.EducationLevel, `IdEducationLevel`)
		//db = conn.Model(&admission).Related(&admission.CampaignStatus, `IdCampaignStatus`)

		c := map[string]interface{}{
			`id`:                     admission.Id,
			`id_campaign`:            admission.IdCampaign,
			`id_organization`:        admission.IdOrganization,
			`id_author`:              admission.IdAuthor,
			`id_campaign_educ_level`: admission.IdCampaignEducLevel,
			`id_direction`:           admission.Direction.Id,
			`name_direction`:         admission.Direction.Name,
			`id_education_level`:     admission.EducationLevel.Id,
			`name_education_level`:   admission.EducationLevel.Name,
			`budget_o`:               admission.BudgetO,
			`budget_oz`:              admission.BudgetOz,
			`budget_z`:               admission.BudgetZ,
			`quota_o`:                admission.QuotaO,
			`quota_oz`:               admission.QuotaOz,
			`quota_z`:                admission.QuotaZ,
			`paid_o`:                 admission.PaidO,
			`paid_oz`:                admission.PaidOz,
			`paid_z`:                 admission.PaidZ,
			`target_o`:               admission.TargetO,
			`target_oz`:              admission.TargetOz,
			`target_z`:               admission.TargetZ,
			`created`:                admission.Created,
		}
		result.Done = true
		result.Items = c
		return
	} else {
		result.Done = true
		message := `КЦП не найдена.`
		result.Message = &message
		result.Items = []digest.AdmissionVolume{}
		return
	}
}
