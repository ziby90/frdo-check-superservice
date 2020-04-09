package handlers

import (
	"persons/config"
	"persons/digest"
	"persons/service"
	"strings"
	"time"
)

type CampaignMain struct {
	Id                  uint      `json:"id"`   // Идентификатор
	UID                 string    `json:"uid"`  // Идентификатор от организации
	Name                string    `json:"name"` // Наименование
	IdCampaignType      uint      `json:"id_campaign_type"`
	NameCampaignType    string    `json:"campaign_type_name"`
	IdCampaignStatus    uint      `json:"id_campaign_status"`
	NameCampaignStatus  string    `json:"campaign_status_name"`
	EducationForms      []uint    `json:"education_forms"`
	EducationFormsName  []string  `json:"education_forms_names"`
	EducationLevels     []uint    `json:"education_levels"`
	EducationLevelsName []string  `json:"education_levels_names"`
	YearStart           int64     `json:"year_start"` // Год начала компании
	YearEnd             int64     `json:"year_end"`   // Год окончания компании
	Created             time.Time `json:"created"`    // Дата создания

}

type CampaignResponse struct {
	Id                 uint      `json:"id"`   // Идентификатор
	UID                string    `json:"uid"`  // Идентификатор от организации
	Name               string    `json:"name"` // Наименование
	IdCampaignType     uint      `json:"id_campaign_type"`
	NameCampaignType   string    `json:"campaign_type_name"`
	IdCampaignStatus   uint      `json:"id_campaign_status"`
	NameCampaignStatus string    `json:"campaign_status_name"`
	YearStart          int64     `json:"year_start"` // Год начала компании
	YearEnd            int64     `json:"year_end"`   // Год окончания компании
	Created            time.Time `json:"created"`    // Дата создания
}

var campaignSearchArray = []string{
	`name`,
}

func (result *Result) GetListCampaign() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaigns []digest.Campaign
	db := conn.Where(`id_organization=?`, result.User.CurrentOrganization.Id).Order(result.Sort.Field + ` ` + result.Sort.Order)
	for _, search := range result.Search {
		if service.SearchStringInSliceString(search[0], campaignSearchArray) >= 0 {
			db = db.Where(`UPPER(`+search[0]+`) LIKE ?`, `%`+strings.ToUpper(search[1])+`%`)
		}
	}
	dbCount := db.Model(&campaigns).Count(&result.Paginator.TotalCount)
	if dbCount.Error != nil {

	}
	result.Paginator.Make()
	db = db.Limit(result.Paginator.Limit).Offset(result.Paginator.Offset).Find(&campaigns)
	var campaignResponses []CampaignResponse
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
		for _, campaign := range campaigns {
			db = conn.Model(&campaign).Related(&campaign.CampaignType, `IdCampaignType`)
			db = conn.Model(&campaign).Related(&campaign.CampaignStatus, `IdCampaignStatus`)
			c := CampaignResponse{
				Id:                 campaign.Id,
				UID:                campaign.Uid,
				Name:               campaign.Name,
				IdCampaignType:     campaign.CampaignType.Id,
				NameCampaignType:   campaign.CampaignType.Name,
				IdCampaignStatus:   campaign.CampaignStatus.Id,
				NameCampaignStatus: campaign.CampaignStatus.Name,
				YearStart:          campaign.YearStart,
				YearEnd:            campaign.YearEnd,
				Created:            campaign.Created,
			}
			campaignResponses = append(campaignResponses, c)
		}
		result.Done = true
		result.Items = campaignResponses
		return
	} else {
		result.Done = true
		message := `Компании не найдены.`
		result.Message = &message
		result.Items = []digest.Campaign{}
		return
	}
}

func (result *ResultInfo) GetInfoCampaign(ID uint) {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaign digest.Campaign
	db := conn.Find(&campaign, ID)
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
		db = conn.Model(&campaign).Related(&campaign.CampaignType, `IdCampaignType`)
		db = conn.Model(&campaign).Related(&campaign.CampaignStatus, `IdCampaignStatus`)

		var campEducForms []digest.CampaignEducForm
		db = conn.Where(`id_campaign=?`, campaign.Id).Find(&campEducForms)
		var campEducLevels []digest.CampaignEducLevel
		db = conn.Where(`id_campaign=?`, campaign.Id).Find(&campEducLevels)

		c := CampaignMain{
			Id:                 campaign.Id,
			UID:                campaign.Uid,
			Name:               campaign.Name,
			IdCampaignType:     campaign.CampaignType.Id,
			NameCampaignType:   campaign.CampaignType.Name,
			IdCampaignStatus:   campaign.CampaignStatus.Id,
			NameCampaignStatus: campaign.CampaignStatus.Name,
			YearStart:          campaign.YearStart,
			YearEnd:            campaign.YearEnd,
			Created:            campaign.Created,
		}
		for _, campEducForm := range campEducForms {
			var educForm digest.EducationForm
			db = conn.Find(&educForm, campEducForm.IdEducationForm)
			c.EducationForms = append(c.EducationForms, educForm.Id)
			c.EducationFormsName = append(c.EducationFormsName, educForm.Name)
		}
		for _, campEducLevel := range campEducLevels {
			var educLevel digest.EducationLevel
			db = conn.Find(&educLevel, campEducLevel.IdEducationLevel)
			c.EducationLevels = append(c.EducationLevels, educLevel.Id)
			c.EducationLevelsName = append(c.EducationLevelsName, educLevel.Name)
		}
		result.Done = true
		result.Items = c
		return
	} else {
		result.Done = true
		message := `Компании не найдены.`
		result.Message = &message
		result.Items = []digest.Campaign{}
		return
	}
}

func (result *ResultInfo) AddCampaign(campaignData CampaignMain, user digest.User) {
	conn := config.Db.ConnGORM
	tx := conn.Begin()
	conn.LogMode(config.Conf.Dblog)

	var campaign digest.Campaign
	campaign.Organization.Id = user.CurrentOrganization.Id
	campaign.IdOrganization = user.CurrentOrganization.Id
	campaign.IdAuthor = user.Id
	campaign.Created = time.Now()
	campaign.Name = campaignData.Name

	campaign.IdCampaignType = campaignData.IdCampaignType
	campaign.CampaignType.Id = campaignData.IdCampaignType
	db := tx.Find(&campaign.CampaignType, campaign.CampaignType.Id)
	if db.Error != nil || !campaign.CampaignType.Actual {
		result.SetErrorResult(`Не найден тип компании`)
		return
	}

	campaign.IdCampaignStatus = campaignData.IdCampaignStatus
	campaign.CampaignStatus.Id = campaignData.IdCampaignStatus
	// проверка типа
	db = tx.Find(&campaign.CampaignStatus, campaign.CampaignStatus.Id)
	if db.Error != nil || !campaign.CampaignStatus.Actual {
		result.SetErrorResult(`Статус комании не найден`)
		return
	}

	campaign.YearEnd = campaignData.YearEnd
	// проверка года окончания
	if int(campaignData.YearEnd) < 1900 || int(campaignData.YearEnd) > time.Now().Year() {
		result.SetErrorResult(`Год окончания за пределами`)
		return
	}

	campaign.YearStart = campaignData.YearStart
	// проверка года начала
	if int(campaignData.YearStart) < 1900 || int(campaignData.YearStart) > time.Now().Year() {
		result.SetErrorResult(`Год начала за пределами`)
		return
	}

	if campaignData.YearStart > campaign.YearEnd {
		result.SetErrorResult(`Год начала не может быть позже года окончания`)
		return
	}

	db = tx.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&campaign)
	if db.Error != nil {
		tx.Rollback()
		m := db.Error.Error()
		result.Message = &m
		return
	}

	for _, educLevelId := range campaignData.EducationLevels {
		var educationLevel digest.EducationLevel
		tx.Find(&educationLevel, educLevelId)
		if !educationLevel.Actual {
			result.SetErrorResult(`Уровень образования не найден`)
			return
		}
		row := conn.Table(`cls.edu_levels_campaign_types`).Where(`id_campaign_types=? AND id_education_level=?`, campaign.CampaignType.Id, educLevelId).Select(`id`).Row()
		var idEducLevelCampaignType uint
		err := row.Scan(&idEducLevelCampaignType)
		if err != nil && idEducLevelCampaignType <= 0 {
			result.SetErrorResult(`Данный уровень образования не соответствует типу приемной компании`)
			return
		}
		campaignEducLevel := digest.CampaignEducLevel{
			IdCampaign:       campaign.Id,
			IdEducationLevel: educLevelId,
		}
		db = tx.Create(&campaignEducLevel)
	}

	for _, educFormId := range campaignData.EducationForms {
		var educationForm digest.EducationForm
		tx.Find(&educationForm, educFormId)
		if !educationForm.Actual {
			result.SetErrorResult(`Форма образования не найдена`)
			return
		}
		campaignEducForm := digest.CampaignEducForm{
			IdCampaign:      campaign.Id,
			IdEducationForm: educFormId,
		}
		db = tx.Create(&campaignEducForm)
	}

	result.Items = campaign.Id
	result.Done = true
	tx.Commit()
}
