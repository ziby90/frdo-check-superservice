package handlers

import (
	"github.com/jinzhu/gorm"
	"persons/config"
	"strings"
)

type CampaignStatusesRespons struct {
	Id   uint   `json:"id"`
	Name string `json:"text"`
}
type BenefitsResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"text"`
}
type CampaignTypesRespons struct {
	Id   uint   `json:"id"`
	Name string `json:"text"`
}

type DirectionsRespons struct {
	Id   uint   `json:"id"`
	Name string `json:"text"`
}

type EducationFormRespons struct {
	Id   uint   `json:"id"`
	Name string `json:"text"`
}

type EducationLevelRespons struct {
	Id   uint   `json:"id"`
	Name string `json:"text"`
}

func (BenefitsResponse) TableName() string {
	return "cls.benefits"
}
func (CampaignStatusesRespons) TableName() string {
	return "cls.campaign_statuses"
}
func (CampaignTypesRespons) TableName() string {
	return "cls.campaign_types"
}
func (DirectionsRespons) TableName() string {
	return "cls.directions"
}
func (EducationFormRespons) TableName() string {
	return "cls.education_forms"
}
func (EducationLevelRespons) TableName() string {
	return "cls.education_levels"
}

func (result *ResultCls) GetClsResponse(clsName string) {
	conn := config.Db.ConnGORM
	conn.LogMode(config.Conf.Dblog)
	var db *gorm.DB
	switch clsName {
	case `benefits`:
		var r []BenefitsResponse
		if result.Search != `` {
			db = conn.Where(`UPPER(name) like ?`, `%`+strings.ToUpper(result.Search)+`%`).Find(&r)
		} else {
			db = conn.Find(&r)
		}
		result.Items = r
		break
	case `campaign_statuses`:
		var r []CampaignStatusesRespons
		if result.Search != `` {
			db = conn.Where(`UPPER(name) like ?`, `%`+strings.ToUpper(result.Search)+`%`).Find(&r)
		} else {
			db = conn.Find(&r)
		}
		result.Items = r
		break
	case `campaign_types`:
		var r []CampaignTypesRespons
		if result.Search != `` {
			db = conn.Where(`UPPER(name) like ?`, `%`+strings.ToUpper(result.Search)+`%`).Find(&r)
		} else {
			db = conn.Find(&r)
		}
		result.Items = r
		break
	case `directions`:
		var r []DirectionsRespons
		if result.Search != `` {
			db = conn.Where(`UPPER(name) like ?`, `%`+strings.ToUpper(result.Search)+`%`).Find(&r)
		} else {
			db = conn.Find(&r)
		}
		result.Items = r
		break
	case `education_levels`:
		var r []EducationLevelRespons
		if result.Search != `` {
			db = conn.Where(`UPPER(name) like ?`, `%`+strings.ToUpper(result.Search)+`%`).Find(&r)
		} else {
			db = conn.Find(&r)
		}
		result.Items = r
		break
	case `education_forms`:
		var r []EducationFormRespons
		if result.Search != `` {
			db = conn.Where(`UPPER(name) like ?`, `%`+strings.ToUpper(result.Search)+`%`).Find(&r)
		} else {
			db = conn.Find(&r)
		}
		result.Items = r
		break
	default:
		message := `Неизвестный справочник.`
		result.Message = &message
		return
	}
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
	if db.RowsAffected > 0 {
		result.Done = true
		return
	}
}
