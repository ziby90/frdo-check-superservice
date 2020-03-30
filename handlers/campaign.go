package handlers

import (
	"persons/config"
	"persons/digest"
	"persons/error_handler"
	"time"
)


var validBase64Err = error_handler.ErrorType{Type: 1, ToUserType: 500}

type CampaignMain struct {
	Id               	uint            `json:"id"`		// Идентификатор
	UID              	string          `json:"uid"`		// Идентификатор от организации
	Name             	string        	`json:"name"` 		// Наименование
	IdCampaignType   	uint			`json:"id_campaign_type"`
	NameCampaignType 	string			`json:"campaign_type_name"`
	IdCampaignStatus 	uint			`json:"id_campaign_status"`
	NameCampaignStatus 	string			`json:"campaign_status_name"`
	EducationForms		[]EducationFormRespons	`json:"education_forms"`
	EducationLevels		[]EducationLevelRespons	`json:"education_levels"`
	YearStart        	int64        	`json:"year_start"`            // Год начала компании
	YearEnd          	int64       	`json:"year_end"`                // Год окончания компании
	Created          	time.Time    	`json:"created"`			// Дата создания
}


type CampaignResponse struct {
	Id               	uint            `json:"id"`		// Идентификатор
	UID              	string          `json:"uid"`		// Идентификатор от организации
	Name             	string        	`json:"name"` 		// Наименование
	IdCampaignType   	uint			`json:"id_campaign_type"`
	NameCampaignType 	string			`json:"campaign_type_name"`
	IdCampaignStatus 	uint			`json:"id_campaign_status"`
	NameCampaignStatus 	string			`json:"campaign_status_name"`
	YearStart        	int64        	`json:"year_start"`            // Год начала компании
	YearEnd          	int64       	`json:"year_end"`                // Год окончания компании
	Created          	time.Time    	`json:"created"`			// Дата создания
}

func (result *Result) GetListCampaign() {
	result.Done = false
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var campaigns []digest.Campaign
	db := conn.Order(result.Params.Sort.Field+` `+result.Params.Sort.Order).Limit(result.Params.Paginator.Limit).Offset(result.Params.Paginator.Offset).Find(&campaigns)
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
	if db.RowsAffected>0 {
		for _, campaign := range campaigns{
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
	}
}

func (result *Result) GetInfoCampaign(ID uint) {
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
	if db.RowsAffected>0 {
		db = conn.Model(&campaign).Related(&campaign.CampaignType, `IdCampaignType`)
		db = conn.Model(&campaign).Related(&campaign.CampaignStatus, `IdCampaignStatus`)

		//var educForms []EducationFormRespons
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
			c.EducationForms = append(c.EducationForms, GetEducFormResponse(campEducForm.IdEducationForm))
		}
		for _, campEducLevel := range campEducLevels {
			c.EducationLevels = append(c.EducationLevels, GetEducLevelResponse(campEducLevel.IdEducationLevel))
		}
		result.Done = true
		result.Items = c
		return
	}
}
