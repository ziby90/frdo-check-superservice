package handlers

import (
	"fmt"
	"persons/config"
	"persons/digest"
	"strings"
	"time"
)

type GroupsSpecialty struct {
	Id      uint          `json:"id"`
	Code    string        `json:"code"`
	Name    string        `json:"name"`
	Count   int           `json:"count"`
	Items   []interface{} `json:"items, omitempty"`
	IdItems []int         `json:"id_items, omitempty"`
}
type EditDistributedAdmissionVolume struct {
	IdAdmissionVolume uint    `json:"id_admission_volume"`
	IdLevelBudget     uint    `json:"id_level_budget"`
	Uid               *string `json:"UID" json:"uid" ` // Идентификатор от организации
	BudgetO           int64   `json:"budget_o"`
	BudgetOz          int64   `json:"budget_oz"`
	BudgetZ           int64   `json:"budget_z"`
	QuotaO            int64   `json:"quota_o"`
	QuotaOz           int64   `json:"quota_oz"`
	QuotaZ            int64   `json:"quota_z"`
	PaidO             int64   `json:"paid_o"`
	PaidOz            int64   `json:"paid_oz"`
	PaidZ             int64   `json:"paid_z"`
	TargetO           int64   `json:"target_o"`
	TargetOz          int64   `json:"target_oz"`
	TargetZ           int64   `json:"target_z"`
}

type SumAdmVolume struct {
	Code             string `json:"code"`
	IdCampaign       uint   `json:"id_campaign"`
	IdEducationLevel uint   `json:"id_education_level"`
	BudgetO          int64  `json:"budget_o"`
	BudgetOz         int64  `json:"budget_oz"`
	BudgetZ          int64  `json:"budget_z"`
	QuotaO           int64  `json:"quota_o"`
	QuotaOz          int64  `json:"quota_oz"`
	QuotaZ           int64  `json:"quota_z"`
	PaidO            int64  `json:"paid_o"`
	PaidOz           int64  `json:"paid_oz"`
	PaidZ            int64  `json:"paid_z"`
	TargetO          int64  `json:"target_o"`
	TargetOz         int64  `json:"target_oz"`
	TargetZ          int64  `json:"target_z"`
	SumDistributed   int64  `json:"sum_distributed"`
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
type AddAdmission struct {
	IdCampaign uint               `json:"id_campaign"`
	Data       []AddAdmissionData `json:"data"`
}
type AddBudget struct {
	LevelsBudget []uint `json:"levels_budget"`
}
type AddAdmissionData struct {
	IdEducationLevel  uint   `json:"id_education_level"`
	IdGroupSpeciality uint   `json:"id_group_speciality"`
	IdSpeciality      []uint `json:"id_speciality"`
}

type GroupsEducations struct {
	Id             uint                       `json:"id"`
	Name           string                     `json:"name"`
	Code           string                     `json:"code"`
	CountSpecialty int                        `json:"count_specialty"`
	CountItems     int                        `json:"count_items"`
	Items          map[string]GroupsSpecialty `json:"groups_specialty"`
}

var AdmissionVolumeSearchArray = []string{
	`code_specialty`,
	`name_specialty`,
	`code_name`,
}

func (result *Result) GetListAdmissionVolume(IdCampaign uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var admissions []digest.AdmissionVolume
	db := conn.Where(`id_organization=? AND id_campaign=?`, result.User.CurrentOrganization.Id, IdCampaign).Order(`id, id_direction, id_education_level`)
	for _, search := range result.Search {
		db = db.Where(`UPPER(`+search[0]+`) LIKE ?`, `%`+strings.ToUpper(search[1])+`%`)
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
		if search[0] == `code_name` {
			db = db.Where(`UPPER( code_specialty || ' ' || name_specialty) LIKE ?`, `%`+strings.ToUpper(search[1])+`%`)
		} else {
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
				`uid`:                  admission.Uid,
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
				`sum_distributed`:      admission.SumDistributed,
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
						Where av.id_campaign=? and av.id_education_level=? AND eg.code=?
						 GROUP BY
						eg.code,
						av.id_campaign,
						av.id_education_level`, value.IdCampaign, value.IdEducationLevel, value.CodeSpecialty).Scan(&sum)
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

func (result *Result) GetAdmissionVolumeById(IdAdmission uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var admission digest.AdmissionVolume
	db := conn.Find(&admission, IdAdmission)
	if db.Error != nil {
		if db.Error.Error() == `record not found` {
			result.Done = true
			message := `КЦП не найден.`
			result.Message = &message
			return
		}
		message := `Ошибка подключения к БД.`
		result.Message = &message
		return
	}

	var response map[string]interface{}
	if admission.Id > 0 {
		response = map[string]interface{}{
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
			`uid`:                  admission.Uid,
		}
		var distrs []interface{}
		var distributed []digest.DistributedAdmissionVolume
		db = conn.Preload(`LevelBudget`).Where(`id_admission_volume=?`, admission.Id).Order(`id_level_budget`).Find(&distributed)
		if len(distributed) > 0 {
			for index, _ := range distributed {
				distr := distributed[index]
				distrs = append(distrs, map[string]interface{}{
					`id`:                   distr.Id,
					`id_admission_volume`:  distr.AdmissionVolumeId,
					`id_level_budget`:      distr.LevelBudget.Id,
					`name_level_budget`:    distr.LevelBudget.Name,
					`id_education_level`:   admission.IdEducationLevel,
					`name_education_level`: admission.NameEducationLevel,
					`id_specialty`:         admission.IdDirection,
					`code_specialty`:       admission.CodeSpecialty,
					`name_specialty`:       admission.NameSpecialty,
					`budget_o`:             distr.BudgetO,
					`budget_oz`:            distr.BudgetOz,
					`budget_z`:             distr.BudgetZ,
					`quota_o`:              distr.QuotaO,
					`quota_oz`:             distr.QuotaOz,
					`quota_z`:              distr.QuotaZ,
					//`paid_o`:               distr.PaidO,
					//`paid_oz`:              distr.PaidOz,
					//`paid_z`:               distr.PaidZ,
					`target_o`:  distr.TargetO,
					`target_oz`: distr.TargetOz,
					`target_z`:  distr.TargetZ,
				})
			}
			response[`distributed`] = distrs
		}

		result.Done = true
		result.Items = response
		return
	} else {
		result.Done = true
		message := `КЦП не найдены.`
		result.Message = &message
		result.Items = make(map[string]interface{})
		return
	}
}

func (result *ResultInfo) AddAdmission(admsData AddAdmission) {
	conn := config.Db.ConnGORM
	var addIdsAdmission []uint
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	fmt.Println(admsData)
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign
	db := tx.Preload(`CampaignType`).Find(&campaign, admsData.IdCampaign)
	if campaign.Id < 1 {
		result.SetErrorResult(`Компания не найдена`)
		tx.Rollback()
		return
	}
	err := CheckAddAdmission(campaign.Id)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	for _, admData := range admsData.Data {
		for _, specialty := range admData.IdSpeciality {
			var admission digest.AdmissionVolume
			admission.Organization.Id = result.User.CurrentOrganization.Id
			admission.IdOrganization = result.User.CurrentOrganization.Id
			admission.IdAuthor = &result.User.Id
			admission.Created = time.Now()
			admission.IdEducationLevel = admData.IdEducationLevel

			row := conn.Table(`cls.edu_levels_campaign_types`).Where(`id_campaign_types=? AND id_education_level=?`, campaign.CampaignType.Id, admData.IdEducationLevel).Select(`id`).Row()
			var idEducLevelCampaignType uint
			err = row.Scan(&idEducLevelCampaignType)
			if err != nil || idEducLevelCampaignType <= 0 {
				result.SetErrorResult(`Данный уровень образования не соответствует типу приемной компании`)
				tx.Rollback()
				return
			}
			admission.IdCampaign = admsData.IdCampaign

			var direction digest.Direction
			db = tx.Where(`id_education_level=? AND id_parent=?`, admData.IdEducationLevel, admData.IdGroupSpeciality).Find(&direction, specialty)
			if direction.Id < 1 {
				result.SetErrorResult(`Специальность не найдена`)
				tx.Rollback()
				return
			}
			var exist digest.AdmissionVolume
			tx.Where(`id_direction=? AND id_campaign=?`, specialty, admsData.IdCampaign).Find(&exist)
			if exist.Id > 0 {
				result.SetErrorResult(`Для специальности "` + direction.Name + `" уже созданы кцп.`)
				tx.Commit()
				return
			}
			admission.IdDirection = specialty
			admission.Actual = true
			db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Table(`cmp.admission_volume`).Create(&admission)
			if db.Error == nil {
				addIdsAdmission = append(addIdsAdmission, admission.Id)
			} else {
				result.SetErrorResult(db.Error.Error())
				tx.Rollback()
				return
			}
		}
	}
	if (len(addIdsAdmission)) > 0 {
		result.Items = addIdsAdmission
		result.Done = true
		tx.Commit()
	} else {
		result.SetErrorResult(`Ошибка добавления`)
		tx.Rollback()
		return
	}

}

func (result *ResultInfo) AddAdmissionBudget(idAdmission uint, data AddBudget) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	var admission digest.AdmissionVolume
	conn.Where(`id=?`, idAdmission).Find(&admission)
	if admission.Id <= 0 {
		result.SetErrorResult(`КЦП не найдены`)
		tx.Rollback()
		return
	}
	err := CheckAddAdmission(admission.IdCampaign)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	var addIdsBudget []uint
	conn.LogMode(config.Conf.Dblog)

	for _, budget := range data.LevelsBudget {
		var count int
		var distributed digest.DistributedAdmissionVolume
		db := conn.Table(`cmp.distributed_admission_volume`).Where(`id_admission_volume=? AND id_level_budget=?`, idAdmission, budget).Count(&count)
		if db.Error != nil {
			result.SetErrorResult(`Ошибка добавления1`)
			tx.Rollback()
			return
		}
		if count > 0 {
			result.SetErrorResult(`Данный уровень бюджета уже существуют`)
			tx.Rollback()
			return
		}

		distributed.IdLevelBudget = budget
		distributed.AdmissionVolumeId = idAdmission
		distributed.IdOrganization = result.User.CurrentOrganization.Id
		distributed.IdAuthor = &result.User.Id
		distributed.Created = time.Now()
		db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&distributed)
		if db.Error == nil {
			addIdsBudget = append(addIdsBudget, distributed.Id)
		} else {
			result.SetErrorResult(db.Error.Error())
			tx.Rollback()
			return
		}
	}
	if (len(addIdsBudget)) > 0 {
		result.Items = addIdsBudget
		result.Done = true
		tx.Commit()
	} else {
		result.Items = data
		result.SetErrorResult(`Ошибка добавления`)
		tx.Rollback()
		return
	}

}

func (result *ResultInfo) EditAdmission(IdAdmission uint, admData digest.AdmissionVolume) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var oldAdmission digest.AdmissionVolume
	db := tx.Preload(`Direction`).Find(&oldAdmission, IdAdmission)
	if oldAdmission.Id == 0 || db.Error != nil {
		result.SetErrorResult(`Кцп не найдены`)
		tx.Rollback()
		return
	}
	err := CheckEditAdmission(oldAdmission.IdCampaign)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	if oldAdmission.IdOrganization != result.User.CurrentOrganization.Id {
		result.SetErrorResult(`Организация не соответствует`)
		tx.Rollback()
		return
	}
	if oldAdmission.IdDirection == admData.IdDirection {
		oldAdmission.IdAuthor = &result.User.Id
		oldAdmission.BudgetO = admData.BudgetO
		oldAdmission.BudgetOz = admData.BudgetOz
		oldAdmission.BudgetZ = admData.BudgetZ
		oldAdmission.QuotaO = admData.QuotaO
		oldAdmission.QuotaOz = admData.QuotaOz
		oldAdmission.QuotaZ = admData.QuotaZ
		oldAdmission.PaidO = admData.PaidO
		oldAdmission.PaidOz = admData.PaidOz
		oldAdmission.PaidZ = admData.PaidZ
		oldAdmission.TargetO = admData.TargetO
		oldAdmission.TargetOz = admData.TargetOz
		oldAdmission.TargetZ = admData.TargetZ
		if admData.Uid != nil {
			var exist digest.AdmissionVolume
			db = tx.Where(`upper(uid)=upper(?) AND id_organization=?`, admData.Uid, result.User.CurrentOrganization.Id).Find(&exist)
			if exist.Id > 0 {
				result.SetErrorResult(`Кцп с данным uid уже существуют`)
				tx.Rollback()
				return
			} else {
				oldAdmission.Uid = admData.Uid
			}
		}

		db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Table(`cmp.admission_volume`).Save(&oldAdmission)
		if db.Error == nil {
			tx.Commit()
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
							SUM(av.quota_z)AS quota_z,
							SUM(vavsg.sum_distributed) AS sum_distributed
					   FROM cmp.admission_volume av  
					   JOIN cls.directions d ON d.id = av.id_direction
					   JOIN cls.v_okso_enlarged_group eg ON d.id_parent = eg.id
						LEFT JOIN cmp.v_admission_volume_specialty_groups vavsg ON vavsg.id = av.id
						Where av.id_campaign=? and av.id_education_level=? AND eg.id=?						
						 GROUP BY
						eg.code,
						av.id_campaign,
						av.id_education_level`, oldAdmission.IdCampaign, oldAdmission.IdEducationLevel, oldAdmission.Direction.IdParent).Scan(&sum)

			result.Done = true
			result.Items = map[string]interface{}{
				`group`: map[string]interface{}{
					`id_specialty`:       oldAdmission.Direction.IdParent,
					`code_specialty`:     sum.Code,
					`id_campaign`:        sum.IdCampaign,
					`id_education_level`: sum.IdEducationLevel,
					`budget_o`:           sum.BudgetO,
					`budget_oz`:          sum.BudgetOz,
					`budget_z`:           sum.BudgetZ,
					`paid_o`:             sum.PaidO,
					`paid_oz`:            sum.PaidOz,
					`paid_z`:             sum.PaidZ,
					`target_o`:           sum.TargetO,
					`target_oz`:          sum.TargetOz,
					`target_z`:           sum.TargetZ,
					`quota_o`:            sum.QuotaO,
					`quota_oz`:           sum.QuotaOz,
					`quota_z`:            sum.QuotaZ,
					`sum_distributed`:    sum.SumDistributed,
				},
				`id_edit_admission`: oldAdmission.Id,
			}
			return
		} else {
			result.SetErrorResult(db.Error.Error())
			tx.Rollback()
			return
		}
	} else {
		result.SetErrorResult(`Не совпадает направление`)
		tx.Rollback()
		return
	}
}

func (result *ResultInfo) EditAdmissionLevelBudget(data EditDistributedAdmissionVolume) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)
	var admission digest.AdmissionVolume
	db := conn.Find(&admission, data.IdAdmissionVolume)
	if admission.Id == 0 || db.Error != nil {
		result.SetErrorResult(`КЦП не найден`)
		tx.Rollback()
		return
	}
	err := CheckEditAdmission(admission.IdCampaign)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	var old digest.DistributedAdmissionVolume
	db = conn.Find(&old, data.IdLevelBudget)
	if old.Id == 0 || db.Error != nil {
		result.SetErrorResult(`Уровень бюджета не найден`)
		tx.Rollback()
		return
	}
	if old.IdOrganization != result.User.CurrentOrganization.Id {
		result.SetErrorResult(`Организация не соответствует`)
		tx.Rollback()
		return
	}
	var new digest.DistributedAdmissionVolume
	new.IdOrganization = old.IdOrganization
	t := time.Now()
	new = old
	new.Changed = &t
	new.IdAuthor = &result.User.Id
	new.BudgetO = data.BudgetO
	new.BudgetOz = data.BudgetOz
	new.BudgetZ = data.BudgetZ
	new.QuotaO = data.QuotaO
	new.QuotaOz = data.QuotaOz
	new.QuotaZ = data.QuotaZ
	new.PaidO = old.PaidO
	new.PaidOz = old.PaidOz
	new.PaidZ = old.PaidZ
	new.TargetO = data.TargetO
	new.TargetOz = data.TargetOz
	new.TargetZ = data.TargetZ
	if data.Uid != nil {
		var exist digest.DistributedAdmissionVolume
		db = tx.Where(`upper(uid)=upper(?) AND id_organization=?`, data.Uid, result.User.CurrentOrganization.Id).Find(&exist)
		if exist.Id > 0 {
			result.SetErrorResult(`Уровень бюджета с данным uid уже существуют`)
			tx.Rollback()
			return
		} else {
			new.Uid = data.Uid
		}
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&new)
	if db.Error == nil {
		tx.Commit()
		result.Done = true
		result.Items = map[string]interface{}{
			`id_edit_level_budget`: new.Id,
		}
		return
	} else {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}

}

func (result *ResultInfo) RemoveAdmission(IdAdmission uint) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var oldAdmission digest.AdmissionVolume
	db := tx.Preload(`Direction`).Find(&oldAdmission, IdAdmission)
	idEducLevel := oldAdmission.IdEducationLevel
	idCampaign := oldAdmission.IdCampaign
	idGroups := oldAdmission.Direction.IdParent
	if oldAdmission.Id == 0 || db.Error != nil {
		result.SetErrorResult(`Кцп не найдены`)
		tx.Rollback()
		return
	}
	err := CheckEditAdmission(oldAdmission.IdCampaign)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	var distributed []digest.DistributedAdmissionVolume
	db = conn.Where(`id_admission_volume=?`, IdAdmission).Find(&distributed)
	if len(distributed) > 0 {
		db = tx.Where(`id_admission_volume=?`, IdAdmission).Delete(&distributed)
	}
	db = tx.Table(`cmp.admission_volume`).Where(`id=?`, IdAdmission).Delete(&oldAdmission)
	if db.Error == nil {
		tx.Commit()
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
							SUM(av.quota_z)AS quota_z,
							SUM(vavsg.sum_distributed) AS sum_distributed
					   FROM cmp.admission_volume av  
					   JOIN cls.directions d ON d.id = av.id_direction
					   JOIN cls.v_okso_enlarged_group eg ON d.id_parent = eg.id
						LEFT JOIN cmp.v_admission_volume_specialty_groups vavsg ON vavsg.id = av.id
						Where av.id_campaign=? and av.id_education_level=? AND eg.id=?						
						 GROUP BY
						eg.code,
						av.id_campaign,
						av.id_education_level`, idCampaign, idEducLevel, idGroups).Scan(&sum)

		result.Done = true
		result.Items = map[string]interface{}{
			`group`: map[string]interface{}{
				`id_specialty`:       oldAdmission.Direction.IdParent,
				`code_specialty`:     sum.Code,
				`id_campaign`:        sum.IdCampaign,
				`id_education_level`: sum.IdEducationLevel,
				`budget_o`:           sum.BudgetO,
				`budget_oz`:          sum.BudgetOz,
				`budget_z`:           sum.BudgetZ,
				`paid_o`:             sum.PaidO,
				`paid_oz`:            sum.PaidOz,
				`paid_z`:             sum.PaidZ,
				`target_o`:           sum.TargetO,
				`target_oz`:          sum.TargetOz,
				`target_z`:           sum.TargetZ,
				`quota_o`:            sum.QuotaO,
				`quota_oz`:           sum.QuotaOz,
				`quota_z`:            sum.QuotaZ,
				`sum_distributed`:    sum.SumDistributed,
			},
			`id_remove_admission`: IdAdmission,
		}
		return
	} else {
		result.SetErrorResult(db.Error.Error())
		tx.Rollback()
		return
	}
}

func (result *ResultInfo) RemoveBudgetAdmission(IdAdmission uint, IdBudget uint) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	defer func() {
		tx.Rollback()
	}()
	conn.LogMode(config.Conf.Dblog)

	var oldAdmission digest.AdmissionVolume
	db := tx.Find(&oldAdmission, IdAdmission)
	if oldAdmission.Id == 0 || db.Error != nil {
		result.SetErrorResult(`Кцп не найдены`)
		tx.Rollback()
		return
	}
	err := CheckEditAdmission(oldAdmission.IdCampaign)
	if err != nil {
		result.SetErrorResult(err.Error())
		tx.Rollback()
		return
	}
	var distributed digest.DistributedAdmissionVolume
	db = conn.Where(`id_admission_volume=? AND id=?`, IdAdmission, IdBudget).Find(&distributed)
	if distributed.Id > 0 {
		db = tx.Where(`id_admission_volume=? AND id=?`, IdAdmission, IdBudget).Delete(&distributed)
		if db.Error == nil {
			result.Done = true
			tx.Commit()
			result.Items = map[string]interface{}{
				`id_admission`: IdAdmission,
				`id_budget`:    IdBudget,
			}
			return
		} else {
			result.SetErrorResult(db.Error.Error())
			tx.Rollback()
			return
		}
	} else {
		result.SetErrorResult(`Не найден уровень бюджета`)
		tx.Rollback()
		return
	}
}

func (result *ResultInfo) GetLevelBudgetAdmission(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var idsLevelsBudget []uint
	db := conn.Table(`cmp.distributed_admission_volume`).Where(`id_admission_volume=?`, ID).Pluck("id_level_budget", &idsLevelsBudget)
	if db.Error == nil {
		result.Done = true
		result.Items = idsLevelsBudget
	} else {
		result.SetErrorResult(`Ошибка БД`)
		return
	}

}
