package handlers

import (
	"fmt"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
)

type GroupsSpecialty struct {
	Id    uint          `json:"id"`
	Code  string        `json:"code"`
	Name  string        `json:"name"`
	Count int           `json:"count"`
	Items []interface{} `json:"items"`
}

type GroupsEducations struct {
	Id    uint                       `json:"id"`
	Name  string                     `json:"name"`
	Code  string                     `json:"code"`
	Count int                        `json:"count"`
	Items map[string]GroupsSpecialty `json:"groups_specialty"`
}

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

	responses := make(map[string]GroupsEducations)

	if db.RowsAffected > 0 {
		for _, admission := range admissions {
			var groupsEducations GroupsEducations
			var groupsSpecialty GroupsSpecialty
			db = conn.Model(&admission).Related(&admission.Direction, `IdDirection`)
			db = conn.Model(&admission).Related(&admission.EducationLevel, `IdEducationLevel`)

			if val, ok := responses[admission.EducationLevel.Code]; ok {
				groupsEducations = val
			} else {
				groupsEducations.Name = admission.EducationLevel.Name
				groupsEducations.Code = admission.EducationLevel.Code
				groupsEducations.Id = admission.EducationLevel.Id
				groupsEducations.Items = make(map[string]GroupsSpecialty)
			}
			if val, ok := groupsEducations.Items[admission.CodeGroups]; ok {
				groupsSpecialty = val
			} else {
				groupsSpecialty.Name = admission.NameGroups
				groupsSpecialty.Code = admission.CodeGroups
				groupsSpecialty.Id = admission.IdDirection
			}

			admVolume := map[string]interface{}{
				`id`:                   admission.Id,
				`id_campaign`:          admission.IdCampaign,
				`id_organization`:      admission.IdOrganization,
				`id_author`:            admission.IdAuthor,
				`id_direction`:         admission.Direction.Id,
				`name_direction`:       admission.Direction.Name,
				`id_education_level`:   admission.EducationLevel.Id,
				`name_education_level`: admission.EducationLevel.Name,
				`budget_o`:             admission.BudgetO,
				`budget_oz`:            admission.BudgetOz,
				`budget_z`:             admission.BudgetZ,
				`quota_o`:              admission.QuotaO,
				`quota_oz`:             admission.QuotaOz,
				`quota_z`:              admission.QuotaZ,
				`paid_o`:               admission.PaidO,
				`paid_oz`:              admission.PaidOz,
				`paid_z`:               admission.PaidZ,
				`target_o`:             admission.TargetO,
				`target_oz`:            admission.TargetOz,
				`target_z`:             admission.TargetZ,
				`created`:              admission.Created,
			}
			var distributedAdmissionVolumes []digest.DistributedAdmissionVolume
			db = conn.Model(&admission).Related(&distributedAdmissionVolumes)
			if len(distributedAdmissionVolumes) > 0 {
				var disturbes []interface{}
				for _, distributedAdmissionVolume := range distributedAdmissionVolumes {
					db = conn.Model(&distributedAdmissionVolume).Related(&distributedAdmissionVolume.LevelBudget, `IdLevelBudget`)
					ditrVolume := map[string]interface{}{
						`id`:                distributedAdmissionVolume.Id,
						`id_organization`:   distributedAdmissionVolume.IdOrganization,
						`id_author`:         distributedAdmissionVolume.IdAuthor,
						`id_level_budget`:   distributedAdmissionVolume.LevelBudget.Id,
						`name_level_budget`: distributedAdmissionVolume.LevelBudget.Name,
						`budget_o`:          admission.BudgetO,
						`budget_oz`:         admission.BudgetOz,
						`budget_z`:          admission.BudgetZ,
						`quota_o`:           admission.QuotaO,
						`quota_oz`:          admission.QuotaOz,
						`quota_z`:           admission.QuotaZ,
						`paid_o`:            admission.PaidO,
						`paid_oz`:           admission.PaidOz,
						`paid_z`:            admission.PaidZ,
						`target_o`:          admission.TargetO,
						`target_oz`:         admission.TargetOz,
						`target_z`:          admission.TargetZ,
						`created`:           admission.Created,
					}
					disturbes = append(disturbes, ditrVolume)
				}
				admVolume[`distributes`] = disturbes
			}
			groupsSpecialty.Items = append(groupsSpecialty.Items, admVolume)
			groupsSpecialty.Count = len(groupsSpecialty.Items)
			groupsEducations.Items[admission.CodeGroups] = groupsSpecialty
			groupsEducations.Count = len(groupsEducations.Items)
			responses[admission.EducationLevel.Code] = groupsEducations

			fmt.Println(groupsSpecialty.Count)
			fmt.Println(groupsEducations.Count)
		}
		result.Done = true
		result.Items = responses
		return
	} else {
		result.Done = true
		message := `КЦП не найдены.`
		result.Message = &message
		result.Items = []GroupsEducations{}
		return
	}
}
