package handlers

import (
	"fmt"
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
)

type GroupsSpecialty struct {
	Id      uint          `json:"id"`
	Code    string        `json:"code"`
	Name    string        `json:"name"`
	Count   int           `json:"count"`
	Items   []interface{} `json:"items, omitempty"`
	IdItems []int         `json:"id_items, omitempty"`
}

type SumAdmVolume struct {
	Code             string `json:"code"`
	IdCampaign       uint   `json:"id_campaign"`
	IdEducationLevel uint   `json:"id_education_level"`
	BudgetO          int64  `json:"budget_o,omitempty"`
	BudgetOz         int64  `json:"budget_oz,omitempty"`
	BudgetZ          int64  `json:"budget_z,omitempty"`
	QuotaO           int64  `json:"quota_o,omitempty"`
	QuotaOz          int64  `json:"quota_oz,omitempty"`
	QuotaZ           int64  `json:"quota_z,omitempty"`
	PaidO            int64  `json:"paid_o,omitempty"`
	PaidOz           int64  `json:"paid_oz,omitempty"`
	PaidZ            int64  `json:"paid_z,omitempty"`
	TargetO          int64  `json:"target_o,omitempty"`
	TargetOz         int64  `json:"target_oz,omitempty"`
	TargetZ          int64  `json:"target_z,omitempty"`
}

type GroupSpecialty struct {
	Id                 uint   `json:"id_specialty"`
	IdCampaign         uint   `json:"id_campaign"`
	CodeSpecialty      string `json:"code_specialty"`
	NameSpecialty      string `json:"name_specialty"`
	IdEducationLevel   uint   `json:"id_education_level"`
	NameEducationLevel string `json:"name_education_level"`
	Count              int    `json:"count"`
	IdItems            []int  `json:"id_items"`
	BudgetO            int64  `json:"budget_o"`
	BudgetOz           int64  `json:"budget_oz"`
	BudgetZ            int64  `json:"budget_z"`
	QuotaO             int64  `json:"quota_o"`
	QuotaOz            int64  `json:"quota_oz"`
	QuotaZ             int64  `json:"quota_z"`
	PaidO              int64  `json:"paid_o"`
	PaidOz             int64  `json:"paid_oz"`
	PaidZ              int64  `json:"paid_z"`
	TargetO            int64  `json:"target_o"`
	TargetOz           int64  `json:"target_oz"`
	TargetZ            int64  `json:"target_z"`
}

type GroupsEducations struct {
	Id             uint                       `json:"id"`
	Name           string                     `json:"name"`
	Code           string                     `json:"code"`
	CountSpecialty int                        `json:"count_specialty"`
	CountItems     int                        `json:"count_items"`
	Items          map[string]GroupsSpecialty `json:"groups_specialty"`
}

func (result *Result) GetListAdmissionVolume(IdCampaign uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var admissions []digest.AdmissionVolume
	db := conn.Where(`id_organization=? AND id_campaign=?`, result.User.CurrentOrganization.Id, IdCampaign).Order(`id, id_direction, id_education_level`)
	for _, search := range result.Search {
		if service.SearchStringInSliceString(search[0], CampaignSearchArray) >= 0 {
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
				`code_direction`:       admission.Direction.Code,
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
				//var disturbes []interface{}
				//for _, distributedAdmissionVolume := range distributedAdmissionVolumes {
				//	db = conn.Model(&distributedAdmissionVolume).Related(&distributedAdmissionVolume.LevelBudget, `IdLevelBudget`)
				//	ditrVolume := map[string]interface{}{
				//		`id`:                distributedAdmissionVolume.Id,
				//		`id_organization`:   distributedAdmissionVolume.IdOrganization,
				//		`id_author`:         distributedAdmissionVolume.IdAuthor,
				//		`id_level_budget`:   distributedAdmissionVolume.LevelBudget.Id,
				//		`name_level_budget`: distributedAdmissionVolume.LevelBudget.Name,
				//		`budget_o`:          admission.BudgetO,
				//		`budget_oz`:         admission.BudgetOz,
				//		`budget_z`:          admission.BudgetZ,
				//		`quota_o`:           admission.QuotaO,
				//		`quota_oz`:          admission.QuotaOz,
				//		`quota_z`:           admission.QuotaZ,
				//		`paid_o`:            admission.PaidO,
				//		`paid_oz`:           admission.PaidOz,
				//		`paid_z`:            admission.PaidZ,
				//		`target_o`:          admission.TargetO,
				//		`target_oz`:         admission.TargetOz,
				//		`target_z`:          admission.TargetZ,
				//		`created`:           admission.Created,
				//	}
				//	disturbes = append(disturbes, ditrVolume)
				//}
				admVolume[`distributes`] = true
			} else {
				admVolume[`distributes`] = false
			}
			groupsSpecialty.Items = append(groupsSpecialty.Items, admVolume)
			groupsSpecialty.Count = len(groupsSpecialty.Items)
			groupsEducations.Items[admission.CodeGroups] = groupsSpecialty
			groupsEducations.CountSpecialty = len(groupsEducations.Items)
			groupsEducations.CountItems += groupsSpecialty.Count

			responses[admission.EducationLevel.Code] = groupsEducations

		}
		var arrayResponses []interface{}
		for i, _ := range responses {
			var arraySpecialty []interface{}
			for s, _ := range responses[i].Items {
				arraySpecialty = append(arraySpecialty, responses[i].Items[s])
			}

			arrayResponses = append(arrayResponses, map[string]interface{}{
				"id":               responses[i].Id,
				"name":             responses[i].Name,
				"code":             responses[i].Code,
				"count_specialty":  responses[i].CountSpecialty,
				"count_items":      responses[i].CountItems,
				"groups_specialty": arraySpecialty,
			})
		}
		result.Done = true
		result.Items = arrayResponses
		return
	} else {
		result.Done = true
		message := `КЦП не найдены.`
		result.Message = &message
		result.Items = []GroupsEducations{}
		return
	}
}

func (result *Result) GetListAdmissionVolumeBySpec(IdCampaign uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var admissions []digest.AdmissionVolume
	db := conn.Where(`id_organization=? AND id_campaign=?`, result.User.CurrentOrganization.Id, IdCampaign).Order(`id_direction`)
	for _, search := range result.Search {
		if service.SearchStringInSliceString(search[0], CampaignSearchArray) >= 0 {
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

	var responses []interface{}
	specialty := make(map[uint]GroupSpecialty)
	var groups []interface{}
	if db.RowsAffected > 0 {
		for _, admission := range admissions {
			fmt.Println(admission.Id)
			var admVolume map[string]interface{}
			var groupSpecialty GroupSpecialty
			//db = conn.Model(&admission).Related(&admission.EducationLevel, `IdEducationLevel`)
			if val, ok := specialty[admission.IdGroups]; ok {
				groupSpecialty = val
			} else {
				groupSpecialty.Id = admission.IdGroups
			}
			admVolume = map[string]interface{}{
				`id`:                   admission.Id,
				`id_campaign`:          admission.IdCampaign,
				`id_education_level`:   admission.IdEducationLevel,
				`name_education_level`: admission.NameEducationLevel,
				`distributed`:          admission.Distributed,
				`id_specialty`:         admission.IdDirection,
				`code_specialty`:       admission.CodeSpecialty,
				`name_specialty`:       admission.NameSpecialty,
				`id_groups`:            admission.IdGroups,
				`code_groups`:          admission.CodeGroups,
				`name_groups`:          admission.NameGroups,
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
			}

			groupSpecialty.Count++
			groupSpecialty.IdCampaign = admission.IdCampaign
			groupSpecialty.IdEducationLevel = admission.IdEducationLevel
			groupSpecialty.NameEducationLevel = admission.NameEducationLevel
			groupSpecialty.CodeSpecialty = admission.CodeGroups
			groupSpecialty.NameSpecialty = admission.NameGroups
			responses = append(responses, admVolume)
			groupSpecialty.IdItems = append(groupSpecialty.IdItems, len(responses)-1)
			specialty[admission.IdGroups] = groupSpecialty
		}
		for _, value := range specialty {
			var sum SumAdmVolume
			conn.Raw(`SELECT
							eg.code,
							av.id_campaign,
							av.id_education_level,
							SUM(av.budget_o) AS budget_o,
							SUM(av.budget_oz) AS budget_oz,
							SUM(av.budget_z) AS budget_z,
							SUM(av.paid_o) AS paid_o,
							SUM(av.paid_oz)AS paid_oz,
							SUM(av.paid_z)AS paid_z,
							SUM(av.target_o)AS target_o,
							SUM(av.target_oz)AS target_oz,
							SUM(av.target_z)AS target_z,
							SUM(av.quota_o)AS quota_o,
							SUM(av.quota_oz)AS quota_oz,
							SUM(av.quota_z)AS quota_z
					   FROM cmp.admission_volume av  
					   JOIN cls.directions d ON d.id = av.id_direction
					   JOIN cls.v_okso_enlarged_group eg ON d.id_parent = eg.id
						Where id_campaign=? and id_education_level=?
						 GROUP BY
						eg.code,
						av.id_campaign,
						av.id_education_level`, value.IdCampaign, value.IdEducationLevel).Scan(&sum)
			value.BudgetO = sum.BudgetO
			value.BudgetOz = sum.BudgetOz
			value.BudgetZ = sum.BudgetZ
			value.QuotaO = sum.QuotaO
			value.QuotaOz = sum.QuotaOz
			value.QuotaZ = sum.QuotaZ
			value.PaidO = sum.PaidO
			value.PaidOz = sum.PaidOz
			value.PaidZ = sum.PaidZ
			value.TargetO = sum.TargetO
			value.TargetOz = sum.TargetOz
			value.TargetZ = sum.TargetZ
			groups = append(groups, value)
		}
		result.Done = true
		result.Items = map[string]interface{}{
			`items`:  responses,
			`groups`: groups,
		}
		return
	} else {
		result.Done = true
		message := `КЦП не найдены.`
		result.Message = &message
		result.Items = make(map[string]interface{})
		return
	}
}
